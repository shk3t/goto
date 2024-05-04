package query

type Scanable interface {
	Scan(dest ...any) error
}
