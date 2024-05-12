package filter

import (
	"goto/src/service"

	"github.com/gofiber/fiber/v2"
)

type ProjectFilter struct {
	FilterBase
	userId   int
	name     string
	language string
	module   string
}

func NewProjectFilter(fctx *fiber.Ctx) *ProjectFilter {
	user := service.GetCurrentUser(fctx)

	pf := ProjectFilter{
		userId:   user.Id,
		name:     fctx.Query("name"),
		language: fctx.Query("language"),
		module:   fctx.Query("module"),
	}

	filterEntries := []FilterEntry{
		{true, pf.userId, "project.user_id = $%d"},
		{pf.name != "", pf.name, "LOWER(project.name) LIKE LOWER('%%' || $%d || '%%')"},
		{pf.language != "", pf.language, "LOWER(project.language) LIKE LOWER('%%' || $%d || '%%')"},
		{
			pf.module != "",
			pf.module,
			"project.id IN (SELECT project_id FROM module WHERE LOWER(name) LIKE LOWER('%%' || $%d || '%%'))",
		},
	}

	pf.FilterBase = *NewFilter(&pf.FilterBase, filterEntries)
	return &pf
}