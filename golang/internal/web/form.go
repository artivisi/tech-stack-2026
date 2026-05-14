package web

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RegistrationForm struct {
	Email    string `form:"email"    validate:"required,min=3,max=254,email_format"`
	FullName string `form:"fullName" validate:"required,min=2,max=100,name_format"`
	Phone    string `form:"phone"    validate:"required,min=7,max=20,phone_format"`
}

var (
	emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	nameRe  = regexp.MustCompile(`^[\p{L}\p{M}\s.'\-]+$`)
	phoneRe = regexp.MustCompile(`^[+0-9 ()\-]+$`)
)

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	v.RegisterValidation("email_format", func(fl validator.FieldLevel) bool {
		return emailRe.MatchString(fl.Field().String())
	})
	v.RegisterValidation("name_format", func(fl validator.FieldLevel) bool {
		return nameRe.MatchString(fl.Field().String())
	})
	v.RegisterValidation("phone_format", func(fl validator.FieldLevel) bool {
		return phoneRe.MatchString(fl.Field().String())
	})
	return v
}

func CollectErrors(err error) map[string]string {
	errors := map[string]string{}
	var vErrs validator.ValidationErrors
	if !asValidationErrors(err, &vErrs) {
		return errors
	}
	for _, e := range vErrs {
		field := e.Field()
		if _, exists := errors[field]; exists {
			continue
		}
		errors[field] = messageFor(field, e.Tag())
	}
	return errors
}

func asValidationErrors(err error, dst *validator.ValidationErrors) bool {
	v, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}
	*dst = v
	return true
}

func messageFor(field, tag string) string {
	switch field {
	case "email":
		if tag == "max" {
			return "email is too long"
		}
		return "valid email is required"
	case "fullName":
		if tag == "name_format" {
			return "full name contains invalid characters"
		}
		return "full name must be 2-100 characters"
	case "phone":
		if tag == "phone_format" {
			return "phone contains invalid characters"
		}
		return "phone must be 7-20 characters"
	}
	return "invalid"
}
