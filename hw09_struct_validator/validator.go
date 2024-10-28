package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	for _, e := range v {
		fmt.Fprintf(&b, "%s: %s", e.Field, e.Err.Error())
	}
	return b.String()
}

type structValidator struct {
	errors ValidationErrors
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	rValue := reflect.ValueOf(v)

	if rValue.Type().Kind() == reflect.Struct {
		validator := newValidator()
		return validator.ValidateStruct(v)
	}

	return nil
}

func newValidator() *structValidator {
	validator := structValidator{
		errors: make(ValidationErrors, 0),
	}

	return &validator
}

func (sv *structValidator) ValidateStruct(value interface{}) error {
	reflectValue := reflect.ValueOf(value)
	//reflectType := reflect.TypeOf(value)
	for i := 0; i < reflectValue.NumField(); i++ {
		fieldReflectValue := reflectValue.Field(i)
		//fieldReflectType := reflectType.Field(i)

		if fieldReflectValue.Type().Kind() == reflect.Struct {
			fieldStructValidator := newValidator()
			fmt.Println("struct field: ", fieldReflectValue.Elem())
			fieldErr := fieldStructValidator.ValidateStruct(fieldReflectValue)

			if fieldErr == nil {
				continue
			}
			var validationErrors ValidationErrors

			if errors.As(fieldErr, &validationErrors) {
				sv.errors = append(sv.errors, validationErrors...)
			} else {
				return fieldErr
			}
		}

	}

	if len(sv.errors) > 0 {
		return sv.errors
	}
	return nil
}
