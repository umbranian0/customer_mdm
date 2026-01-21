package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerConflict = errors.New("customer already exists")
	ErrCustomerBadInput = errors.New("invalid customer input")
)

type ConflictError struct {
	Resource string
	Fields   []string
}

func (e *ConflictError) Error() string {
	if len(e.Fields) == 0 {
		if e.Resource != "" {
			return fmt.Sprintf("%s already exists", e.Resource)
		}
		return ErrCustomerConflict.Error()
	}
	if e.Resource == "" {
		return fmt.Sprintf("conflict on fields: %s", strings.Join(e.Fields, ", "))
	}
	return fmt.Sprintf("%s already exists: %s", e.Resource, strings.Join(e.Fields, ", "))
}

func NewCustomerConflict(fields ...string) error {
	return &ConflictError{Resource: "customer", Fields: fields}
}
