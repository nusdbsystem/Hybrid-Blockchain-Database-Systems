package raftkv

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func NewSetCommand(key, value string) *command {
	return &command{
		Op:    "set",
		Key:   key,
		Value: value,
	}
}
