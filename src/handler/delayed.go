package handler

import (
	"context"
	q "goto/src/database/query"
	m "goto/src/model"
	"goto/src/service"
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

type _ m.DelayedTask

// @tags Отложенные задания
// @summary Список моих отложенных заданий
// @security BearerAuth
// @param start query int false "Вернуть с"
// @param take query int false "Количество возвращаемых элементов"
// @produce json
// @success 200 {array} m.DelayedTasks
// @router /delayed-tasks [get]
func GetDelayedTasks(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)
	pager := service.NewPager(fctx)
	delayedTasks := q.GetUserDelayedTasks(ctx, user.Id, pager)
	return fctx.JSON(delayedTasks)
}

// @tags Отложенные задания
// @summary Детализация отложенного задания
// @security BearerAuth
// @param id path int true "Идентификатор отложенного задания"
// @produce json
// @success 200 {object} m.DelayedTask
// @router /delayed-tasks/{id} [get]
func GetDelayedTask(fctx *fiber.Ctx) error {
	ctx := context.Background()
	user := service.GetCurrentUser(fctx)

	id, err := sc.Atoi(fctx.Params("id"))
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString("Id is not correct")
	}

	delayedTask := q.GetUserDelayedTask(ctx, id, user.Id)
	if delayedTask == nil {
		return fctx.Status(404).SendString("Delayed task not found")
	}

	return fctx.JSON(delayedTask)
}