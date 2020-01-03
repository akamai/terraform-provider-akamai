package cps

import (
	"testing"

	//cps "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/stretchr/testify/assert"
)

func TestSHAConversion(t *testing.T) {
	assert.NotEqual(t, "dosty", "everts")
	assert.Equal(t, SHA256, SHA("SHA-256"))
}
