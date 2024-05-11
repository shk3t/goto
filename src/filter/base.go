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
	SqlCondition string
	SqlArgs      []any
}

func NewFilter(f *FilterBase, filterEntries []FilterEntry) *FilterBase {

	i := 1
	conditions := []string{}
	for _, fe := range filterEntries {
		if fe.IsActive {
			conditions = append(conditions, fmt.Sprintf(fe.QueryPart, i))
			f.SqlArgs = append(f.SqlArgs, fe.Value)
			i++
		}
	}

	if len(conditions) > 0 {
		f.SqlCondition = " " + s.Join(conditions, " AND\n")
	} else {
		f.SqlCondition = " TRUE"
	}

	return f
}