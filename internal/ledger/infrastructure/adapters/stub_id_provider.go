package adapters

type StubIdProvider struct {
	Id string
}

func NewStubIdProvider(id string) *StubIdProvider {
	return &StubIdProvider{
		Id: id,
	}
}

func (p *StubIdProvider) Generate() string {
	return p.Id
}
