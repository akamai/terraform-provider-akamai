package dns

import (
	"context"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"strconv"
	"sync"
)

var (
	zoneRecordsetsWriteLock sync.Mutex
)

// Recordsets contains operations available on a recordsets
// See: https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html
type RecordSets interface {
	// NewRecordSetResponse returns new response object
	NewRecordSetResponse(context.Context, string) *RecordSetResponse
	// GetRecordsets retrieves recordsets with Query Args. No formatting of arg values
	// See: See: https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html#getzonerecordsets
	GetRecordsets(context.Context, string, ...RecordsetQueryArgs) (*RecordSetResponse, error)
	// CreateRecordsets creates multiple recordsets
	// See: https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html#postzonerecordsets
	CreateRecordsets(context.Context, *Recordsets, string, ...bool) error
	// UpdateRecordsets sreplaces list of recordsets
	// See: https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html#putzonerecordsets
	UpdateRecordsets(context.Context, *Recordsets, string, ...bool) error
}

// Recordset Query args struct
type RecordsetQueryArgs struct {
	Page     int
	PageSize int
	Search   string
	ShowAll  bool
	SortBy   string
	Types    string
}

// Recordsets Struct. Used for Create and Update Recordsets
type Recordsets struct {
	Recordsets []Recordset `json:"recordsets"`
}

type Recordset struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	TTL   int      `json:"ttl"`
	Rdata []string `json:"rdata"`
} //`json:"recordsets"`

type MetadataH struct {
	LastPage      int  `json:"lastPage"`
	Page          int  `json:"page"`
	PageSize      int  `json:"pageSize"`
	ShowAll       bool `json:"showAll"`
	TotalElements int  `json:"totalElements"`
} //`json:"metadata"`

type RecordSetResponse struct {
	Metadata   MetadataH   `json:"metadata"`
	Recordsets []Recordset `json:"recordsets"`
}

// Validate validates Recordsets
func (rs *Recordsets) Validate() error {

	if len(rs.Recordsets) < 1 {
		return fmt.Errorf("Request initiated with empty recordsets list")
	}
	for _, rec := range rs.Recordsets {
		err := validation.Errors{
			"Name":  validation.Validate(rec.Name, validation.Required),
			"Type":  validation.Validate(rec.Type, validation.Required),
			"TTL":   validation.Validate(rec.TTL, validation.Required),
			"Rdata": validation.Validate(rec.Rdata, validation.Required),
		}.Filter()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *dns) NewRecordSetResponse(ctx context.Context, name string) *RecordSetResponse {
	recordset := &RecordSetResponse{}
	return recordset
}

// Get RecordSets with Query Args. No formatting of arg values!
func (p *dns) GetRecordsets(ctx context.Context, zone string, queryArgs ...RecordsetQueryArgs) (*RecordSetResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetRecordsets")

	if len(queryArgs) > 1 {
		return nil, fmt.Errorf("invalid arguments GetRecordsets QueryArgs")
	}

	var recordsetResp RecordSetResponse
	getURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", zone)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetRecordsets request: %w", err)
	}

	q := req.URL.Query()
	if len(queryArgs) > 0 {
		if queryArgs[0].Page > 0 {
			q.Add("page", strconv.Itoa(queryArgs[0].Page))
		}
		if queryArgs[0].PageSize > 0 {
			q.Add("pageSize", strconv.Itoa(queryArgs[0].PageSize))
		}
		if queryArgs[0].Search != "" {
			q.Add("search", queryArgs[0].Search)
		}
		q.Add("showAll", strconv.FormatBool(queryArgs[0].ShowAll))
		if queryArgs[0].SortBy != "" {
			q.Add("sortBy", queryArgs[0].SortBy)
		}
		if queryArgs[0].Types != "" {
			q.Add("types", queryArgs[0].Types)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := p.Exec(req, &recordsetResp)
	if err != nil {
		return nil, fmt.Errorf("GetRecordsets request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &recordsetResp, nil
}

// Create Recordstes
func (p *dns) CreateRecordsets(ctx context.Context, recordsets *Recordsets, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordsetsWriteLock.Lock()
		defer zoneRecordsetsWriteLock.Unlock()
	}

	logger := p.Log(ctx)
	logger.Debug("CreateRecordsets")

	if err := recordsets.Validate(); err != nil {
		return err
	}

	reqbody, err := convertStructToReqBody(recordsets)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	postURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create CreateRecordsets request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("CreateRecordsets request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}

func (p *dns) UpdateRecordsets(ctx context.Context, recordsets *Recordsets, zone string, recLock ...bool) error {
	// This lock will restrict the concurrency of API calls
	// to 1 save request at a time. This is needed for the Soa.Serial value which
	// is required to be incremented for every subsequent update to a zone
	// so we have to save just one request at a time to ensure this is always
	// incremented properly

	if localLock(recLock) {
		zoneRecordsetsWriteLock.Lock()
		defer zoneRecordsetsWriteLock.Unlock()
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRecordsets")

	if err := recordsets.Validate(); err != nil {
		return err
	}

	reqbody, err := convertStructToReqBody(recordsets)
	if err != nil {
		return fmt.Errorf("failed to generate request body: %w", err)
	}

	putURL := fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", zone)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, reqbody)
	if err != nil {
		return fmt.Errorf("failed to create UpdateRecordsets request: %w", err)
	}

	resp, err := p.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("UpdateRecordsets request failed: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return p.Error(resp)
	}

	return nil
}
