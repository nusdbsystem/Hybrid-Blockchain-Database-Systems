package veritastm

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Value struct {
	Val     string
	Version int64
}

func Encode(val string, ver int64) (string, error) {
	v, err := json.Marshal(&Value{
		Val:     val,
		Version: ver,
	})
	return string(v), err
}

func Decode(entry string) (*Value, error) {
	var v Value
	if err := json.Unmarshal([]byte(entry), &v); err != nil {
		return nil, err
	}
	return &v, nil
}
