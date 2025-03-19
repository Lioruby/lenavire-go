package ports

type LedgerActivityChannel interface {
	Send(message string) error
}
