package gohetznerdns

import (
	"fmt"
	"strings"
)

func iif[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func validateNotNil[T any](parameterName string, value *T) error {
	var err *Error = nil
	if value == nil {
		err = &Error{Code: 900, Message: fmt.Sprintf("%s is nil", parameterName)}
	}
	if err != nil {
		return err.Error()
	}
	return nil
}

func validateNotEmpty(parameterName string, value *string) error {
	var err *Error = nil
	if _err := validateNotNil(parameterName, value); _err != nil {
		return _err
	}
	if len(strings.TrimSpace(*value)) == 0 {
		err = &Error{Code: 901, Message: fmt.Sprintf("%s is empty", parameterName)}
	}
	if err != nil {
		return err.Error()
	}
	return nil
}
