package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Get returns the global validator instance
func Get() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

// Struct validates a struct and returns the error
func Struct(s interface{}) error {
	return Get().Struct(s)
}
