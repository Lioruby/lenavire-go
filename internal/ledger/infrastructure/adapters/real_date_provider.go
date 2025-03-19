package adapters

import "time"

type RealDateProvider struct{}

func NewRealDateProvider() *RealDateProvider {
	return &RealDateProvider{}
}

func (p *RealDateProvider) Now() time.Time {
	return time.Now()
}
