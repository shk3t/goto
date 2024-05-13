package handler

import (
	"context"
	"goto/src/config"
	q "goto/src/database/query"
	f "goto/src/filter"
	m "goto/src/model"
	"goto/src/service"
	u "goto/src/utils"
	"os"
	"os/exec"
	"path/filepath"
	sc "strconv"
	s "strings"

	cp "github.com/otiai10/copy"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSolutions(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)
	pager := service.NewPager(fctx)
	filter := f.NewSolutionFilter(fctx)
	solutions := q.GetSolutions(ctx, user.Id, pager, filter)
	return fctx.JSON(solutions.Min())
}

func GetSolution(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	solution := q.GetUserSolution(ctx, id, user.Id)
	if solution == nil {
		return fctx.Status(404).SendString("Solution not found")
	}

	return fctx.JSON(solution)
}

func SubmitSolution(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)

	solutionBody := m.SolutionInput{}
	if err := fctx.BodyParser(&solutionBody); err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Bad solution format")
	}

	task := q.GetTask(ctx, solutionBody.TaskId)
	if task == nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Task not found")
	}

	validSolutionFiles := validateFileNames(fctx, solutionBody.Files, task.Files)
	if validSolutionFiles == nil {
		return nil
	}

	solution := &m.Solution{
		UserId: user.Id,
		Files:  validSolutionFiles,
		Task:   *task.Min(),
	}
	q.SaveSolution(ctx, solution)
	solution = q.GetUserSolution(ctx, solution.Id, solution.UserId)

	go checkSolution(solution, task)

	return fctx.JSON(solution)
}

func validateFileNames(
	fctx *fiber.Ctx,
	solutionFiles m.SolutionFiles,
	taskFiles m.TaskFiles,
) (m.SolutionFiles) {
	taskFileNames := make([]string, len(taskFiles))
	for i, tf := range taskFiles {
		taskFileNames[i] = tf.Name
	}
	solutionFilesByNames := map[string]m.SolutionFile{}
	for _, sf := range solutionFiles {
		solutionFilesByNames[sf.Name] = sf
	}
	solutionFileNames := u.MapKeys(solutionFilesByNames)

	missingFileNames := u.Difference(taskFileNames, solutionFileNames)
	if len(missingFileNames) > 0 {
		fctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Missing files", "details": missingFileNames})
		return nil
	}

	redundantFileNames := u.Difference(solutionFileNames, taskFileNames)
	for _, redName := range redundantFileNames {
		delete(solutionFilesByNames, redName)
	}

	return u.MapValues(solutionFilesByNames)
}

func checkSolution(solution *m.Solution, task *m.Task) {
	ctx := context.Background()

	solutionFilesByNames := map[string]m.SolutionFile{}
	for _, sf := range solution.Files {
		solutionFilesByNames[sf.Name] = sf
	}

	project := q.GetProjectShallow(ctx, task.ProjectId)
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
		solution.Result = service.ParseComposeOutput(output, project.Dir)
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

	solution.Status = service.ParseStatus(solution.Result)
	q.SaveSolution(ctx, solution)
}