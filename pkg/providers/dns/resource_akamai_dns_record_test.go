package dns

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResDnsRecord(t *testing.T) {
	dnsClient := dns.Client(session.Must(session.New()))

	var rec *dns.RecordBody

	notFound := &dns.Error{
		StatusCode: http.StatusNotFound,
	}

	// This test performs a full life-cycle (CRUD) test
	t.Run("lifecycle test", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, notFound)

		parseCall := client.On("ParseRData",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]string"),
		).Return(nil)

		procCall := client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("string"),
		).Return(nil, nil)

		updateArguments := func(args mock.Arguments) {
			rec = args.Get(1).(*dns.RecordBody)
			getCall.ReturnArguments = mock.Arguments{rec, nil}
			parseCall.ReturnArguments = mock.Arguments{
				dnsClient.ParseRData(context.Background(), rec.RecordType, rec.Target),
			}
			procCall.ReturnArguments = mock.Arguments{rec.Target, nil}
		}

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			updateArguments(args)
		})

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			updateArguments(args)
		})

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]bool"),
		).Return(nil).Run(func(mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{nil, notFound}
		})

		dataSourceName := "akamai_dns_record.a_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("TXT record test", func(t *testing.T) {
		client := &dns.Mock{}

		target1 := "\"Hel\\\\lo\\\"world\""
		target2 := "\"extralongtargetwhichis\" \"intwoseparateparts\""

		client.On("GetRecord",
			mock.Anything,
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			mock.Anything,
			&dns.RecordBody{
				Name:       "exampleterraform.io",
				RecordType: "TXT",
				TTL:        300,
				Active:     false,
				Target:     []string{target1, target2},
			},
			"exampleterraform.io",
			[]bool{false},
		).Return(nil)

		client.On("GetRecord",
			mock.Anything,
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(&dns.RecordBody{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{target1, target2},
		}, nil).Once()

		client.On("ParseRData",
			mock.Anything,
			"TXT",
			[]string{target1, target2},
		).Return(map[string]interface{}{
			"target": []string{target1, target2},
		}).Times(2)

		client.On("ProcessRdata",
			mock.Anything,
			[]string{target1, target2},
			"TXT",
		).Return([]string{target1, target2}).Times(2)

		client.On("GetRecord",
			mock.Anything,
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(&dns.RecordBody{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{target1, target2},
		}, nil).Once()

		client.On("DeleteRecord",
			mock.Anything,
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]bool"),
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/create_basic_txt.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "TXT"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "target.0", "Hel\\lo\"world"),
							resource.TestCheckResourceAttr(resourceName, "target.1", "\"extralongtargetwhichis\" \"intwoseparateparts\""),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestTargetDiffSuppress(t *testing.T) {
	t.Run("target is computed and recordType is AAAA", func(t *testing.T) {
		config := schema.TestResourceDataRaw(t, getResourceDNSRecordSchema(), map[string]interface{}{"recordtype": "AAAA"})
		assert.False(t, dnsRecordTargetSuppress("target.#", "0", "", config))
	})
}

func TestResolveTxtRecordTargets(t *testing.T) {
	denormalized := []string{"onetwo", "\"one\" \"two\""}
	normalized := []string{"\"onetwo\"", "\"one\" \"two\"", "\"one\" \"two\""}
	expected := []string{"onetwo", "\"one\" \"two\"", "\"one\" \"two\""}

	res, err := resolveTxtRecordTargets(denormalized, normalized)
	require.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestResolveTargets(t *testing.T) {
	compare := func(dt, nt string) (bool, error) {
		if dt == "error" {
			return false, fmt.Errorf("oops")
		}
		return nt == strings.ToLower(dt), nil
	}

	tests := map[string]struct {
		denormalized []string
		normalized   []string
		expected     []string
		withError    bool
	}{
		"replaces equal targets": {
			denormalized: []string{"a", "B", "C"},
			normalized:   []string{"a", "b", "c", "d"},
			expected:     []string{"a", "B", "C", "d"},
		},
		"preserves additional normalized targets": {
			denormalized: []string{"a", "b"},
			normalized:   []string{"a", "b", "c", "d"},
			expected:     []string{"a", "b", "c", "d"},
		},
		"does not append additional denormalized targets": {
			denormalized: []string{"a", "b", "C", "D"},
			normalized:   []string{"a", "b"},
			expected:     []string{"a", "b"},
		},
		"preserves normalized targets when elements shift": {
			denormalized: []string{"a", "B", "C"},
			normalized:   []string{"a", "b", "bb", "c"},
			expected:     []string{"a", "B", "bb", "c"},
		},
		"preserves normalized targets when order changes": {
			denormalized: []string{"a", "B", "C", "d"},
			normalized:   []string{"d", "c", "b", "a"},
			expected:     []string{"d", "c", "b", "a"},
		},
		"returns error when normalization failed": {
			denormalized: []string{"error"},
			normalized:   []string{"a"},
			withError:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := resolveTargets(tc.denormalized, tc.normalized, compare)
			if tc.withError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, res)
		})
	}
}
