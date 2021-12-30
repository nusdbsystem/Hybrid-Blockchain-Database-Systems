package connectors

import "context"

type Testconnector struct {
}

func (t *Testconnector) Read(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (t *Testconnector) Write(ctx context.Context, key, value string) (string, error) {
	return "", nil
}

func (t *Testconnector) Verify(ctx context.Context, opt, key, tx string) (bool, error) {
	return true, nil
}
