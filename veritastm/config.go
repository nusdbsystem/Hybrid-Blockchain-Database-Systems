package veritastm

type Config struct {
	Signature   string
	Parties     map[string]struct{}
	BlockSize   int
	LedgerPath  string
	ABCIRPCAddr string
}
