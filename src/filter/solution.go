package filter

import (
	"goto/src/service"
	u "goto/src/utils"
	sc "strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SolutionFilter struct {
	FilterBase
	userId   int
	taskId   int
	dateFrom time.Time
	dateTo   time.Time
	status   string
	name     string
	language string
	module   string
}

func NewSolutionFilter(fctx *fiber.Ctx) *SolutionFilter {
	user := service.GetCurrentUser(fctx)

	sf := SolutionFilter{
		userId:   user.Id,
		taskId:   u.Default(sc.Atoi(fctx.Query("taskId"))),
		dateFrom: u.Default(time.Parse(time.DateTime, fctx.Query("dateFrom"))),
		dateTo:   u.Default(time.Parse(time.DateTime, fctx.Query("dateTo"))),
		status:   fctx.Query("status"),
		name:     fctx.Query("name"),
		language: fctx.Query("language"),
		module:   fctx.Query("module"),
	}

	filterEntries := []FilterEntry{
		{true, sf.userId, "solution.user_id = $%d"},
		{sf.taskId != 0, sf.taskId, "solution.task_id = $%d"},
		{!sf.dateFrom.IsZero(), sf.dateFrom, "solution.updated_at >= $%d"},
		{!sf.dateTo.IsZero(), sf.dateTo, "solution.updated_at <= $%d"},
		{sf.status != "", sf.status, "LOWER(solution.status) = LOWER($%d)"},
		{sf.name != "", sf.name, "LOWER(task.name) LIKE LOWER('%%' || $%d || '%%')"},
		{sf.language != "", sf.language, "LOWER(project.language) LIKE LOWER('%%' || $%d || '%%')"},
		{
			sf.module != "",
			sf.module,
			"project.id IN (SELECT project_id FROM module WHERE LOWER(name) LIKE LOWER('%%' || $%d || '%%'))",
		},
	}

	sf.FilterBase = *NewFilter(&sf.FilterBase, filterEntries)
	return &sf
}