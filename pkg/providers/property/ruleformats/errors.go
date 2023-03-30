package ruleformats

import (
	"errors"
	"fmt"
	"strings"
)

// NotFoundError defines an error that is used when value is not found in the schema.ResourceData.
type NotFoundError struct {
	key string
}

// TypeAssertionError defines an error that is used when type assertion failed.
type TypeAssertionError struct {
	want string
	got  string
	key  string
}

// TooManyElementsError defines an error that is used when too many elements were found in the list.
type TooManyElementsError struct {
	names    []string
	expected int
}

var (
	// ErrNotFound is used when value not found in schema.ResourceData.
	ErrNotFound = errors.New("not found")
	// ErrTypeAssertion is used when type assertion failed.
	ErrTypeAssertion = errors.New("type assertion")
	// ErrTooManyElements is used when too many elements were found in a list.
	ErrTooManyElements = errors.New("too many elements")
)

// Error returns NotFoundError as a string.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.key)
}

// Unwrap unwraps NotFoundError into ErrNotFound.
func (e *NotFoundError) Unwrap() error {
	return ErrNotFound
}

// Error returns TypeAssertionError as a string.
func (e *TypeAssertionError) Error() string {
	if len(e.key) > 0 {
		return fmt.Sprintf("%s: type assertion failed, want '%s', got '%s'", e.key, e.want, e.got)
	}
	return fmt.Sprintf("type assertion failed, want '%s', got '%s'", e.want, e.got)
}

// Unwrap unwraps TypeAssertionError into ErrTypeAssertion.
func (e *TypeAssertionError) Unwrap() error {
	return ErrTypeAssertion
}

// Error returns TooManyElementsError as a string.
func (e *TooManyElementsError) Error() string {
	return fmt.Sprintf("expected %d element(s), got %d - [%v]", e.expected, len(e.names), strings.Join(e.names, ", "))
}

// Unwrap unwraps TooManyElementsError into ErrTooManyElements.
func (e *TooManyElementsError) Unwrap() error {
	return ErrTooManyElements
}

func typeof(v any) string {
	return fmt.Sprintf("%T", v)
}
