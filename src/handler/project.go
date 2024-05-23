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
	return saveProject(fctx)
}

func saveProject(fctx *fiber.Ctx) error {
	var err error
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)

	id := 0
	action := "create"
	if fctx.Method() == "PUT" {
		action = "update"
		id, err = sc.Atoi(fctx.Params("id"))
		if err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
		}
		project := q.GetProjectShallow(ctx, id)
		if project == nil || project.UserId != user.Id {
			return fctx.Status(404).SendString("Project not found")
		}
	}

	body := struct{ Url string }{}
	if err := fctx.BodyParser(&body); err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	postfix := uuid.New().String()
	var delayedTask *m.DelayedTask
	var projectDir, projectName string

	if body.Url != "" {
		urlParts := s.Split(body.Url, "/")
		projectDir = urlParts[len(urlParts)-1]
		projectName = projectDir + "_" + postfix

		gitCloneCmd := exec.Command("git", "clone", body.Url, projectName)
		gitCloneCmd.Dir = config.MediaPath
		if err := gitCloneCmd.Run(); err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString("Invalid url")
		}

	} else {
		file, err := fctx.FormFile("file")
		if err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString("Use `file` as a key for uploaded file")
		}

		projectDir, extension := u.SplitExt(file.Filename)
		projectName = projectDir + "_" + postfix

		archivePath := filepath.Join(config.MediaPath, projectName+extension)
		if err = fctx.SaveFile(file, archivePath); err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		err = service.Unzip(archivePath, true)
		os.Remove(archivePath)
		if err != nil {
			return fctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
	}

	delayedTask = &m.DelayedTask{
		UserId: user.Id,
		Action: action + " project",
		Target: projectDir,
	}
	q.SaveDelayedTask(ctx, delayedTask)
	go postSaveProject(id, projectName, delayedTask)

	return fctx.JSON(delayedTask)
}

func postSaveProject(projectId int, projectName string, delayedTask *m.DelayedTask) {
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
	project.Id = projectId
	project.Dir = projectName
	project.UserId = delayedTask.UserId
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
		deleteProject(ctx, 0, project.Dir, project.Containerization)
		return
	}

	for _, cmd := range prepareCmds {
		cmd.Dir = projectPath
		if err = cmd.Run(); err != nil {
			errorDelayedTask(ctx, delayedTask, err)
			deleteProject(ctx, 0, project.Dir, project.Containerization)
			return
		}
	}

	projectOldDir := ""
	if project.Id != 0 {
		projectOldDir = q.GetProjectShallow(ctx, project.Id).Dir
	}
	err = q.SaveProject(ctx, project, gotoConfig)
	if err != nil {
		errorDelayedTask(ctx, delayedTask, err)
		deleteProject(ctx, 0, project.Dir, project.Containerization)
		return
	}

	okMessage := "Created"
	if projectOldDir != "" {
		deleteProject(ctx, 0, projectOldDir, project.Containerization)
		okMessage = "Updated"
	}

	okDelayedTask(ctx, delayedTask, project.Id, okMessage)
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

func DeleteProject(fctx *fiber.Ctx) error {
	ctx := context.Background()

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	user := service.GetCurrentUser(fctx)
	project := q.GetProjectShallow(ctx, id)
	if project == nil || project.UserId != user.Id {
		return fctx.Status(404).SendString("Project not found")
	}

	deleteProject(ctx, project.Id, project.Dir, project.Containerization)

	return fctx.SendStatus(fiber.StatusOK)
}

func deleteProject(ctx context.Context, id int, dir string, containerization string) {
	if id != 0 {
		q.DeleteProject(ctx, id)
	}
	cleanupContainers(dir, containerization)
	os.RemoveAll(filepath.Join(config.MediaPath, dir))
}

func cleanupContainers(dir string, containerization string) {
	switch containerization {
	case "docker":
		exec.Command("docker", "image", "remove", "-f", dir).Run()
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
		removeCmd.Dir = filepath.Join(config.MediaPath, dir)
		removeCmd.Run()
		exec.Command("docker", "system", "prune", "-f").Run()
	}
}