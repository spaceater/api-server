package utility

import (
	"encoding/json"
	"fmt"
)

type MessageData struct {
	Type string      `json:"-"`
	Data interface{} `json:"-"`
}

type ErrorData struct {
	Error string `json:"error"`
}

func (md *MessageData) ToBytes() ([]byte, error) {
	jsonData, err := json.Marshal(md.Data)
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil, "%s:%s", md.Type, string(jsonData)), nil
}
