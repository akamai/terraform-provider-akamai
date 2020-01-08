package ccu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestPurge_Invalidate(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/ccu/v3/invalidate/url/production")
	mock.
		Post("/ccu/v3/invalidate/url/production").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"detail": "Request accepted",
				"estimatedSeconds": 5,
				"httpStatus": 201,
				"purgeId": "674e54ae-3131-11e8-ba75-615d2757a3f3",
				"supportId": "17PY1522094889114372-178558144"
			}`)

	Init(config)
	purge := NewPurge([]string{"https://www.daveyshafik.com"})
	res, err := purge.Invalidate(PurgeByUrl, NetworkProduction)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.Detail, "Request accepted")
	assert.Equal(t, res.EstimatedSeconds, 5)
	assert.Equal(t, res.HTTPStatus, 201)
	assert.Equal(t, res.PurgeID, "674e54ae-3131-11e8-ba75-615d2757a3f3")
	assert.Equal(t, res.SupportID, "17PY1522094889114372-178558144")
}

func TestPurge_Delete(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/ccu/v3/delete/url/production")
	mock.
		Post("/ccu/v3/delete/url/production").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"detail": "Request accepted",
				"estimatedSeconds": 5,
				"httpStatus": 201,
				"purgeId": "674e54ae-3131-11e8-ba75-615d2757a3f3",
				"supportId": "17PY1522094889114372-178558144"
			}`)

	Init(config)
	purge := NewPurge([]string{"https://www.daveyshafik.com"})
	res, err := purge.Delete(PurgeByUrl, NetworkProduction)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.Detail, "Request accepted")
	assert.Equal(t, res.EstimatedSeconds, 5)
	assert.Equal(t, res.HTTPStatus, 201)
	assert.Equal(t, res.PurgeID, "674e54ae-3131-11e8-ba75-615d2757a3f3")
	assert.Equal(t, res.SupportID, "17PY1522094889114372-178558144")
}

func TestPurge_Invalidate_CpCode(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/ccu/v3/invalidate/cpcode/production")
	mock.
		Post("/ccu/v3/invalidate/cpcode/production").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"detail": "Request accepted",
				"estimatedSeconds": 5,
				"httpStatus": 201,
				"purgeId": "674e54ae-3131-11e8-ba75-615d2757a3f3",
				"supportId": "17PY1522094889114372-178558144"
			}`)

	Init(config)
	purge := NewPurge([]string{"cpc_12345"})
	res, err := purge.Invalidate(PurgeByCpCode, NetworkProduction)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.Detail, "Request accepted")
	assert.Equal(t, res.EstimatedSeconds, 5)
	assert.Equal(t, res.HTTPStatus, 201)
	assert.Equal(t, res.PurgeID, "674e54ae-3131-11e8-ba75-615d2757a3f3")
	assert.Equal(t, res.SupportID, "17PY1522094889114372-178558144")
}

func TestPurge_Invalidate_CacheTag(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/ccu/v3/invalidate/tag/production")
	mock.
		Post("/ccu/v3/invalidate/tag/production").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"detail": "Request accepted",
				"estimatedSeconds": 5,
				"httpStatus": 201,
				"purgeId": "674e54ae-3131-11e8-ba75-615d2757a3f3",
				"supportId": "17PY1522094889114372-178558144"
			}`)

	Init(config)
	purge := NewPurge([]string{"https://www.daveyshafik.com"})
	res, err := purge.Invalidate(PurgeByCacheTag, NetworkProduction)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.Detail, "Request accepted")
	assert.Equal(t, res.EstimatedSeconds, 5)
	assert.Equal(t, res.HTTPStatus, 201)
	assert.Equal(t, res.PurgeID, "674e54ae-3131-11e8-ba75-615d2757a3f3")
	assert.Equal(t, res.SupportID, "17PY1522094889114372-178558144")
}

func TestPurge_Delete_CpCode(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/ccu/v3/delete/cpcode/production")
	mock.
		Post("/ccu/v3/delete/cpcode/production").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"detail": "Request accepted",
				"estimatedSeconds": 5,
				"httpStatus": 201,
				"purgeId": "674e54ae-3131-11e8-ba75-615d2757a3f3",
				"supportId": "17PY1522094889114372-178558144"
			}`)

	Init(config)
	purge := NewPurge([]string{"cpc_12345"})
	res, err := purge.Delete(PurgeByCpCode, NetworkProduction)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.Detail, "Request accepted")
	assert.Equal(t, res.EstimatedSeconds, 5)
	assert.Equal(t, res.HTTPStatus, 201)
	assert.Equal(t, res.PurgeID, "674e54ae-3131-11e8-ba75-615d2757a3f3")
	assert.Equal(t, res.SupportID, "17PY1522094889114372-178558144")
}

func TestPurge_Delete_CacheTag(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/ccu/v3/delete/tag/production")
	mock.
		Post("/ccu/v3/delete/tag/production").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
				"detail": "Request accepted",
				"estimatedSeconds": 5,
				"httpStatus": 201,
				"purgeId": "674e54ae-3131-11e8-ba75-615d2757a3f3",
				"supportId": "17PY1522094889114372-178558144"
			}`)

	Init(config)
	purge := NewPurge([]string{"https://www.daveyshafik.com"})
	res, err := purge.Delete(PurgeByCacheTag, NetworkProduction)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.Detail, "Request accepted")
	assert.Equal(t, res.EstimatedSeconds, 5)
	assert.Equal(t, res.HTTPStatus, 201)
	assert.Equal(t, res.PurgeID, "674e54ae-3131-11e8-ba75-615d2757a3f3")
	assert.Equal(t, res.SupportID, "17PY1522094889114372-178558144")
}
