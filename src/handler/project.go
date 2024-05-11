package handler

import (
	"context"
	"errors"
	"goto/src/config"
	q "goto/src/database/query"
	f "goto/src/filter"
	m "goto/src/model"
	"goto/src/service"
	u "goto/src/utils"
	"os"
	sc "strconv"
	s "strings"

	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetProjects(fctx *fiber.Ctx) error {
	ctx := context.Background()
	pager := service.NewPager(fctx)
	filter := f.NewProjectFilter(fctx)
	projects := q.GetProjects(ctx, pager, filter)
	return fctx.JSON(projects.Min())
}

func GetProject(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	project := q.GetUserProject(ctx, id, user.Id)
	if project == nil {
		return fctx.Status(404).SendString("Project not found")
	}

	return fctx.JSON(project.Public())
}

func LoadProject(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)
	body := struct{ Url string }{}
	if err := fctx.BodyParser(&body); err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	postfix := uuid.New().String()
	var delayedTask *m.DelayedTask

	if body.Url != "" {
		urlParts := s.Split(body.Url, "/")
		projectDir := urlParts[len(urlParts)-1]
		projectName := projectDir + "_" + postfix

		gitCloneCmd := exec.Command("git", "clone", body.Url, projectName)
		gitCloneCmd.Dir = config.MediaPath
		if err := gitCloneCmd.Run(); err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString("Invalid url")
		}

		delayedTask = &m.DelayedTask{
			UserId: user.Id,
			Action: "create project",
			Target: projectDir,
		}
		q.SaveDelayedTask(ctx, delayedTask)
		go postCreateProject(user, delayedTask, projectName)

	} else {
		file, err := fctx.FormFile("file")
		if err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString("Use `file` as a key for uploaded file")
		}

		projectDir, extension := u.SplitExt(file.Filename)
		projectName := projectDir + "_" + postfix

		archivePath := filepath.Join(config.MediaPath, projectName+extension)
		if err = fctx.SaveFile(file, archivePath); err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		delayedTask = &m.DelayedTask{
			UserId: user.Id,
			Action: "create project",
			Target: projectDir,
		}
		q.SaveDelayedTask(ctx, delayedTask)
		go postCreateProjectZip(user, delayedTask, projectName, archivePath)
	}

	return fctx.JSON(delayedTask)
}

func postCreateProject(user *m.User, delayedTask *m.DelayedTask, projectName string) {
	ctx := context.Background()

	projectPath := filepath.Join(config.MediaPath, projectName)
	gotoConfigPath := filepath.Join(projectPath, config.GotoConfigName)
	gotoConfig, err := m.LoadGotoConfig(gotoConfigPath)
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

	if err = q.CreateProject(ctx, project); err != nil {
		errorDelayedTask(ctx, delayedTask, err)
		deleteProject(ctx, project)
		return
	}

	okDelayedTask(ctx, delayedTask, project.Id, "Created") // TODO return good Id
}

func postCreateProjectZip(
	user *m.User,
	delayedTask *m.DelayedTask,
	projectName string,
	archivePath string,
) {
	ctx := context.Background()

	if err := service.Unzip(archivePath, true); err != nil {
		errorDelayedTask(ctx, delayedTask, err)
		return
	}
	os.Remove(archivePath)

	postCreateProject(user, delayedTask, projectName)
}

func errorDelayedTask(ctx context.Context, delayedTask *m.DelayedTask, err error) {
	delayedTask.Status = "error"
	delayedTask.Details = err.Error()
	q.SaveDelayedTask(ctx, delayedTask)
}

func okDelayedTask(
	ctx context.Context,
	delayedTask *m.DelayedTask,
	targetId int,
	details string,
) {
	delayedTask.TargetId = &targetId
	delayedTask.Status = "ok"
	delayedTask.Details = details
	q.SaveDelayedTask(ctx, delayedTask)
}

func UpdateProject(fctx *fiber.Ctx) error {
	return errors.New("TODO")
}

func DeleteProject(fctx *fiber.Ctx) error {
	ctx := context.Background()

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	user := service.GetCurrentUser(fctx)
	project := q.GetUserProject(ctx, id, user.Id)
	if project == nil {
		return fctx.Status(404).SendString("Project not found")
	}

	deleteProject(ctx, project)

	return fctx.SendStatus(fiber.StatusOK)
}

func deleteProject(ctx context.Context, project *m.Project) {
	if project.Id != 0 {
		q.DeleteProject(ctx, project.Id)
	}
	cleanupContainers(project)
	os.RemoveAll(filepath.Join(config.MediaPath, project.Dir))
}

func cleanupContainers(project *m.Project) {
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