package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testprovider"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var providerBlock = `provider "akamai" {
  edgerc = "../testutils/edgerc"
}
`

func TestStateChecker(t *testing.T) {
	checker := NewStateChecker("akamai_test.sample-resource")

	tests := map[string]struct {
		checks resource.TestCheckFunc
		config string
		error  *regexp.Regexp
	}{
		"happy path - correct attribute value, two checks passing": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckEqual("input", "my output value").
				CheckEqual("output", "my output value").
				Build(),
		},
		"happy path - overwrite attribute check with last CheckEqual": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckEqual("output", "wrong value").
				CheckEqual("output", "my output value").
				Build(),
		},
		"happy path - CheckMissing overwrites CheckEqual for the same attribute": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckEqual("not_existing_attribute", "wrong value").
				CheckMissing("not_existing_attribute").
				Build(),
		},
		"happy path - CheckEqual overwrites CheckMissing for the same attribute": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckMissing("output").
				CheckEqual("output", "my output value").
				Build(),
		},
		"happy path - two attributes are missing": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckMissing("no_such_attribute").
				CheckMissing("no_such_attribute_2").
				Build(),
		},
		"happy path - attribute value is an empty string": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = ""
				}`,
			checks: checker.CheckEqual("output", "").Build(),
		},
		"happy path - attribute value is a number": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "1"
				}`,
			checks: checker.CheckEqual("output", "1").Build(),
		},
		"expect error - check for missing attributes but one is present": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckMissing("no_such_attribute").
				CheckMissing("output").
				Build(),
			error: regexp.MustCompile("Attribute 'output' found when not expected"),
		},
		"expect error - one check fails": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckEqual("input", "my output value").
				CheckEqual("output", "my output value - updated").
				Build(),
			error: regexp.MustCompile("Attribute 'output' expected \"my output value - updated\", got \"my output value\""),
		},
		"expect error - ": {
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckEqual("no_such_attribute", "no_such_value").
				Build(),
			error: regexp.MustCompile("Attribute 'no_such_attribute' not found"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(testprovider.NewMockSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      providerBlock + tc.config,
						Check:       tc.checks,
						ExpectError: tc.error,
					},
				},
			})
		})
	}
}

func TestImportChecker(t *testing.T) {
	checker := NewImportChecker()

	tests := map[string]struct {
		importID string
		config   string
		checks   func(s []*terraform.InstanceState) error
		error    *regexp.Regexp
	}{
		"happy path - correct attribute value": {
			importID: "1,abcdefghijklmnop",
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.CheckEqual("output", "abcdefghijklmnop").Build(),
		},
		"happy path - attribute is missing": {
			importID: "1,abcdefghijklmnop",
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.CheckMissing("no_such_attribute").Build(),
		},
		"happy path - attribute value is an empty string": {
			importID: "1,",
			config: `
				resource "akamai_test" "sample-resource" {
				  input = ""
				}`,
			checks: checker.CheckEqual("output", "").Build(),
		},
		"happy path - attribute value is a number": {
			importID: "1,1",
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "1"
				}`,
			checks: checker.CheckEqual("output", "1").Build(),
		},
		"happy path - overwrite check with correct value": {
			importID: "1,abcdefghijklmnop",
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckEqual("output", "wrong value").
				CheckEqual("output", "abcdefghijklmnop").
				Build(),
		},
		"happy path - overwrite check missing with correct value check": {
			importID: "1,abcdefghijklmnop",
			config: `
				resource "akamai_test" "sample-resource" {
				  input = "my output value"
				}`,
			checks: checker.
				CheckMissing("output").
				CheckEqual("output", "abcdefghijklmnop").
				Build(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(testprovider.NewMockSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:           providerBlock + tc.config,
						ImportStateCheck: tc.checks,
						ImportStateId:    tc.importID,
						ImportState:      true,
						ResourceName:     "akamai_test.sample-resource",
						ExpectError:      tc.error,
					},
				},
			})
		})
	}
}

