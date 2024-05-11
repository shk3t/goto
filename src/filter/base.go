package filter

type FilterEntry struct {
	IsActive  bool
	Value     any
	QueryPart string
}

type FilterBase struct {
	SqlCondition string
	SqlArgs      []any
}