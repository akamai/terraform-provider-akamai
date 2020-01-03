package papi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestCustomOverrides_GetCustomOverrides(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/papi/v1/custom-overrides")
	mock.
		Get("/papi/v1/custom-overrides").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"accountId": "act_1-1TJZFB",
				"customOverrides": {
					"items": [
						{
							"overrideId": "cbo_12345",
							"displayName": "MDC Behavior",
							"description": "Multiple Domain Configuration can be used to ...",
							"name": "mdc",
							"status": "ACTIVE",
							"updatedByUser": "jsikkela",
							"updatedDate": "2017-04-24T12:34:56Z"
						}
					]
				}
			}`)

	Init(config)

	overrides := NewCustomOverrides()
	err := overrides.GetCustomOverrides()

	assert.NoError(t, err)
	assert.Len(t, overrides.CustomOverrides.Items, 1)
	assert.Equal(t, "cbo_12345", overrides.CustomOverrides.Items[0].OverrideID)
	assert.Equal(t, "mdc", overrides.CustomOverrides.Items[0].Name)
	assert.Equal(t, "ACTIVE", overrides.CustomOverrides.Items[0].Status)
	assert.Equal(t, "MDC Behavior", overrides.CustomOverrides.Items[0].DisplayName)
	assert.Equal(t, "Multiple Domain Configuration can be used to ...", overrides.CustomOverrides.Items[0].Description)
	time, _ := time.Parse("2006-01-02T15:04:05Z", "2017-04-24T12:34:56Z")
	assert.Equal(t, time, overrides.CustomOverrides.Items[0].UpdatedDate)
	assert.Equal(t, "jsikkela", overrides.CustomOverrides.Items[0].UpdatedByUser)
}

func TestCustomOverride_GetCustomOverride(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/papi/v1/custom-overrides/cbo_12345")
	mock.
		Get("/papi/v1/custom-overrides/cbo_12345").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"accountId": "act_1-1TJZFB",
				"customOverrides": {
					"items": [
						{
							"overrideId": "cbo_12345",
							"displayName": "MDC Behavior",
							"description": "Multiple Domain Configuration can be used to ...",
							"name": "mdc",
							"status": "ACTIVE",
							"updatedByUser": "jsikkela",
							"updatedDate": "2017-04-24T12:34:56Z"
						}
					]
				}
			}`)

	Init(config)

	override := NewCustomOverride(NewCustomOverrides())
	override.OverrideID = "cbo_12345"
	err := override.GetCustomOverride()

	assert.NoError(t, err)
	assert.Equal(t, "cbo_12345", override.OverrideID)
	assert.Equal(t, "mdc", override.Name)
	assert.Equal(t, "ACTIVE", override.Status)
	assert.Equal(t, "MDC Behavior", override.DisplayName)
	assert.Equal(t, "Multiple Domain Configuration can be used to ...", override.Description)
	time, _ := time.Parse("2006-01-02T15:04:05Z", "2017-04-24T12:34:56Z")
	assert.Equal(t, time, override.UpdatedDate)
	assert.Equal(t, "jsikkela", override.UpdatedByUser)
}
