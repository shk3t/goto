package filter

import (
	"goto/src/service"
	u "goto/src/utils"
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskFilter struct {
	FilterBase
	userId   int
	my       bool
	name     string
	language string
	module   string
}

func NewTaskFilter(fctx *fiber.Ctx) *TaskFilter {
	user := service.GetCurrentUser(fctx)

	tf := TaskFilter{
		userId:   user.Id,
		my:       u.Default(sc.ParseBool(fctx.Query("my"))),
		name:     fctx.Query("name"),
		language: fctx.Query("language"),
		module:   fctx.Query("module"),
	}

	filterEntries := []FilterEntry{
		{tf.my, tf.userId, "project.user_id = $%d"},
		{tf.name != "", tf.name, "LOWER(task.name) LIKE LOWER('%%' || $%d || '%%')"},
		{tf.language != "", tf.language, "LOWER(project.language) LIKE LOWER('%%' || $%d || '%%')"},
		{
			tf.module != "",
			tf.module,
			"project.id IN (SELECT project_id FROM module WHERE LOWER(name) LIKE LOWER('%%' || $%d || '%%'))",
		},
	}

	tf.FilterBase = *NewFilter(&tf.FilterBase, filterEntries)
	return &tf
}