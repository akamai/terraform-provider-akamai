//go:build tools

// Package tools is used for managing developer tools for this project
package tools

//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install github.com/terraform-linters/tflint
//go:generate go install golang.org/x/tools/cmd/goimports

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/terraform-linters/tflint"
	_ "golang.org/x/tools/cmd/goimports"
)
