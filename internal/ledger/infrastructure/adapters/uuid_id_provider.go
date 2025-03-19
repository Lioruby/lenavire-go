package adapters

import "github.com/google/uuid"

type UUIDIdProvider struct{}

func NewUUIDIdProvider() *UUIDIdProvider {
	return &UUIDIdProvider{}
}

func (p *UUIDIdProvider) Generate() string {
	return uuid.New().String()
}
