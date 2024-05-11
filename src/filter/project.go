package filter

import (
	"goto/src/service"

	"github.com/gofiber/fiber/v2"
)

type ProjectFilter struct {
	FilterBase
	UserId   int
	Name     string
	Language string
	Module   string
}

func NewProjectFilter(fctx *fiber.Ctx) *ProjectFilter {
	user := service.GetCurrentUser(fctx)

	pf := ProjectFilter{
		UserId:   user.Id,
		Name:     fctx.Query("name"),
		Language: fctx.Query("language"),
		Module:   fctx.Query("module"),
	}

	filterEntries := []FilterEntry{
		{pf.Name != "", pf.Name, "LOWER(project.name) LIKE LOWER('%%' || $%d || '%%')"},
		{pf.Language != "", pf.Language, "LOWER(project.language) LIKE LOWER('%%' || $%d || '%%')"},
		{
			pf.Module != "",
			pf.Module,
			"project.id IN (SELECT project_id FROM module WHERE LOWER(name) LIKE LOWER('%%' || $%d || '%%'))",
		},
		{true, pf.UserId, "project.user_id = $%d"},
	}

	pf.FilterBase = *NewFilter(&pf.FilterBase, filterEntries)
	return &pf
}