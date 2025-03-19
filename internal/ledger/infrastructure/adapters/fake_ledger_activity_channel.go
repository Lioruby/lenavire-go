package adapters

type FakeLedgerActivityChannel struct {
	Messages []string
}

func NewFakeLedgerActivityChannel() *FakeLedgerActivityChannel {
	return &FakeLedgerActivityChannel{
		Messages: make([]string, 0),
	}
}

func (c *FakeLedgerActivityChannel) Send(message string) error {
	c.Messages = append(c.Messages, message)
	return nil
}
