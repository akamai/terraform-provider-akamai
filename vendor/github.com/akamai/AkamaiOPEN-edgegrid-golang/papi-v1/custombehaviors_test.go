package papi

import (
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestCustomBehaviors_GetCustomBehaviors(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/papi/v1/custom-behaviors")
	mock.
		Get("/papi/v1/custom-behaviors").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"accountId": "act_1-1TJZFB",
				"customBehaviors": {
					"items": [
						{
							"behaviorId": "cbe_12345",
							"name": "DLR",
							"status": "ACTIVE",
							"displayName": "Custom Download Receipt",
							"description": "Setting custom download receipt. Uses PMUSER_LOG variable.",
							"updatedDate": "2017-04-24T12:34:56Z",
							"updatedByUser": "jsikkela"
						}
					]
				}
			}`)

	Init(config)

	behaviors := NewCustomBehaviors()
	err := behaviors.GetCustomBehaviors()

	assert.NoError(t, err)
	assert.Len(t, behaviors.CustomBehaviors.Items, 1)
	assert.Equal(t, "cbe_12345", behaviors.CustomBehaviors.Items[0].BehaviorID)
	assert.Equal(t, "DLR", behaviors.CustomBehaviors.Items[0].Name)
	assert.Equal(t, "ACTIVE", behaviors.CustomBehaviors.Items[0].Status)
	assert.Equal(t, "Custom Download Receipt", behaviors.CustomBehaviors.Items[0].DisplayName)
	assert.Equal(t, "Setting custom download receipt. Uses PMUSER_LOG variable.", behaviors.CustomBehaviors.Items[0].Description)
	time, _ := time.Parse("2006-01-02T15:04:05Z", "2017-04-24T12:34:56Z")
	assert.Equal(t, time, behaviors.CustomBehaviors.Items[0].UpdatedDate)
	assert.Equal(t, "jsikkela", behaviors.CustomBehaviors.Items[0].UpdatedByUser)
}

func TestCustomBehavior_GetCustomBehavior(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/papi/v1/custom-behaviors/cbe_12345")
	mock.
		Get("/papi/v1/custom-behaviors/cbe_12345").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"accountId": "act_1-1TJZFB",
				"customBehaviors": {
					"items": [
						{
							"behaviorId": "cbe_12345",
							"name": "DLR",
							"status": "ACTIVE",
							"displayName": "Custom Download Receipt",
							"description": "Setting custom download receipt. Uses PMUSER_LOG variable.",
							"updatedDate": "2017-04-24T12:34:56Z",
							"updatedByUser": "jsikkela"
						}
					]
				}
			}`)

	Init(config)

	behavior := NewCustomBehavior(NewCustomBehaviors())
	behavior.BehaviorID = "cbe_12345"
	err := behavior.GetCustomBehavior()

	assert.NoError(t, err)
	assert.Equal(t, "cbe_12345", behavior.BehaviorID)
	assert.Equal(t, "DLR", behavior.Name)
	assert.Equal(t, "ACTIVE", behavior.Status)
	assert.Equal(t, "Custom Download Receipt", behavior.DisplayName)
	assert.Equal(t, "Setting custom download receipt. Uses PMUSER_LOG variable.", behavior.Description)
	time, _ := time.Parse("2006-01-02T15:04:05Z", "2017-04-24T12:34:56Z")
	assert.Equal(t, time, behavior.UpdatedDate)
	assert.Equal(t, "jsikkela", behavior.UpdatedByUser)
}
