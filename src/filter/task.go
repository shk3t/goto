package filter

import (
	"goto/src/service"
	u "goto/src/utils"
	sc "strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskFilter struct {
	FilterBase
	UserId   int
	My       bool
	Name     string
	Language string
	Module   string
}

func NewTaskFilter(fctx *fiber.Ctx) *TaskFilter {
	user := service.GetCurrentUser(fctx)

	tf := TaskFilter{
		UserId:   user.Id,
		My:       u.Default(sc.ParseBool(fctx.Query("my"))),
		Name:     fctx.Query("name"),
		Language: fctx.Query("language"),
		Module:   fctx.Query("module"),
	}

	filterEntries := []FilterEntry{
		{tf.Name != "", tf.Name, "LOWER(task.name) LIKE LOWER('%%' || $%d || '%%')"},
		{tf.Language != "", tf.Language, "LOWER(project.language) LIKE LOWER('%%' || $%d || '%%')"},
		{
			tf.Module != "",
			tf.Module,
			"project.id IN (SELECT project_id FROM module WHERE LOWER(name) LIKE LOWER('%%' || $%d || '%%'))",
		},
		{tf.My, tf.UserId, "project.user_id = $%d"},
	}

	tf.FilterBase = *NewFilter(&tf.FilterBase, filterEntries)
	return &tf
}