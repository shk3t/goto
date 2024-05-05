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

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetProjects(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	projects := query.GetUserProjects(ctx, user.Id)

	response := []model.ProjectPublic{}
	for _, p := range projects {
		response = append(response, model.ProjectPublic{
			ProjectBase: p.ProjectBase,
			Id:          p.Id,
			Tasks:       p.Tasks,
		})
	}

	return c.JSON(response)
}

func GetProject(c *fiber.Ctx) error {
	ctx := context.Background()
	user := GetCurrentUser(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	project := query.GetUserProject(ctx, id, user.Id)
	if project == nil {
		return c.Status(404).SendString("Project not found")
	}

	return c.JSON(model.ProjectPublic{
		ProjectBase: project.ProjectBase,
		Id:          project.Id,
		Tasks:       project.Tasks,
	})
}

func LoadProject(c *fiber.Ctx) error {
	user := GetCurrentUser(c)
	body := struct{ Url string }{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
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

		go postCreateProject(user, projectName)

	} else {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Use `file` as a key for uploaded file")
		}

		nameLess, extension := utils.SplitExt(file.Filename)
		projectName := nameLess + "_" + postfix

		archivePath := filepath.Join(config.MediaPath, projectName+extension)
		if err = c.SaveFile(file, archivePath); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		go postCreateProjectZip(user, projectName, archivePath)
	}

	return c.SendStatus(fiber.StatusOK)
}

func postCreateProject(user *model.User, projectName string) {
	ctx := context.Background()

	projectPath := filepath.Join(config.MediaPath, projectName)
	gotoConfigPath := filepath.Join(projectPath, config.GotoConfigName)
	gotoConfig, err := model.LoadGotoConfig(gotoConfigPath)
	if err != nil {
		log.Println(err)
		os.RemoveAll(projectPath)
		return
	}

	project := model.NewProjectFromConfig(gotoConfig)
	project.Dir = projectName
	project.User = *user
	if err = query.CreateProject(ctx, project); err != nil {
		log.Println(err)
		os.RemoveAll(projectPath)
		return
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
		err = errors.New("Specified containerization type is not implemented")
		return
	}

	for _, cmd := range prepareCmds {
		cmd.Dir = projectPath
		if err = cmd.Run(); err != nil {
			log.Println(err)
			query.DeleteProject(ctx, project.Id)
			os.RemoveAll(projectPath)
			return
		}
	}
}

func postCreateProjectZip(user *model.User, projectName string, archivePath string) {
	if err := utils.Unzip(archivePath, true); err != nil {
		log.Println(err)
		return
	}
	if err := os.Remove(archivePath); err != nil {
		log.Println(err)
		return
	}

	postCreateProject(user, projectName)
}

func UpdateProject(c *fiber.Ctx) error {
	return errors.New("TODO")
}

func DeleteProject(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	user := GetCurrentUser(c)
	project := query.GetUserProject(ctx, id, user.Id)
	if project == nil {
		return c.Status(404).SendString("Project not found")
	}
	projectPath := filepath.Join(config.MediaPath, project.Dir)

	query.DeleteProject(ctx, id)

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
		removeCmd.Dir = projectPath
		removeCmd.Run()
		exec.Command("docker", "system", "prune", "-f").Run()
	}

	os.RemoveAll(projectPath)

	return c.SendStatus(fiber.StatusOK)
}