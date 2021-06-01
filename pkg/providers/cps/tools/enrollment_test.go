package tools

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContactInfo(t *testing.T) {
	hashFunc := func(i interface{}) int { return 0 }
	tests := map[string]struct {
		given     *schema.Set
		expected  *cps.Contact
		withError bool
	}{
		"basic test": {
			given: schema.NewSet(hashFunc, []interface{}{map[string]interface{}{
				"first_name":       "R1",
				"last_name":        "D1",
				"title":            "mr",
				"organization":     "Akamai",
				"email":            "r1d1@akamai.com",
				"phone":            "123123123",
				"address_line_one": "abc",
				"address_line_two": "def",
				"city":             "Cambridge",
				"region":           "MA",
				"postal_code":      "12345",
				"country_code":     "US",
			}}),
			expected: &cps.Contact{
				AddressLineOne:   "abc",
				AddressLineTwo:   "def",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
				Title:            "mr",
			},
		},
		"set value is of invalid type": {
			given:     schema.NewSet(schema.HashString, []interface{}{"abc"}),
			withError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := GetContactInfo(test.given)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetCSR(t *testing.T) {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"common_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sans": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"csr": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"city": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organizational_unit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"state": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
	tests := map[string]struct {
		givenCSR  *schema.Set
		givenCN   string
		givenSANS *schema.Set
		expected  *cps.CSR
		withError bool
	}{
		"basic test": {
			givenCSR: schema.NewSet(schema.HashResource(resource), []interface{}{map[string]interface{}{
				"city":                "Cambridge",
				"state":               "MA",
				"country_code":        "US",
				"organization":        "Akamai",
				"organizational_unit": "test_ou",
			}}),
			givenCN:   "test.com",
			givenSANS: schema.NewSet(schema.HashString, []interface{}{"a.com", "b.com"}),
			expected: &cps.CSR{
				C:    "US",
				CN:   "test.com",
				L:    "Cambridge",
				O:    "Akamai",
				OU:   "test_ou",
				SANS: []string{"a.com", "b.com"},
				ST:   "MA",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rd := resource.TestResourceData()
			err := rd.Set("csr", test.givenCSR)
			require.NoError(t, err)
			err = rd.Set("common_name", test.givenCN)
			require.NoError(t, err)
			err = rd.Set("sans", test.givenSANS)
			require.NoError(t, err)

			res, err := GetCSR(rd)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}
