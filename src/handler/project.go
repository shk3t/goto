package handler

import (
	"context"
	"errors"
	"goto/src/config"
	"goto/src/database/query"
	"goto/src/model"
	"goto/src/utils"
	"log"
	"os"
	"strconv"
	"strings"

	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func PostCreateProject(projectName string) {
	ctx := context.Background()

	projectPath := filepath.Join(config.MediaPath, projectName)
	gotoConfigPath := filepath.Join(projectPath, config.GotoConfigName)
	gotoConfig, err := model.LoadGotoConfig(gotoConfigPath)
	if err != nil {
		os.RemoveAll(projectPath)
		log.Println(err)
		return
	}

	project := model.NewProjectFromConfig(gotoConfig)
	project.Dir = projectName
	if err = query.CreateProject(ctx, project); err != nil {
		os.RemoveAll(projectPath)
		log.Println(err)
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
		log.Println(err)
		return
	}

	buildCmd.Dir = projectPath
	err = buildCmd.Run()
	if err != nil {
		query.DeleteProject(ctx, project.Id)
		os.RemoveAll(projectPath)
		log.Println(err)
		return
	}
}

func PostCreateProjectZip(projectName string, archivePath string) {
	if err := utils.Unzip(archivePath, true); err != nil {
		log.Println(err)
		return
	}
	if err := os.Remove(archivePath); err != nil {
		log.Println(err)
		return
	}

	PostCreateProject(projectName)
}

func LoadProject(c fiber.Ctx) error {
	body := struct{ Url string }{}
	c.Bind().Body(&body)
	url := body.Url

	postfix := uuid.New().String()

	if url != "" {
		urlParts := strings.Split(url, "/")
		repoName := urlParts[len(urlParts)-1]
		projectName := repoName + "_" + postfix

		gitCloneCmd := exec.Command("git", "clone", url, projectName)
		gitCloneCmd.Dir = config.MediaPath
		if err := gitCloneCmd.Run(); err != nil {
			return c.Status(400).SendString("Invalid url")
		}

		go PostCreateProject(projectName)

	} else {
		file, err := c.FormFile("project")
		if err != nil {
			return c.Status(400).SendString("Use `project` as a key for uploaded file")
		}

		nameLess, extension := utils.SplitExt(file.Filename)
		projectName := nameLess + "_" + postfix

		archivePath := filepath.Join(config.MediaPath, projectName+extension)
		if err = c.SaveFile(file, archivePath); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		go PostCreateProjectZip(projectName, archivePath)
	}

	return c.SendString("OK")
}

func DeleteProject(c fiber.Ctx) error {
	ctx := context.Background()

	projectId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Id is not correct")
	}

	project, err := query.GetProject(ctx, projectId)
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

	return c.SendString("OK")
}