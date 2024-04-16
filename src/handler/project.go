package handler

import (
	// "goto/src/model"
	"goto/src/config"
	"goto/src/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
)

func LoadProject(c fiber.Ctx) error {
	// project := new(model.Project)

	file, err := c.FormFile("project")
	if err != nil {
		c.Status(400).SendString(err.Error())
	}

	archivePath := filepath.Join(config.MediaPath, file.Filename)
	if err := c.SaveFile(file, archivePath); err != nil {
		c.Status(400).SendString(err.Error())
	}

	go func() {
		if err := utils.Unzip(archivePath); err != nil {
			log.Println(err)
		}
		if err := os.Remove(archivePath); err != nil {
			log.Println(err)
		}
	}()

	return c.SendString("OK")
}