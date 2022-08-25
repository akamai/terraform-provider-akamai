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
		"no sans provided": {
			givenCSR: schema.NewSet(schema.HashResource(resource), []interface{}{map[string]interface{}{
				"city":                "Cambridge",
				"state":               "MA",
				"country_code":        "US",
				"organization":        "Akamai",
				"organizational_unit": "test_ou",
			}}),
			givenCN:   "test.com",
			givenSANS: nil,
			expected: &cps.CSR{
				C:    "US",
				CN:   "test.com",
				L:    "Cambridge",
				O:    "Akamai",
				OU:   "test_ou",
				SANS: nil,
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

func TestGetNetworkConfig(t *testing.T) {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"network_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_mutual_authentication": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_ca_list_to_client": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ocsp_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"set_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"disallowed_tls_versions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"clone_dns_names": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"geography": {
							Type:     schema.TypeString,
							Required: true,
						},
						"must_have_ciphers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ocsp_stapling": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"preferred_ciphers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"quic_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"sans": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"secure_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sni_only": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
		},
	}
	truePtr := true
	tests := map[string]struct {
		givenNetworkConfig *schema.Set
		givenSANS          *schema.Set
		givenSniOnly       bool
		expected           *cps.NetworkConfiguration
		withError          bool
	}{
		"basic test": {
			givenNetworkConfig: schema.NewSet(schema.HashResource(resource), []interface{}{map[string]interface{}{
				"client_mutual_authentication": schema.NewSet(schema.HashResource(resource), []interface{}{map[string]interface{}{
					"send_ca_list_to_client": true,
					"ocsp_enabled":           true,
					"set_id":                 "123",
				}}),
				"disallowed_tls_versions": []string{"TLSv1"},
				"clone_dns_names":         true,
				"geography":               "core",
				"must_have_ciphers":       "ak-akamai-default",
				"ocsp_stapling":           "on",
				"preferred_ciphers":       "ak-akamai-default",
				"quic_enabled":            true,
			}}),
			givenSANS:    schema.NewSet(schema.HashString, []interface{}{"a.com", "b.com"}),
			givenSniOnly: true,
			expected: &cps.NetworkConfiguration{
				ClientMutualAuthentication: &cps.ClientMutualAuthentication{
					AuthenticationOptions: &cps.AuthenticationOptions{
						OCSP:               &cps.OCSP{Enabled: &truePtr},
						SendCAListToClient: &truePtr,
					},
					SetID: "123",
				},
				DisallowedTLSVersions: []string{"TLSv1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: true,
					DNSNames:      []string{"a.com", "b.com"},
				},
				Geography:        "core",
				MustHaveCiphers:  "ak-akamai-default",
				OCSPStapling:     "on",
				PreferredCiphers: "ak-akamai-default",
				QuicEnabled:      true,
				SecureNetwork:    "enhanced_tls",
				SNIOnly:          true,
			},
		},
		"only required values with sni_only=true": {
			givenNetworkConfig: schema.NewSet(schema.HashResource(resource), []interface{}{map[string]interface{}{
				"geography": "core",
			}}),
			givenSANS:    nil,
			givenSniOnly: true,
			expected: &cps.NetworkConfiguration{
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      nil,
				},
				Geography:     "core",
				SecureNetwork: "enhanced_tls",
				SNIOnly:       true,
			},
		},
		"only required values with sni_only=false": {
			givenNetworkConfig: schema.NewSet(schema.HashResource(resource), []interface{}{map[string]interface{}{
				"geography": "core",
			}}),
			givenSANS: nil,
			expected: &cps.NetworkConfiguration{
				DNSNameSettings: nil,
				Geography:       "core",
				SecureNetwork:   "enhanced_tls",
				SNIOnly:         false,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rd := resource.TestResourceData()
			err := rd.Set("network_configuration", test.givenNetworkConfig)
			require.NoError(t, err)
			err = rd.Set("secure_network", "enhanced_tls")
			require.NoError(t, err)
			err = rd.Set("sni_only", test.givenSniOnly)
			require.NoError(t, err)
			err = rd.Set("sans", test.givenSANS)
			require.NoError(t, err)

			res, err := GetNetworkConfig(rd)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetOrg(t *testing.T) {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address_line_one": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address_line_two": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"city": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"postal_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
	tests := map[string]struct {
		givenOrg *schema.Set
		expected *cps.Org
	}{
		"basic test": {
			givenOrg: schema.NewSet(schema.HashResource(resource), []interface{}{map[string]interface{}{
				"name":             "Akamai",
				"phone":            "123123123",
				"address_line_one": "test line 1",
				"address_line_two": "test line 2",
				"city":             "Cambridge",
				"region":           "MA",
				"postal_code":      "12345",
				"country_code":     "US",
			}}),
			expected: &cps.Org{
				AddressLineOne: "test line 1",
				AddressLineTwo: "test line 2",
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "123123123",
				PostalCode:     "12345",
				Region:         "MA",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rd := resource.TestResourceData()
			err := rd.Set("organization", test.givenOrg)

			res, err := GetOrg(rd)
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestContactInfoToMap(t *testing.T) {
	tests := map[string]struct {
		given    cps.Contact
		expected map[string]interface{}
	}{
		"basic test": {
			given: cps.Contact{
				AddressLineOne:   "test line 1",
				AddressLineTwo:   "test line 2",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
				Title:            "Mr",
			},
			expected: map[string]interface{}{
				"first_name":       "R1",
				"last_name":        "D1",
				"title":            "Mr",
				"organization":     "Akamai",
				"email":            "r1d1@akamai.com",
				"phone":            "123123123",
				"address_line_one": "test line 1",
				"address_line_two": "test line 2",
				"city":             "Cambridge",
				"region":           "MA",
				"postal_code":      "12345",
				"country_code":     "US",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := ContactInfoToMap(test.given)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestOrgToMap(t *testing.T) {
	tests := map[string]struct {
		given    cps.Org
		expected map[string]interface{}
	}{
		"basic test": {
			given: cps.Org{
				AddressLineOne: "test line 1",
				AddressLineTwo: "test line 2",
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "123123123",
				PostalCode:     "12345",
				Region:         "MA",
			},
			expected: map[string]interface{}{
				"name":             "Akamai",
				"phone":            "123123123",
				"address_line_one": "test line 1",
				"address_line_two": "test line 2",
				"city":             "Cambridge",
				"region":           "MA",
				"postal_code":      "12345",
				"country_code":     "US",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := OrgToMap(test.given)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestCSRToMap(t *testing.T) {
	tests := map[string]struct {
		given    cps.CSR
		expected map[string]interface{}
	}{
		"basic test": {
			given: cps.CSR{
				C:  "US",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			expected: map[string]interface{}{
				"country_code":        "US",
				"city":                "Cambridge",
				"state":               "MA",
				"organization":        "Akamai",
				"organizational_unit": "WebEx",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := CSRToMap(test.given)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestNetworkConfigToMap(t *testing.T) {
	truePtr := true
	tests := map[string]struct {
		given    cps.NetworkConfiguration
		expected map[string]interface{}
	}{
		"basic test": {
			given: cps.NetworkConfiguration{
				ClientMutualAuthentication: &cps.ClientMutualAuthentication{
					AuthenticationOptions: &cps.AuthenticationOptions{
						OCSP:               &cps.OCSP{Enabled: &truePtr},
						SendCAListToClient: &truePtr,
					},
					SetID: "123",
				},
				DisallowedTLSVersions: []string{"TLSv1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: true,
					DNSNames:      []string{"a.com", "b.com"},
				},
				Geography:        "core",
				MustHaveCiphers:  "ak-akamai-default",
				OCSPStapling:     "on",
				PreferredCiphers: "ak-akamai-default",
				QuicEnabled:      true,
				SecureNetwork:    "enhanced_tls",
				SNIOnly:          true,
			},
			expected: map[string]interface{}{
				"client_mutual_authentication": []interface{}{map[string]interface{}{
					"send_ca_list_to_client": true,
					"ocsp_enabled":           true,
					"set_id":                 "123",
				}},
				"disallowed_tls_versions": []string{"TLSv1"},
				"clone_dns_names":         true,
				"geography":               "core",
				"must_have_ciphers":       "ak-akamai-default",
				"ocsp_stapling":           "on",
				"preferred_ciphers":       "ak-akamai-default",
				"quic_enabled":            true,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := NetworkConfigToMap(test.given)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetChangeIDFromPendingChanges(t *testing.T) {
	tests := map[string]struct {
		givenChanges []string
		expected     int
		withError    bool
	}{
		"basic test": {
			givenChanges: []string{"/cps/enrollments/1/changes/2"},
			expected:     2,
		},
		"no pending changes provided": {
			givenChanges: nil,
			withError:    true,
		},
		"invalid change ID": {
			givenChanges: []string{"/cps/enrollments/1/changes/abc"},
			withError:    true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := GetChangeIDFromPendingChanges(test.givenChanges)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetEnrollmentID(t *testing.T) {
	tests := map[string]struct {
		givenLocation string
		expected      int
		withError     bool
	}{
		"basic test": {
			givenLocation: "/cps/v2/enrollments/5555",
			expected:      5555,
		},
		"invalid enrollment ID": {
			givenLocation: "/cps/v2/enrollments/abc",
			withError:     true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := GetEnrollmentID(test.givenLocation)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}
