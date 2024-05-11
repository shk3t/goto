package handler

import (
	"context"
	"errors"
	"goto/src/config"
	"goto/src/database/query"
	"goto/src/model"
	"goto/src/utils"
	"os"
	sc "strconv"
	s "strings"

	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetProjects(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)
	pager := utils.NewPager(c)

	projects := query.GetUserProjects(ctx, user.Id, pager)

	response := make([]model.ProjectMin, len(projects))
	for i, p := range projects {
		response[i] = *p.Min()
	}
	return c.JSON(response)
}

func GetProject(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	id, err := sc.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	project := query.GetUserProject(ctx, id, user.Id)
	if project == nil {
		return c.Status(404).SendString("Project not found")
	}

	return c.JSON(project.Public())
}

func LoadProject(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)
	body := struct{ Url string }{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	postfix := uuid.New().String()
	var delayedTask *model.DelayedTask

	if body.Url != "" {
		urlParts := s.Split(body.Url, "/")
		projectDir := urlParts[len(urlParts)-1]
		projectName := projectDir + "_" + postfix

		gitCloneCmd := exec.Command("git", "clone", body.Url, projectName)
		gitCloneCmd.Dir = config.MediaPath
		if err := gitCloneCmd.Run(); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid url")
		}

		delayedTask = &model.DelayedTask{
			UserId: user.Id,
			Action: "create project",
			Target: projectDir,
		}
		query.SaveDelayedTask(ctx, delayedTask)
		go postCreateProject(user, delayedTask, projectName)

	} else {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Use `file` as a key for uploaded file")
		}

		projectDir, extension := utils.SplitExt(file.Filename)
		projectName := projectDir + "_" + postfix

		archivePath := filepath.Join(config.MediaPath, projectName+extension)
		if err = c.SaveFile(file, archivePath); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		delayedTask = &model.DelayedTask{
			UserId: user.Id,
			Action: "create project",
			Target: projectDir,
		}
		query.SaveDelayedTask(ctx, delayedTask)
		go postCreateProjectZip(user, delayedTask, projectName, archivePath)
	}

	return c.JSON(delayedTask)
}

func postCreateProject(user *model.User, delayedTask *model.DelayedTask, projectName string) {
	ctx := context.Background()

	projectPath := filepath.Join(config.MediaPath, projectName)
	gotoConfigPath := filepath.Join(projectPath, config.GotoConfigName)
	gotoConfig, err := model.LoadGotoConfig(gotoConfigPath)
	if err != nil {
		errorDelayedTask(ctx, delayedTask, err)
		os.RemoveAll(projectPath)
		return
	}

	project := gotoConfig.Project()
	project.Dir = projectName
	project.UserId = user.Id
	for _, t := range project.Tasks {
		for j, tf := range t.Files {
			stubPath := filepath.Join(config.MediaPath, project.Dir, project.StubDir, tf.Path)
			stubBytes, _ := os.ReadFile(stubPath)
			t.Files[j].Stub = string(stubBytes)
		}
	}

	var prepareCmds []*exec.Cmd
	switch gotoConfig.Containerization {
	case "docker":
		prepareCmds = []*exec.Cmd{
			exec.Command("docker", "buildx", "build", "-t", projectName, "."),
		}
	case "docker-compose":
		prepareCmds = []*exec.Cmd{
			exec.Command("docker", "compose", "pull"),
			exec.Command("docker", "compose", "build"),
		}
	default:
		err = errors.New("Specified containerization type is not supported")
		errorDelayedTask(ctx, delayedTask, err)
		deleteProject(ctx, project)
		return
	}

	for _, cmd := range prepareCmds {
		cmd.Dir = projectPath
		if err = cmd.Run(); err != nil {
			errorDelayedTask(ctx, delayedTask, err)
			deleteProject(ctx, project)
			return
		}
	}

	if err = query.CreateProject(ctx, project); err != nil {
		errorDelayedTask(ctx, delayedTask, err)
		deleteProject(ctx, project)
		return
	}

	okDelayedTask(ctx, delayedTask, project.Id, "Created")  // TODO return good Id
}

func postCreateProjectZip(
	user *model.User,
	delayedTask *model.DelayedTask,
	projectName string,
	archivePath string,
) {
	ctx := context.Background()

	if err := utils.Unzip(archivePath, true); err != nil {
		errorDelayedTask(ctx, delayedTask, err)
		return
	}
	os.Remove(archivePath)

	postCreateProject(user, delayedTask, projectName)
}

func errorDelayedTask(ctx context.Context, delayedTask *model.DelayedTask, err error) {
	delayedTask.Status = "error"
	delayedTask.Details = err.Error()
	query.SaveDelayedTask(ctx, delayedTask)
}

func okDelayedTask(
	ctx context.Context,
	delayedTask *model.DelayedTask,
	targetId int,
	details string,
) {
	delayedTask.TargetId = &targetId
	delayedTask.Status = "ok"
	delayedTask.Details = details
	query.SaveDelayedTask(ctx, delayedTask)
}

func UpdateProject(c *fiber.Ctx) error {
	return errors.New("TODO")
}

func DeleteProject(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := sc.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	user := GetCurrentUser(c)
	project := query.GetUserProject(ctx, id, user.Id)
	if project == nil {
		return c.Status(404).SendString("Project not found")
	}

	deleteProject(ctx, project)

	return c.SendStatus(fiber.StatusOK)
}

func deleteProject(ctx context.Context, project *model.Project) {
	if project.Id != 0 {
		query.DeleteProject(ctx, project.Id)
	}
    cleanupContainers(project)
	os.RemoveAll(filepath.Join(config.MediaPath, project.Dir))
}

func cleanupContainers(project *model.Project) {
	switch project.Containerization {
	case "docker":
		exec.Command("docker", "image", "remove", "-f", project.Dir).Run()
		exec.Command("docker", "system", "prune", "-f").Run()
	case "docker-compose":
		removeCmd := exec.Command(
			"docker",
			"compose",
			"down",
			"--rmi",
			"all",
			"-v",
			"--remove-orphans",
		)
		removeCmd.Dir = filepath.Join(config.MediaPath, project.Dir)
		removeCmd.Run()
		exec.Command("docker", "system", "prune", "-f").Run()
	}
}