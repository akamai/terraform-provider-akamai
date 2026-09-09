package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPasswordPolicy(t *testing.T) {
	t.Parallel()
	expectGetPasswordPolicy := func(client *iam.Mock, timesToRun int) {
		passwordPolicyResponse := iam.GetPasswordPolicyResponse{
			CaseDiff:        0,
			MaxRepeating:    2,
			MinDigits:       1,
			MinLength:       8,
			MinLetters:      1,
			MinNonAlpha:     0,
			MinReuse:        4,
			PwClass:         "aka90",
			RotateFrequency: 90,
		}
		client.On("GetPasswordPolicy", testutils.MockContext).Return(&passwordPolicyResponse, nil).Times(timesToRun)
	}

	expectGetPasswordPolicyWithError := func(client *iam.Mock, timesToRun int) {
		client.On("GetPasswordPolicy", testutils.MockContext).Return(nil, fmt.Errorf("get password policy failed")).Times(timesToRun)
	}

	tests := map[string]struct {
		init  func(*iam.Mock)
		error *regexp.Regexp
	}{
		"happy path": {
			init: func(m *iam.Mock) {
				expectGetPasswordPolicy(m, 3)
			},
		},
		"error listing password policy": {
			init: func(m *iam.Mock) {
				expectGetPasswordPolicyWithError(m, 1)
			},
			error: regexp.MustCompile("get password policy failed"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.IAM)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataPasswordPolicy/default.tf"),
						Check:       checkPasswordPolicyAttrs(),
						ExpectError: tc.error,
					},
				},
			})
			client.IAM.AssertExpectations(t)
		})
	}
}

func checkPasswordPolicyAttrs() resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "pw_class", "aka90"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "case_dif", "0"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "max_repeating", "2"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "min_digits", "1"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "min_length", "8"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "min_letters", "1"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "min_non_alpha", "0"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "min_reuse", "4"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_iam_password_policy.test", "rotate_frequency", "90"))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
