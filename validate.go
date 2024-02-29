package hetznerdns

import (
	"fmt"
	"strings"
)

type argError struct {
	arg    string
	reason string
}

func newArgError(arg, reason string) *argError {
	return &argError{
		arg:    arg,
		reason: reason,
	}
}

func (e *argError) Error() string {
	return fmt.Sprintf("%s is invalid because %s", e.arg, e.reason)
}

func validateNotEmpty(name, value string) error {
	if len(strings.TrimSpace(value)) == 0 {
		return newArgError(name, "cannot be empty")
	}
	return nil
}
