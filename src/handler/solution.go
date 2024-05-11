package handler

import (
	"context"
	"goto/src/config"
	"goto/src/database/query"
	"goto/src/model"
	"goto/src/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	s "strings"

	cp "github.com/otiai10/copy"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSolutions(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)
	pager := utils.NewPager(c)
	solutions := query.GetUserSolutions(ctx, user.Id, pager)

	response := make([]model.SolutionMin, len(solutions))
	for i, s := range solutions {
		response[i] = *s.Min()
	}
	return c.JSON(response)
}

func GetSolution(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	solution := query.GetUserSolution(ctx, id, user.Id)
	if solution == nil {
		return c.Status(404).SendString("Solution not found")
	}

	return c.JSON(solution)
}

func SubmitSolution(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	solutionBody := model.SolutionInput{}
	if err := c.BodyParser(&solutionBody); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad solution format")
	}

	task := query.GetTask(ctx, solutionBody.TaskId)
	if task == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Task not found")
	}

	taskFileNames := make([]string, len(task.Files))
	for i, tf := range task.Files {
		taskFileNames[i] = tf.Name
	}
	solutionFileNames := make([]string, len(solutionBody.Files))
	for i, sf := range solutionBody.Files {
		solutionFileNames[i] = sf.Name
	}
	missingFileNames := utils.Difference(taskFileNames, solutionFileNames)
	if len(missingFileNames) > 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Missing files", "details": missingFileNames})
	}

	solution := &model.Solution{
		UserId: user.Id,
		TaskId: solutionBody.TaskId,
		Files:  solutionBody.Files,
	}
	solution = query.SaveSolution(ctx, solution)
	go checkSolution(solution, task)

	return c.JSON(solution)
}

func checkSolution(solution *model.Solution, task *model.Task) {
	ctx := context.Background()

	solutionFilesByNames := map[string]model.SolutionFile{}
	for _, sf := range solution.Files {
		solutionFilesByNames[sf.Name] = sf
	}

	project := query.GetProjectShallow(ctx, task.ProjectId)
	tempDir := uuid.New().String()
	projectTempPath := filepath.Join(config.TempPath, tempDir)
	projectPath := filepath.Join(config.MediaPath, project.Dir)
	cp.Copy(projectPath, projectTempPath)
	defer os.RemoveAll(projectTempPath)

	srcPath := filepath.Join(projectTempPath, project.SrcDir)
	for _, tf := range task.Files {
		path := filepath.Join(srcPath, tf.Path)
		code := solutionFilesByNames[tf.Name].Code
		os.WriteFile(path, []byte(code), os.ModePerm)
	}

	switch project.Containerization {
	case "docker":
		buildCmd := exec.Command("docker", "build", "-q", projectTempPath)
		tempImage, _ := buildCmd.Output()
		runCmd := exec.Command(
			"docker",
			"run",
			"-e",
			"TARGET="+task.RunTarget,
			"--rm",
			"-t",
			s.TrimSuffix(string(tempImage), "\n"),
		)
		output, _ := runCmd.Output()
		solution.Result = string(output)
		exec.Command("docker", "system", "prune", "-f").Run()
	case "docker-compose":
		upCmd := exec.Command("docker", "compose", "up", "--build", "--abort-on-container-exit")
		upCmd.Env = append(upCmd.Env, "TARGET="+task.RunTarget)
		upCmd.Dir = projectTempPath
		output, _ := upCmd.Output()
		solution.Result = utils.ParseComposeOutput(output, project.Dir)
		downCmd := exec.Command(
			"docker",
			"compose",
			"down",
			"--rmi",
			"local",
			"-v",
			"--remove-orphans",
		)
		downCmd.Dir = projectTempPath
		downCmd.Run()
		exec.Command("docker", "system", "prune", "-f").Run()
	}

	solution.Status = utils.ParseStatus(solution.Result)
	query.SaveSolution(ctx, solution)
}