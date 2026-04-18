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

		for _, fe := range errs {
			var msg string
			switch fe.Tag() {
			case "required":
				msg = "This field is mandatory"
			case "email":
				msg = "Invalid email format"
			case "min":
				msg = "Minimum " + fe.Param() + " character"
			default:
				msg = "An error occurred in this field"
			}

			field, _ := t.Elem().FieldByName(fe.Field())
			key := field.Tag.Get("json")

			if key == "" || key == "-" {
				key = strings.ToLower(fe.Field())
			}
			errorMessages[key] = msg
		}
	}
	return errorMessages
}
