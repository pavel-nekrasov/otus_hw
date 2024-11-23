package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	ValidationError struct {
		Field string
		Err   error
	}

	ValidationErrors []ValidationError

	structValidator struct {
		prefix           string
		criticalError    error
		validationErrors ValidationErrors
	}

	basicTypes interface {
		~int | ~string
	}

	integerValidator func(f string, v int, criteria string) error
	stringValidator  func(f string, v string, criteria string) error
)

var intRules = map[string]integerValidator{
	"in":  validateIntInclusion,
	"max": validateIntMax,
	"min": validateIntMin,
}

var strRules = map[string]stringValidator{
	"in":     validateStringInclusion,
	"len":    validateStringLen,
	"regexp": validateStringRegex,
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err.Error())
}

func (v ValidationErrors) Error() string {
	var b strings.Builder
	for i, e := range v {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	reflectValue := reflect.ValueOf(v)

	// разыменовать указатель (если передавалcя указатель)
	if reflectValue.Type().Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.Type().Kind() == reflect.Struct {
		validator := newStructValidator("")
		return validator.ValidateStruct(reflectValue.Interface())
	}

	return nil
}

func newStructValidator(prefix string) *structValidator {
	validator := structValidator{
		prefix:           prefix,
		validationErrors: make(ValidationErrors, 0),
	}

	return &validator
}

func (sv *structValidator) ValidateStruct(value interface{}) error {
	structValue := reflect.ValueOf(value)
	structType := reflect.TypeOf(value)

	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		tag := structField.Tag.Get("validate")

		if tag == "" {
			continue
		}

		//exhaustive:ignore
		switch structField.Type.Kind() {
		case reflect.Int, reflect.String:
			sv.handleScalar(structField, tag, structValue.Field(i))
		case reflect.Slice:
			sv.handleVector(structField, tag, structValue.Field(i))
		case reflect.Struct:
			sv.handleStruct(structField, tag, structValue.Field(i))
		case reflect.Ptr:
			sv.handlePtr(structField, tag, structValue.Field(i))
		default:
		}
		if sv.criticalError != nil {
			return sv.criticalError
		}
	}

	if len(sv.validationErrors) > 0 {
		return sv.validationErrors
	}
	return nil
}

func (sv *structValidator) fullFieldName(fieldName string) string {
	if sv.prefix == "" {
		return fieldName
	}
	return fmt.Sprintf("%v.%v", sv.prefix, fieldName)
}

func (sv *structValidator) handleScalar(field reflect.StructField, tag string, value reflect.Value) {
	fullFieldName := sv.fullFieldName(field.Name)

	//exhaustive:ignore
	switch field.Type.Kind() {
	case reflect.Int:
		sv.filterErr(validateTags(fullFieldName, tag, int(value.Int()), intRules))
	case reflect.String:
		sv.filterErr(validateTags(fullFieldName, tag, value.String(), strRules))
	}
}

func (sv *structValidator) handleVector(field reflect.StructField, tag string, value reflect.Value) {
	fullFieldName := sv.fullFieldName(field.Name)

	//exhaustive:ignore
	switch field.Type.Elem().Kind() {
	case reflect.Int:
		for _, v := range value.Interface().([]int) {
			sv.filterErr(validateTags(fullFieldName, tag, v, intRules))
			if sv.criticalError != nil {
				return
			}
		}
	case reflect.String:
		for _, v := range value.Interface().([]string) {
			sv.filterErr(validateTags(fullFieldName, tag, v, strRules))
			if sv.criticalError != nil {
				return
			}
		}
	}
}

func (sv *structValidator) handleStruct(field reflect.StructField, tag string, value reflect.Value) {
	if tag != "nested" {
		return
	}
	fullFieldName := sv.fullFieldName(field.Name)
	innerValidator := newStructValidator(fullFieldName)
	sv.filterErr(innerValidator.ValidateStruct(value.Interface()))
}

