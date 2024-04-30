package handler

import (
	"context"
	"errors"
	"goto/src/config"
	"goto/src/database/query"
	"goto/src/model"
	"goto/src/utils"
	"os"
	"strconv"
	"strings"

	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func postCreateProject(c *fiber.Ctx, projectName string) {
	ctx := context.Background()

	projectPath := filepath.Join(config.MediaPath, projectName)
	gotoConfigPath := filepath.Join(projectPath, config.GotoConfigName)
	gotoConfig, err := model.LoadGotoConfig(gotoConfigPath)
	if err != nil {
		os.RemoveAll(projectPath)
		return
	}

	project := model.NewProjectFromConfig(gotoConfig)
	project.Dir = projectName
	project.User = *GetCurrentUser(c)
	if err = query.CreateProject(ctx, project); err != nil {
		os.RemoveAll(projectPath)
		return
	}

	var buildCmd *exec.Cmd
	switch gotoConfig.Containerization {
	case "docker":
		buildCmd = exec.Command("docker", "buildx", "build", "-t", projectName, ".")
	case "docker-compose":
		buildCmd = exec.Command("docker", "compose", "build")
	default:
		err = errors.New("Specified containerization type is not implemented")
		return
	}

	buildCmd.Dir = projectPath
	err = buildCmd.Run()
	if err != nil {
		query.DeleteProject(ctx, project.Id)
		os.RemoveAll(projectPath)
		return
	}
}

func postCreateProjectZip(c *fiber.Ctx, projectName string, archivePath string) {
	if err := utils.Unzip(archivePath, true); err != nil {
		return
	}
	if err := os.Remove(archivePath); err != nil {
		return
	}

	postCreateProject(c, projectName)
}

func LoadProject(c *fiber.Ctx) error {
	body := struct{ Url string }{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	postfix := uuid.New().String()

	if body.Url != "" {
		urlParts := strings.Split(body.Url, "/")
		repoName := urlParts[len(urlParts)-1]
		projectName := repoName + "_" + postfix

		gitCloneCmd := exec.Command("git", "clone", body.Url, projectName)
		gitCloneCmd.Dir = config.MediaPath
		if err := gitCloneCmd.Run(); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid url")
		}

		go postCreateProject(c, projectName)

	} else {
		file, err := c.FormFile("project")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Use `project` as a key for uploaded file")
		}

		nameLess, extension := utils.SplitExt(file.Filename)
		projectName := nameLess + "_" + postfix

		archivePath := filepath.Join(config.MediaPath, projectName+extension)
		if err = c.SaveFile(file, archivePath); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		go postCreateProjectZip(c, projectName, archivePath)
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteProject(c *fiber.Ctx) error {
	ctx := context.Background()

	projectId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

    user := GetCurrentUser(c)
	project, err := query.GetUserProject(ctx, projectId, user.Id)
	if err != nil {
		return c.Status(404).SendString("Project not found")
	}

	query.DeleteProject(ctx, projectId)

	projectPath := filepath.Join(config.MediaPath, project.Dir)
	os.RemoveAll(projectPath)

	switch project.Containerization {
	case "docker":
		exec.Command("docker", "image", "remove", "-f", project.Dir).Run()
		exec.Command("docker", "system", "prune", "-f").Run()
	case "docker-compose":
		removeCmd := exec.Command("docker", "compose", "remove", "-fsv")
		removeCmd.Dir = config.MediaPath
		removeCmd.Run()
		exec.Command("docker", "system", "prune", "-f").Run()
	}

	return c.SendStatus(fiber.StatusOK)
}