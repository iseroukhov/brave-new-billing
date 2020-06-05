package payment

const (
	StatusExpected = iota + 1
	StatusPaid
	StatusError
)

type Status struct {
	ID   int64
	Name string
	Slug string
}
