package adapters

import "time"

type StubDateProvider struct {
	Date string
}

func NewStubDateProvider(date string) *StubDateProvider {
	return &StubDateProvider{
		Date: date,
	}
}

func (p *StubDateProvider) Now() time.Time {
	date, _ := time.Parse("2006-01-02", p.Date)
	return date
}
