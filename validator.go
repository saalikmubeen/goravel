package goravel

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type ErrorsMap map[string][]string

// Add adds an error message to a given Validation field
func (e ErrorsMap) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message for a given field
func (e ErrorsMap) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

// Validation creates a custom Validation  struct, embeds a url.Values object
type Validation struct {
	// Validation Struct will contain all the key value pairs of the url.Values
	// object(field name and value pairs of html Validation that was submitted by the user)
	Data   url.Values // map[string][]string
	Errors ErrorsMap  // map[string][]string
}

// New initializes a custom Validation struct
func (g *Goravel) Validator(data url.Values) *Validation {
	return &Validation{
		Data:   data,
		Errors: ErrorsMap(map[string][]string{}),
	}
}

// Valid returns true if there are no errors, otherwise false
func (v *Validation) IsValid() bool {
	return len(v.Errors) == 0
}

// Check adds an error message to the validation object if the condition is not met
func (v *Validation) Check(ok bool, field, message string) {
	if !ok {
		v.Errors.Add(field, message)
	}
}

// To check if a Validation field in POST request is empty or not
func (v *Validation) HasValue(field string, value ...string) bool {
	validationValue := ""
	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	if validationValue == "" {
		v.Errors.Add(field, "This field can't be blank")
	}

	return validationValue != ""
}

// To check if any of the Validation field in POST request is empty or not
func (v *Validation) HasRequired(requiredFields ...string) {

	// every required field should be present in the url.Values object
	for _, field := range requiredFields {
		value := v.Data.Get(field)

		if strings.TrimSpace(value) == "" {
			v.Errors.Add(field, "This field can't be blank")
		}
	}
}

func (v *Validation) HasMinLength(field string, length int, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	if len(validationValue) < length {
		v.Errors.Add(field, fmt.Sprintf("This field must be %d charcaters long or more", length))
		return false
	}
	return true
}

func (v *Validation) IsValidEmail(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	if !govalidator.IsEmail(validationValue) {
		v.Errors.Add(field, "Invaid Email")
		return false
	}
	return true

}

func (v *Validation) IsValidUrl(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	if !govalidator.IsURL(validationValue) {
		v.Errors.Add(field, "Invaid URL")
		return false
	}
	return true

}

func (v *Validation) IsValidPassword(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	if !govalidator.IsByteLength(validationValue, 8, 50) {
		v.Errors.Add(field, "Password must be between 8 and 50 characters")
		return false
	}
	return true
}

func (v *Validation) IsValidUsername(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	if !govalidator.IsAlphanumeric(validationValue) {
		v.Errors.Add(field, "Username must be alphanumeric")
		return false
	}
	return true
}

func (v *Validation) IsInt(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	_, err := strconv.Atoi(validationValue)
	if err != nil {
		v.Errors.Add(field, "This field must be an integer")
		return false
	}

	return true
}

func (v *Validation) IsFloat(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	_, err := strconv.ParseFloat(validationValue, 64)
	if err != nil {
		v.Errors.Add(field, "This field must be a float")
		return false
	}

	return true
}

func (v *Validation) IsDateISO(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	_, err := time.Parse("2006-01-02", validationValue)
	if err != nil {
		v.Errors.Add(field, "This field must be a date in the format YYYY-MM-DD")
		return false
	}

	return true
}

func (v *Validation) NoWhitespace(field string, value ...string) bool {
	validationValue := ""

	if len(value) > 0 {
		validationValue = value[0]
	} else {
		validationValue = v.Data.Get(field)
	}

	if govalidator.HasWhitespace(validationValue) {
		v.Errors.Add(field, "This field must not contain any whitespace")
		return false
	}

	return true
}
