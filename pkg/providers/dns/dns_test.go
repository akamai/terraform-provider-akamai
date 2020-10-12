package dns

import (
	"context"
	"net"

	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configdns"
	"github.com/stretchr/testify/mock"
)

type mockdns struct {
	mock.Mock
}

func (d *mockdns) ListZones(ctx context.Context, query ...dns.ZoneListQueryArgs) (*dns.ZoneListResponse, error) {
	var args mock.Arguments

	if len(query) > 0 {
		args = d.Called(ctx, query[0])
	} else {
		args = d.Called(ctx)
	}

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.ZoneListResponse), args.Error(1)
}

func (d *mockdns) NewZone(ctx context.Context, params dns.ZoneCreate) *dns.ZoneCreate {
	args := d.Called(ctx, params)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*dns.ZoneCreate)
}

func (d *mockdns) NewZoneResponse(ctx context.Context, param string) *dns.ZoneResponse {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*dns.ZoneResponse)
}

func (d *mockdns) NewChangeListResponse(ctx context.Context, param string) *dns.ChangeListResponse {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*dns.ChangeListResponse)
}

func (d *mockdns) NewZoneQueryString(ctx context.Context, param1 string, param2 string) *dns.ZoneQueryString {
	args := d.Called(ctx, param1, param1)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*dns.ZoneQueryString)
}

func (d *mockdns) GetZone(ctx context.Context, name string) (*dns.ZoneResponse, error) {
	args := d.Called(ctx, name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.ZoneResponse), args.Error(1)
}

func (d *mockdns) GetChangeList(ctx context.Context, param string) (*dns.ChangeListResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.ChangeListResponse), args.Error(1)
}

func (d *mockdns) GetMasterZoneFile(ctx context.Context, param string) (string, error) {
	args := d.Called(ctx, param)

	return args.String(0), args.Error(1)
}

func (d *mockdns) CreateZone(ctx context.Context, param1 *dns.ZoneCreate, param2 dns.ZoneQueryString, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param1, param2, param3[0])
	} else {
		args = d.Called(ctx, param1, param2)
	}

	return args.Error(0)
}

func (d *mockdns) SaveChangelist(ctx context.Context, param *dns.ZoneCreate) error {
	args := d.Called(ctx, param)

	return args.Error(0)
}

func (d *mockdns) SubmitChangelist(ctx context.Context, param *dns.ZoneCreate) error {
	args := d.Called(ctx, param)

	return args.Error(0)
}

func (d *mockdns) UpdateZone(ctx context.Context, param1 *dns.ZoneCreate, param2 dns.ZoneQueryString) error {
	args := d.Called(ctx, param1, param2)

	return args.Error(0)
}

func (d *mockdns) DeleteZone(ctx context.Context, param1 *dns.ZoneCreate, param2 dns.ZoneQueryString) error {
	args := d.Called(ctx, param1, param2)

	return args.Error(0)
}

func (d *mockdns) ValidateZone(ctx context.Context, param1 *dns.ZoneCreate) error {
	args := d.Called(ctx, param1)

	return args.Error(0)
}

func (d *mockdns) GetZoneNames(ctx context.Context, param string) (*dns.ZoneNamesResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dns.ZoneNamesResponse), args.Error(1)
}

func (d *mockdns) GetZoneNameTypes(ctx context.Context, param1 string, param2 string) (*dns.ZoneNameTypesResponse, error) {
	args := d.Called(ctx, param1, param2)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dns.ZoneNameTypesResponse), args.Error(1)
}

func (d *mockdns) NewTsigKey(ctx context.Context, param string) *dns.TSIGKey {
	args := d.Called(ctx, param)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dns.TSIGKey)
}

func (d *mockdns) NewTsigQueryString(ctx context.Context) *dns.TSIGQueryString {
	args := d.Called(ctx)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dns.TSIGQueryString)
}

func (d *mockdns) ListTsigKeys(ctx context.Context, param *dns.TSIGQueryString) (*dns.TSIGReportResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.TSIGReportResponse), args.Error(1)
}

func (d *mockdns) GetTsigKeyZones(ctx context.Context, param *dns.TSIGKey) (*dns.ZoneNameListResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.ZoneNameListResponse), args.Error(1)
}

