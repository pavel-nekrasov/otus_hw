package customerrors

import "fmt"

type (
	ParamError struct {
		Param string
		Err   error
	}

	ValidationError struct {
		Field string
		Err   error
	}

	NotFound struct {
		Message string
	}
)

func (v ParamError) Error() string {
	return fmt.Sprintf("%s: %s", v.Param, v.Err.Error())
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err.Error())
}

func (e NotFound) Error() string {
	return e.Message
}
