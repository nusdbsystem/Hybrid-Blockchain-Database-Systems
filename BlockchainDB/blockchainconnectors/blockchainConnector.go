package connectors

import "context"

type BlockchainConnector interface {
	Read(context.Context, string) (string, error)
	Write(context.Context, string, string) (string, error)
	Verify(context.Context, string, string, string) (bool, error)
}
