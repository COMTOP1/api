package error

import (
	"encoding/json"
	"errors"
	"fmt"
)

func New(text string) error {
	var error1 errorString
	err := json.Unmarshal([]byte(text), &error1)
	if err != nil {
		return errors.New(fmt.Sprintf("error message not compatible with json: %v", err))
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
