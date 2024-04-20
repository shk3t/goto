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
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func LoadProject(c fiber.Ctx) error {
	file, err := c.FormFile("project")
	if err != nil {
		c.Status(400).SendString(err.Error())
	}

	postfix := uuid.New().String()
	nameLess, extension := utils.SplitExt(file.Filename)
	projectName := nameLess + "_" + postfix

	archivePath := filepath.Join(config.MediaPath, projectName+extension)
	if err := c.SaveFile(file, archivePath); err != nil {
		c.Status(400).SendString(err.Error())
	}

	go func() {
		if err := utils.Unzip(archivePath, true); err != nil {
			log.Println(err)
			return
		}
		if err := os.Remove(archivePath); err != nil {
			log.Println(err)
			return
		}

		projectPath := filepath.Join(config.MediaPath, projectName)
		gotoConfigPath := filepath.Join(projectPath, config.GotoConfigName)
		gotoConfig, err := model.LoadGotoConfig(gotoConfigPath)
		if err != nil {
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
			log.Println(err)
			return
		}

		project := model.NewProjectFromConfig(gotoConfig)
        project.Dir = projectName
		query.CreateProject(context.Background(), project)
	}()

	return c.SendString("OK")
}