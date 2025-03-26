package validators

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestFileReadableValidator_FileNotExists(t *testing.T) {
	ctx := context.Background()

	pathValidator := FileReadable()

	req := validator.StringRequest{
		ConfigValue: types.StringValue("dont-exist.txt"),
	}

	resp := validator.StringResponse{
		Diagnostics: diag.Diagnostics{},
	}

	pathValidator.ValidateString(ctx, req, &resp)
	assert.True(t, resp.Diagnostics.HasError(), "Expected an error but got none")
}

func TestFileReadableValidator_FileExists(t *testing.T) {
	ctx := context.Background()

	file := createTempFile(0644)

	pathValidator := FileReadable()

	req := validator.StringRequest{
		ConfigValue: types.StringValue(file.Name()),
	}

	resp := validator.StringResponse{
		Diagnostics: diag.Diagnostics{},
	}

	pathValidator.ValidateString(ctx, req, &resp)

	assert.False(t, resp.Diagnostics.HasError(), "Expected no error but got one")
}

func createTempFile(mode os.FileMode) *os.File {
	file, err := os.CreateTemp("", "exists-*.txt")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if err := os.Chmod(file.Name(), mode); err != nil {
		panic(err)
	}

	return file
}
