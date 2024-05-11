package filter

import (
	"fmt"
	"goto/src/service"
	u "goto/src/utils"
	sc "strconv"
	s "strings"

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
		{tf.Name != "", tf.Name, "task.name LIKE '%%' || $%d || '%%'"},
		{tf.Language != "", tf.Language, "project.language LIKE '%%' || $%d || '%%'"},
		{
			tf.Module != "",
			tf.Module,
			"project.id IN (SELECT project_id FROM module WHERE name LIKE '%%' || $%d || '%%')",
		},
		{tf.My, tf.UserId, "project.user_id = $%d"},
	}

	i := 1
	conditions := []string{}
	for _, fe := range filterEntries {
		if fe.IsActive {
			conditions = append(conditions, fmt.Sprintf(fe.QueryPart, i))
			tf.SqlArgs = append(tf.SqlArgs, fe.Value)
			i++
		}
	}

	if len(conditions) > 0 {
		tf.SqlCondition = " " + s.Join(conditions, " AND\n")
	} else {
		tf.SqlCondition = " TRUE"
	}

	return &tf
}