func TestStateCheckerMethods(t *testing.T) {
	tests := map[string]struct {
		checker            StateChecker
		expectedAttributes map[string]checkData
		error              bool
	}{
		"basic": {
			checker: NewStateChecker("test").
				CheckEqual("attr", "val").
				CheckEqual("attr2", "val2").
				CheckMissing("missing"),
			expectedAttributes: map[string]checkData{
				"attr": {
					isMissing: false,
					value:     "val",
				},
				"attr2": {
					isMissing: false,
					value:     "val2",
				},
				"missing": {
					isMissing: true,
					value:     "",
				},
			},
		},
		"only missing": {
			checker: NewStateChecker("test").
				CheckMissing("missing").
				CheckMissing("missing2").
				CheckMissing("missing3"),
			expectedAttributes: map[string]checkData{
				"missing": {
					isMissing: true,
					value:     "",
				},
				"missing2": {
					isMissing: true,
					value:     "",
				},
				"missing3": {
					isMissing: true,
					value:     "",
				},
			},
		},
		"only existing attributes": {
			checker: NewStateChecker("test").
				CheckEqual("attr", "val").
				CheckEqual("attr2", "val2").
				CheckEqual("attr3", "val3"),
			expectedAttributes: map[string]checkData{
				"attr": {
					isMissing: false,
					value:     "val",
				},
				"attr2": {
					isMissing: false,
					value:     "val2",
				},
				"attr3": {
					isMissing: false,
					value:     "val3",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expectedAttributes, test.checker.attributes)
		})
	}
}

func TestImportCheckerMethods(t *testing.T) {
	importTests := map[string]struct {
		checker            ImportChecker
		expectedAttributes map[string]checkData
	}{
		"basic": {
			checker: NewImportChecker().
				CheckEqual("attr", "val").
				CheckEqual("attr2", "val2").
				CheckMissing("missing"),
			expectedAttributes: map[string]checkData{
				"attr": {
					isMissing: false,
					value:     "val",
				},
				"attr2": {
					isMissing: false,
					value:     "val2",
				},
				"missing": {
					isMissing: true,
					value:     "",
				},
			},
		},
		"only missing": {
			checker: NewImportChecker().
				CheckMissing("missing").
				CheckMissing("missing2").
				CheckMissing("missing3"),
			expectedAttributes: map[string]checkData{
				"missing": {
					isMissing: true,
					value:     "",
				},
				"missing2": {
					isMissing: true,
					value:     "",
				},
				"missing3": {
					isMissing: true,
					value:     "",
				},
			},
		},
		"only existing attributes": {
			checker: NewImportChecker().
				CheckEqual("attr", "val").
				CheckEqual("attr2", "val2").
				CheckEqual("attr3", "val3"),
			expectedAttributes: map[string]checkData{
				"attr": {
					isMissing: false,
					value:     "val",
				},
				"attr2": {
					isMissing: false,
					value:     "val2",
				},
				"attr3": {
					isMissing: false,
					value:     "val3",
				},
			},
		},
	}

	for name, test := range importTests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expectedAttributes, test.checker.attributes)
		})
	}
}

func TestStateCheckersPanic(t *testing.T) {
	stateCheck := func() {
		_ = NewStateChecker("test").Build()
	}
	assert.PanicsWithValue(t, "there must be at least one check in order to build the checker", stateCheck)

	importCheck := func() {
		checkFunc := NewImportChecker().Build()

		givenState := &terraform.InstanceState{
			Attributes: map[string]string{},
		}
		_ = checkFunc([]*terraform.InstanceState{givenState})
	}
	assert.PanicsWithValue(t, "there must be at least one check in order to build the checker", importCheck)
}

func TestAssertAttributeFor(t *testing.T) {
	givenState := &terraform.InstanceState{
		Attributes: map[string]string{
			"attr": "val",
		},
	}

	tests := map[string]struct {
		key         string
		checkData   checkData
		expectError error
	}{
		"'attr' exists and has correct value - no error": {
			key: "attr",
			checkData: checkData{
				isMissing: false,
				value:     "val",
			},
			expectError: nil,
		},
		"'not-existing-attr' does not exists - no error": {
			key: "not-existing-attr",
			checkData: checkData{
				isMissing: true,
			},
			expectError: nil,
		},
		"'not-existing-attr' does not exists, but expecting to exist - expect error": {
			key: "not-existing-attr",
			checkData: checkData{
				isMissing: false,
				value:     "val",
			},
			expectError: fmt.Errorf(`attribute "not-existing-attr" was not present, but should have a value: "val"`),
		},
		"'attr' does exists, but has a different value - expect error": {
			key: "attr",
			checkData: checkData{
				isMissing: false,
				value:     "actual value",
			},
			expectError: fmt.Errorf(`attribute "attr" has incorrect value "val", but should have "actual value"`),
		},
		"'attr' does exists, but should not - expect error": {
			key: "attr",
			checkData: checkData{
				isMissing: true,
				value:     "val",
			},
			expectError: fmt.Errorf(`attribute "attr" was present and has a value: "val", but shouldn't be`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := assertAttributeFor(givenState, test.key, test.checkData)
			if test.expectError != nil {
				assert.EqualError(t, err, test.expectError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
