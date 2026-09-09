package cps

import (
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// testPollChangeStatusInterval is the polling interval for change status checks in tests.
	testPollChangeStatusInterval = 1 * time.Millisecond
	// testPollGetEnrollmentInterval is the polling interval for enrollment checks in tests.
	testPollGetEnrollmentInterval = 1 * time.Millisecond
)

func TestMain(m *testing.M) {
	testutils.TestRunner(m)
}

type (
	// CustomPollingSubprovider is a CPS subprovider with customizable polling intervals for testing.
	CustomPollingSubprovider struct {
		Subprovider
		pollChangeStatusInterval  time.Duration
		pollGetEnrollmentInterval time.Duration
	}
)

// NewCustomPollingSubprovider creates a CPS subprovider with custom polling intervals for testing.
func NewCustomPollingSubprovider(pollChangeStatusInterval, pollGetEnrollmentInterval time.Duration) *CustomPollingSubprovider {
	return &CustomPollingSubprovider{
		pollChangeStatusInterval:  pollChangeStatusInterval,
		pollGetEnrollmentInterval: pollGetEnrollmentInterval,
	}
}

// SDKResources overrides the embedded Subprovider's SDKResources to use test polling intervals.
func (p *CustomPollingSubprovider) SDKResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cps_dv_enrollment":          resourceCPSDVEnrollment(p.pollChangeStatusInterval, p.pollGetEnrollmentInterval),
		"akamai_cps_dv_validation":          resourceCPSDVValidation(p.pollChangeStatusInterval),
		"akamai_cps_third_party_enrollment": resourceCPSThirdPartyEnrollment(p.pollChangeStatusInterval, p.pollGetEnrollmentInterval),
		"akamai_cps_upload_certificate":     resourceCPSUploadCertificate(p.pollChangeStatusInterval),
	}
}

// FrameworkActions overrides the embedded Subprovider's FrameworkActions to use a test polling interval.
func (p *CustomPollingSubprovider) FrameworkActions() []func() action.Action {
	return []func() action.Action{
		NewForceCertificateRenewalAction(p.pollChangeStatusInterval),
	}
}

func TestFrameworkActions(t *testing.T) {
	t.Parallel()

	actions := NewSubprovider().FrameworkActions()
	require.Len(t, actions, 1)

	a, ok := actions[0]().(*ForceCertificateRenewalAction)
	require.True(t, ok)
	assert.Equal(t, defaultPollChangeStatusInterval, a.pollChangeStatusInterval)
}
