package akamai

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type (
	// Error is the akamai error interface
	Error struct {
		summary  string
		notFound bool
	}
)

var (
	// ErrDuplicateSchemaKey is returned when a duplicate schema key is detected during merge
	ErrDuplicateSchemaKey = &Error{"duplicate schema key", false}

	// ErrNoConfiguredProviders is returned when no providers are configured
	ErrNoConfiguredProviders = &Error{"One or more Akamai Edgegrid provider configurations must be defined", false}

	// ErrCacheEntryNotFound returns a cache entry error
	ErrCacheEntryNotFound = func(key string) error { return &Error{fmt.Sprintf("cache entry %q not found", key), true} }

	// ErrProviderNotLoaded returned and panic'd when a requested provider is not loaded
	// Users should never see this, unit tests and sanity checks should pick this up
	ErrProviderNotLoaded = func(name string) error { return &Error{fmt.Sprintf("Provider %q not loaded", name), false} }

	// NoticeDeprecatedUseAlias is returned for schema configurations that are deprecated
	// Terraform now supports section aliases
	// TODO: Add alias example to the examples directory
	NoticeDeprecatedUseAlias = func(n string) string {
		return fmt.Sprintf(`The setting %q has been deprecated. See:
	https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-configurations`, n)
	}
)

// Diagnostic converts the error a diagnostic
func (e Error) Diagnostic(detail ...string) diag.Diagnostic {
	d := diag.Diagnostic{
		Severity: diag.Error,
		Summary:  e.Error(),
	}

	if len(detail) > 0 {
		d.Detail = detail[0]
	}

	return d
}

// Diagnostics converts the error to a diag.DiaDiagnostics
func (e Error) Diagnostics(detail ...string) diag.Diagnostics {
	return []diag.Diagnostic{e.Diagnostic(detail...)}
}

// Error implements the error interface
func (e Error) Error() string {
	return e.summary
}

// IsNotFoundError is returned if the error has the notFound flag set
func IsNotFoundError(e error) bool {
	if e, ok := e.(*Error); ok {
		return e.notFound
	}

	return false
}
