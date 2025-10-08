// Package validators contains custom terraform schema validators
package validators

import "regexp"

// CertificatePEMRegex is a regex to validate PEM encoded certificates.
var CertificatePEMRegex = regexp.MustCompile(
	`^\s*-----BEGIN CERTIFICATE-----\r?\n[A-Za-z0-9+/=\r\n]+-----END CERTIFICATE-----\s*$`)

// ToolchainPEMRegex is a regex to validate a PEM encoding multiple certificates.
var ToolchainPEMRegex = regexp.MustCompile(
	`^\s*(?:-----BEGIN CERTIFICATE-----\r?\n[A-Za-z0-9+/=\r\n]+-----END CERTIFICATE-----\s*)+$`)
