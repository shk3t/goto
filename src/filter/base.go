package filter

import (
	"fmt"
	s "strings"
)

type FilterEntry struct {
	IsActive  bool
	Value     any
	QueryPart string
}

type FilterBase struct {
	QueryCondition string
	QueryArgs      []any
}

func NewFilter(f *FilterBase, filterEntries []FilterEntry) *FilterBase {

	i := 1
	conditions := []string{}
	for _, fe := range filterEntries {
		if fe.IsActive && fe.Value != nil {
			conditions = append(conditions, fmt.Sprintf(fe.QueryPart, i))
			f.QueryArgs = append(f.QueryArgs, fe.Value)
			i++
		} else if fe.IsActive && fe.Value == nil {
			conditions = append(conditions, fe.QueryPart)
		}
	}

	if len(conditions) > 0 {
		f.QueryCondition = " " + s.Join(conditions, " AND\n") + " "
	} else {
		f.QueryCondition = " TRUE "
	}

	return f
}