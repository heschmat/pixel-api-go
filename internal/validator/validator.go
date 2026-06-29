package validator

import "slices"

type Validator struct {
	Errors map[string]string
}

// instantiate a new Validator type with a map of validation errors
// the function returns a pointer to a Validator
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) addError(key, msg string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = msg
	}
}

func (v *Validator) Check(ok bool, key, msg string) {
	if !ok {
		v.addError(key, msg)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// generic functions -----
// true: if a specific value is in a list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]string)

	for _, value := range values {
		uniqueValues[value] = ""
	}

	return len(values) == len(uniqueValues)
}
