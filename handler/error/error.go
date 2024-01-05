package error

import (
	"encoding/json"
	"fmt"
)

func New(text string) error {
	var error1 errorString
	err := json.Unmarshal([]byte(text), &error1)
	if err != nil {
		return fmt.Errorf("error message not compatible with json: %w", err)
	}
	return &error1
}

// errorString is a trivial implementation of error.
type errorString struct {
	S string `json:"error"`
}

func (e *errorString) Error() string {
	return e.S
}
