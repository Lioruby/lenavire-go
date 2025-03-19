package ports

type IdProvider interface {
	Generate() string
}
