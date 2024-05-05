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
	"time"

	cp "github.com/otiai10/copy"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSolutions(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)
	solutions := query.GetUserSolutions(ctx, user.Id)
	return c.JSON(solutions)
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

	solution := model.SolutionBase{}
	if err := c.BodyParser(&solution); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad solution format")
	}

	task := query.GetTask(ctx, solution.TaskId)
	if task == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Task not found")
	}

	taskFileNames := utils.MapKeys(task.Files)
	solutionFileNames := utils.MapKeys(solution.Files)
	missingFileNames := utils.Difference(taskFileNames, solutionFileNames)
	if len(missingFileNames) > 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Missing files", "details": missingFileNames})
	}

	project := query.GetProjectShallow(ctx, task.ProjectId)
	tempDir := uuid.New().String()
	projectTempPath := filepath.Join(config.TempPath, tempDir)
	projectPath := filepath.Join(config.MediaPath, project.Dir)
	cp.Copy(projectPath, projectTempPath)
	defer os.RemoveAll(projectTempPath)

	srcPath := filepath.Join(projectTempPath, project.SrcDir)
	for taskFileName, taskFilePath := range task.Files {
		path := filepath.Join(srcPath, taskFilePath)
		code := solution.Files[taskFileName]
		os.WriteFile(path, []byte(code), os.ModePerm)
	}

	var result string
	switch project.Containerization {
	case "docker":
		buildCmd := exec.Command("docker", "build", "-q", projectTempPath)
		tempImage, _ := buildCmd.Output()
		runCmd := exec.Command("docker", "run", "--rm", "-t", string(tempImage))
		output, _ := runCmd.Output()
		result = string(output)

		// TODO

		exec.Command("docker", "system", "prune", "-f").Run()
	case "docker-compose":
		upCmd := exec.Command("docker", "compose", "up", "--build", "--abort-on-container-exit")
		upCmd.Env = append(upCmd.Env, "TARGET="+task.RunTarget)
		upCmd.Dir = projectTempPath
		output, _ := upCmd.Output()

		result = utils.ParseComposeOutput(output, project.Dir)

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

	return c.JSON(model.Solution{
		SolutionBase: solution,
		UserId:       user.Id,
		Status:       utils.ParseStatus(result),
		Result:       result,
		UpdatedAt:    time.Now(),
	})
}