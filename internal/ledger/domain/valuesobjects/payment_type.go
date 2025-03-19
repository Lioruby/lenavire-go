package valuesobjects

type PaymentType string

const (
	OneTime   PaymentType = "one-time"
	Recurring PaymentType = "recurring"
)
