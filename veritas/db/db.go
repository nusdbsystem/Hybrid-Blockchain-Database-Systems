package db

type DB interface {
	Get(key string) (value string, err error)
	Set(key, value string) error
}
