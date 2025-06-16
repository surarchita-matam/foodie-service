package utils

import (
	"errors"
	"reflect"
	"regexp"
)

var validate = &Validator{}

type Validator struct{}

func (v *Validator) Struct(data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("validation only works on structs")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get validation tags
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Check required tag
		if validateTag == "required" || validateTag == "required,email" || validateTag == "required,password" {
			switch field.Kind() {
			case reflect.String:
				if field.String() == "" {
					return errors.New(fieldType.Name + " is required")
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if field.Int() <= 0 {
					return errors.New(fieldType.Name + " must be greater than 0")
				}
			case reflect.Float32, reflect.Float64:
				if field.Float() <= 0 {
					return errors.New(fieldType.Name + " must be greater than 0")
				}
			case reflect.Slice:
				if field.IsNil() || field.Len() == 0 {
					return errors.New(fieldType.Name + " cannot be empty")
				}
				// Validate each item in the slice
				for j := 0; j < field.Len(); j++ {
					if err := v.Struct(field.Index(j).Interface()); err != nil {
						return err
					}
				}
			}
		}

		// Check email validation
		if validateTag == "required,email" || validateTag == "email" {
			if field.Kind() == reflect.String {
				emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
				if !emailRegex.MatchString(field.String()) {
					return errors.New(fieldType.Name + " must be a valid email address")
				}
			}
		}
	}
	return nil
}

func Validate(data interface{}) error {
	return validate.Struct(data)
}
