package handler

import (
	"context"
	q "goto/src/database/query"
	f "goto/src/filter"
	m "goto/src/model"
	"goto/src/service"
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

type _ m.TasksMin

// @tags Задачи
// @summary Список задач
// @security BearerAuth
// @param start query int false "Вернуть с"
// @param take query int false "Количество возвращаемых элементов"
// @param my query bool false "Созданные мной"
// @param name query string false "Название"
// @param language query string false "Язык"
// @param module query string false "Название модуля"
// @produce json
// @success 200 {array} m.TasksMin
// @router /tasks [get]
func GetTasks(fctx *fiber.Ctx) error {
	ctx := context.Background()
	pager := service.NewPager(fctx)
	filter := f.NewTaskFilter(fctx)
	tasks := q.GetTasks(ctx, pager, filter)
	return fctx.JSON(tasks.Min())
}

// @tags Задачи
// @summary Детализация задачи
// @security BearerAuth
// @param id path int true "Идентификатор задачи"
// @produce json
// @success 200 {object} m.TaskPrivate
// @router /tasks/{id} [get]
func GetTask(fctx *fiber.Ctx) error {
	ctx := context.Background()

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	task := q.GetTask(ctx, id)
	if task == nil {
		return fctx.Status(404).SendString("Task not found")
	}

	return fctx.JSON(task.Private())
}