package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error, obj interface{}) map[string]string {
	errorMessages := make(map[string]string)

	if errs, ok := err.(validator.ValidationErrors); ok {

		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		for _, fe := range errs {
			var msg string
			switch fe.Tag() {
			case "required":
				msg = "this field is required"
			case "email":
				msg = "invalid email format"
			case "min":
				msg = "minimum " + fe.Param() + " characters"
			default:
				msg = "invalid value"
			}

			field, found := t.FieldByName(fe.Field())
			if !found {
				continue
			}
			key := strings.Split(field.Tag.Get("json"), ",")[0]

			if key == "" || key == "-" {
				key = strings.ToLower(fe.Field())
			}
			errorMessages[key] = msg
		}
	}
	return errorMessages
}
