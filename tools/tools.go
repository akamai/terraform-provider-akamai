//go:build tools

// Package tools is used for managing developer tools for this project
package tools

//go:generate go install golang.org/x/tools/cmd/goimports

import (
	_ "golang.org/x/tools/cmd/goimports"
)
