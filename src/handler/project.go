package handler

import (
	// "goto/src/model"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v3"
)

func LoadProject(c fiber.Ctx) error {
	// project := new(model.Project)

	file, err := c.FormFile("project")
	if err != nil {
		c.Status(400).SendString(err.Error())
	}

	err = c.SaveFile(file, filepath.Join("media", file.Filename))
	if err != nil {
		c.Status(400).SendString(err.Error())
	}

	go func() {
		timer := time.NewTimer(10 * time.Second)
		<-timer.C
		fmt.Println("BG TASK IS DONE")
	}()

	// if err := c.Bind().Body(project); err != nil {
	// 	return c.Status(500).JSON(err)
	// }

	return c.SendString("OK")
}