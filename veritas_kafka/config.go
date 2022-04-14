package veritas

type Config struct {
	Signature string
	Topic     string
	Parties   map[string]struct{}
	BlockSize int
}