func (d *mockdns) GetTsigKeyAliases(ctx context.Context, param string) (*dns.ZoneNameListResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.ZoneNameListResponse), args.Error(1)
}

func (d *mockdns) TsigKeyBulkUpdate(ctx context.Context, param1 *dns.TSIGKeyBulkPost) error {
	args := d.Called(ctx, param1)

	return args.Error(0)
}

func (d *mockdns) GetTsigKey(ctx context.Context, param string) (*dns.TSIGKeyResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.TSIGKeyResponse), args.Error(1)
}

func (d *mockdns) DeleteTsigKey(ctx context.Context, param1 string) error {
	args := d.Called(ctx, param1)

	return args.Error(0)
}

func (d *mockdns) UpdateTsigKey(ctx context.Context, param1 *dns.TSIGKey, param2 string) error {
	args := d.Called(ctx, param1, param2)

	return args.Error(0)
}

func (d *mockdns) GetAuthorities(ctx context.Context, param string) (*dns.AuthorityResponse, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.AuthorityResponse), args.Error(1)
}

func (d *mockdns) GetNameServerRecordList(ctx context.Context, param string) ([]string, error) {
	args := d.Called(ctx, param)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (d *mockdns) NewAuthorityResponse(ctx context.Context, param string) *dns.AuthorityResponse {
	args := d.Called(ctx, param)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dns.AuthorityResponse)
}

func (d *mockdns) RecordToMap(ctx context.Context, param *dns.RecordBody) map[string]interface{} {
	args := d.Called(ctx, param)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]interface{})
}

func (d *mockdns) NewRecordBody(ctx context.Context, param dns.RecordBody) *dns.RecordBody {
	args := d.Called(ctx, param)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dns.RecordBody)
}

func (d *mockdns) GetRecordList(ctx context.Context, param string, param2 string, param3 string) (*dns.RecordSetResponse, error) {
	args := d.Called(ctx, param, param2, param3)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.RecordSetResponse), args.Error(1)
}

func (d *mockdns) GetRdata(ctx context.Context, param string, param2 string, param3 string) ([]string, error) {
	args := d.Called(ctx, param, param2, param3)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}

func (d *mockdns) ProcessRdata(ctx context.Context, param []string, param2 string) []string {
	args := d.Called(ctx, param, param2)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

func (d *mockdns) ParseRData(ctx context.Context, param string, param2 []string) map[string]interface{} {
	args := d.Called(ctx, param, param2)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]interface{})
}

func (d *mockdns) GetRecord(ctx context.Context, param string, param2 string, param3 string) (*dns.RecordBody, error) {
	args := d.Called(ctx, param, param2, param3)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.RecordBody), args.Error(1)
}

func (d *mockdns) CreateRecord(ctx context.Context, param *dns.RecordBody, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *mockdns) DeleteRecord(ctx context.Context, param *dns.RecordBody, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *mockdns) UpdateRecord(ctx context.Context, param *dns.RecordBody, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *mockdns) FullIPv6(ctx context.Context, param1 net.IP) string {
	args := d.Called(ctx, param1)

	return args.String(0)
}

func (d *mockdns) PadCoordinates(ctx context.Context, param1 string) string {
	args := d.Called(ctx, param1)

	return args.String(0)
}

func (d *mockdns) NewRecordSetResponse(ctx context.Context, param string) *dns.RecordSetResponse {
	args := d.Called(ctx, param)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dns.RecordSetResponse)
}

func (d *mockdns) GetRecordsets(ctx context.Context, param string, param2 ...dns.RecordsetQueryArgs) (*dns.RecordSetResponse, error) {
	var args mock.Arguments

	if len(param2) > 0 {
		args = d.Called(ctx, param, param2)
	} else {
		args = d.Called(ctx, param)
	}

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dns.RecordSetResponse), args.Error(1)
}

func (d *mockdns) CreateRecordsets(ctx context.Context, param *dns.Recordsets, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}

func (d *mockdns) UpdateRecordsets(ctx context.Context, param *dns.Recordsets, param2 string, param3 ...bool) error {
	var args mock.Arguments

	if len(param3) > 0 {
		args = d.Called(ctx, param, param2, param3)
	} else {
		args = d.Called(ctx, param, param2)
	}

	return args.Error(0)
}
