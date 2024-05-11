package filter

type ProjectFilter struct {
	Name     string
	Language string
	Module   string
}

// func (pf *ProjectFilter) SqlCondition() string {
// 	return `
//         project.name LIKE %$1% AND
//         project.language LIKE %$2% AND
//         project.id IN (SELECT project_id FROM module WHERE name LIKE %$3%)`
// }
//
// func (pf *ProjectFilter) SqlArgs() []any {
// 	return []any{pf.Name, pf.Language, pf.Module}
// }