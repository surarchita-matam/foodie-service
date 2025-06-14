package utils

import (
	"errors"
	"reflect"
)

var validate = &Validator{}

type Validator struct{}

func (v *Validator) Struct(data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Struct {
		return errors.New("validation only works on structs")
	}
	return nil
}

func Validate(data interface{}) error {
	return validate.Struct(data)
}