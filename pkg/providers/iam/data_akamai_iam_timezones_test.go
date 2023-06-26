package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataTimezones(t *testing.T) {
	workDir := t.Name()
	t.Run("happy path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(test.TattleT{T: t})

		client.On("SupportedTimezones", mock.Anything).Return([]iam.Timezone{
			{
				Timezone:    "Asia/Kolkata",
				Description: "Asia/Kolkata",
				Offset:      "+5:30",
				Posix:       "Asia/Kolkata",
			},
			{
				Timezone:    "America/Caracas",
				Description: "America/Caracas",
				Offset:      "-4",
				Posix:       "America/Caracas",
			},
			{
				Timezone:    "Europe/Budapest",
				Description: "Europe/Budapest",
				Offset:      "+1",
				Posix:       "Europe/Budapest",
			},
		}, nil)

		useClient(client, func() {
			resourceName := "data.akamai_iam_timezones.test"
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				IsUnitTest:        true,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/%s/timezones.tf", workDir),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckResourceAttr(resourceName, "timezones.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "timezones.0.timezone", "America/Caracas"),
							resource.TestCheckResourceAttr(resourceName, "timezones.0.description", "America/Caracas"),
							resource.TestCheckResourceAttr(resourceName, "timezones.0.posix", "America/Caracas"),
							resource.TestCheckResourceAttr(resourceName, "timezones.0.offset", "-4"),
							resource.TestCheckResourceAttr(resourceName, "timezones.1.timezone", "Asia/Kolkata"),
							resource.TestCheckResourceAttr(resourceName, "timezones.1.description", "Asia/Kolkata"),
							resource.TestCheckResourceAttr(resourceName, "timezones.1.posix", "Asia/Kolkata"),
							resource.TestCheckResourceAttr(resourceName, "timezones.1.offset", "+5:30"),
							resource.TestCheckResourceAttr(resourceName, "timezones.2.timezone", "Europe/Budapest"),
							resource.TestCheckResourceAttr(resourceName, "timezones.2.description", "Europe/Budapest"),
							resource.TestCheckResourceAttr(resourceName, "timezones.2.posix", "Europe/Budapest"),
							resource.TestCheckResourceAttr(resourceName, "timezones.2.offset", "+1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("fail path", func(t *testing.T) {
		client := &iam.Mock{}
		client.Test(test.TattleT{T: t})
		client.On("SupportedTimezones", mock.Anything).Return([]iam.Timezone{}, fmt.Errorf("supported timezones: timezones could not be fetched"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				IsUnitTest:        true,
				Steps: []resource.TestStep{
					{
						Config:      test.Fixture("testdata/%s/timezones.tf", workDir),
						ExpectError: regexp.MustCompile("supported timezones: timezones could not be fetched"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
