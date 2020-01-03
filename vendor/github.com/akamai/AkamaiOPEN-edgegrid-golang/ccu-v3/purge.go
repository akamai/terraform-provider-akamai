package ccu

import (
	"errors"
	fmt "fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type PurgeTypeValue string
type NetworkValue string

var (
	PurgeByUrl      PurgeTypeValue = "url"
	PurgeByCpCode   PurgeTypeValue = "cpcode"
	PurgeByCacheTag PurgeTypeValue = "tag"

	NetworkStaging    NetworkValue = "staging"
	NetworkProduction NetworkValue = "production"
)

type Purge struct {
	Objects []string `json:"objects""`
}

func NewPurge(objects []string) *Purge {
	return &Purge{
		Objects: objects,
	}
}

func (p *Purge) Invalidate(purgeByType PurgeTypeValue, network NetworkValue) (*PurgeResponse, error) {
	return p.purge("invalidate", purgeByType, network)
}

func (p *Purge) Delete(purgeByType PurgeTypeValue, network NetworkValue) (*PurgeResponse, error) {
	return p.purge("delete", purgeByType, network)
}

func (p *Purge) purge(purgeMethod string, purgeByType PurgeTypeValue, network NetworkValue) (*PurgeResponse, error) {
	if len(p.Objects) == 0 {
		return nil, errors.New("one of more purge objects must be defined")
	}

	url := fmt.Sprintf(
		"/ccu/v3/%s/%s/%s",
		purgeMethod,
		purgeByType,
		network,
	)

	req, err := client.NewJSONRequest(Config, "POST", url, p)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	purge := &PurgeResponse{}
	if err = client.BodyJSON(res, purge); err != nil {
		return nil, err
	}

	return purge, nil
}

type PurgeResponse struct {
	PurgeID          string `json:"purgeId"`
	EstimatedSeconds int    `json:"estimatedSeconds"`
	HTTPStatus       int    `json:"httpStatus"`
	Detail           string `json:"detail"`
	SupportID        string `json:"supportId"`
}