func (sv *structValidator) handlePtr(field reflect.StructField, tag string, value reflect.Value) {
	if tag != "nested" {
		return
	}
	if field.Type.Elem().Kind() == reflect.Struct {
		fullFieldName := sv.fullFieldName(field.Name)
		innerValidator := newStructValidator(fullFieldName)
		sv.filterErr(innerValidator.ValidateStruct(value.Elem().Interface()))
	}
}

func (sv *structValidator) filterErr(err error) {
	if err == nil {
		return
	}

	var validationErrors ValidationErrors
	if errors.As(err, &validationErrors) {
		sv.validationErrors = append(sv.validationErrors, validationErrors...)
		return
	}
	sv.criticalError = err
}

func validateTags[T basicTypes, H ~func(string, T, string) error](field string,
	tag string,
	val T,
	validators map[string]H,
) error {
	validationErrors := make(ValidationErrors, 0)
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		err := validateTag(field, rule, val, validators)
		if err == nil {
			continue
		}
		var validationErr ValidationError
		if errors.As(err, &validationErr) {
			validationErrors = append(validationErrors, validationErr)
			continue
		}
		return err
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateTag[T basicTypes, H ~func(string, T, string) error](field string,
	rule string,
	val T,
	validators map[string]H,
) error {
	ruleParts := strings.Split(rule, ":")
	if len(ruleParts) != 2 {
		return fmt.Errorf("%v: wrong notation: bad format %v", field, rule)
	}
	validator, ok := validators[ruleParts[0]]
	if !ok {
		return fmt.Errorf("%v: wrong notation: unknown rule \"%v\"", field, ruleParts[0])
	}
	return validator(field, val, ruleParts[1])
}

func validateStringLen(field string, value string, criteria string) error {
	intLen, err := strconv.Atoi(criteria)
	if err != nil {
		return fmt.Errorf("%v: wrong notation: len \"%v\" is not an integer", field, criteria)
	}

	if len(value) != intLen {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("wrong len %v, expected %v", len(value), intLen),
		}
	}

	return nil
}

func validateStringInclusion(field string, value string, criteria string) error {
	expectedValues := strings.Split(criteria, ",")

	if !sliceContains(expectedValues, value) {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("wrong value \"%v\", should be one of %v", value, expectedValues),
		}
	}

	return nil
}

func validateStringRegex(field string, value string, criteria string) error {
	r, err := regexp.Compile(criteria)
	if err != nil {
		return fmt.Errorf("%v: wrong notation: regex \"%v\" is incorrect", field, criteria)
	}

	if !r.Match([]byte(value)) {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("wrong value \"%v\"", value),
		}
	}

	return nil
}

func validateIntInclusion(field string, value int, criteria string) error {
	expectedValuesS := strings.Split(criteria, ",")

	expectedValues := make([]int, 0)
	for _, s := range expectedValuesS {
		i, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("%v: wrong notation: \"%v\" is not an integer", field, s)
		}
		expectedValues = append(expectedValues, i)
	}

	if !sliceContains(expectedValues, value) {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("wrong value %v, should be one of %v", value, expectedValues),
		}
	}

	return nil
}

func validateIntMax(field string, value int, criteria string) error {
	maxCriteria, err := strconv.Atoi(criteria)
	if err != nil {
		return fmt.Errorf("%v: wrong notation: max \"%v\" is not an integer", field, criteria)
	}

	if value > maxCriteria {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("wrong value %v, should less or equal to %v", value, criteria),
		}
	}

	return nil
}

func validateIntMin(field string, value int, criteria string) error {
	minCriteria, err := strconv.Atoi(criteria)
	if err != nil {
		return fmt.Errorf("%v: wrong notation: min \"%v\" is not an integer", field, criteria)
	}

	if value < minCriteria {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("wrong value %v, should greater or equal to %v", value, criteria),
		}
	}

	return nil
}

func sliceContains[T basicTypes](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
