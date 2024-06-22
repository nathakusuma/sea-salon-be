package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() Validator {
	return Validator{
		validator: validator.New(),
	}
}

func getJSONFieldName(structObj any, fieldName string) (string, error) {
	t := reflect.TypeOf(structObj)
	field, found := t.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field not found")
	}

	return field.Tag.Get("json"), nil
}

func (v Validator) Validate(data any) []string {
	errs := v.validator.Struct(data)
	if errs != nil {
		errMsgs := make([]string, 0)
		for _, err := range errs.(validator.ValidationErrors) {
			fieldName, err2 := getJSONFieldName(data, err.Field())
			if err2 != nil {
				fieldName = err.Field()
			}
			errMsg := fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'",
				fieldName, err.Value(), err.Tag())
			errMsgs = append(errMsgs, errMsg)
		}
		return errMsgs
	}

	return nil
}
