package gtm

import (
	"context"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/stretchr/testify/mock"
)

type mockgtm struct {
	mock.Mock
}

func (p *mockgtm) NullFieldMap(ctx context.Context, _ *gtm.Domain) (*gtm.NullFieldMapStruct, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.NullFieldMapStruct), args.Error(1)
}

func (p *mockgtm) NewDomain(ctx context.Context, _ string, _ string) *gtm.Domain {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.Domain)
}

func (p *mockgtm) GetDomainStatus(ctx context.Context, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) ListDomains(ctx context.Context) ([]*gtm.DomainItem, error) {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*gtm.DomainItem), args.Error(1)
}

func (p *mockgtm) GetDomain(ctx context.Context, domain string) (*gtm.Domain, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.Domain), args.Error(1)
}

func (p *mockgtm) CreateDomain(ctx context.Context, domain *gtm.Domain, queryArgs map[string]string) (*gtm.DomainResponse, error) {
	args := p.Called(ctx, domain, queryArgs)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.DomainResponse), args.Error(1)
}

func (p *mockgtm) DeleteDomain(ctx context.Context, domain *gtm.Domain) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) UpdateDomain(ctx context.Context, domain *gtm.Domain, queryArgs map[string]string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, domain, queryArgs)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) GetProperty(ctx context.Context, prop string, domain string) (*gtm.Property, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.Property), args.Error(1)
}

func (p *mockgtm) DeleteProperty(ctx context.Context, prop *gtm.Property, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) CreateProperty(ctx context.Context, prop *gtm.Property, domain string) (*gtm.PropertyResponse, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.PropertyResponse), args.Error(1)
}

func (p *mockgtm) UpdateProperty(ctx context.Context, prop *gtm.Property, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, prop, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) NewTrafficTarget(ctx context.Context) *gtm.TrafficTarget {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.TrafficTarget)
}

func (p *mockgtm) NewStaticRRSet(ctx context.Context) *gtm.StaticRRSet {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.StaticRRSet)
}

func (p *mockgtm) NewLivenessTest(ctx context.Context, a string, b string, c int, d float32) *gtm.LivenessTest {
	args := p.Called(ctx, a, b, c, d)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.LivenessTest)
}

func (p *mockgtm) NewProperty(ctx context.Context, prop string) *gtm.Property {
	args := p.Called(ctx, prop)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.Property)
}

func (p *mockgtm) ListProperties(ctx context.Context, domain string) ([]*gtm.Property, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*gtm.Property), args.Error(1)
}

func (p *mockgtm) GetDatacenter(ctx context.Context, dcid int, domain string) (*gtm.Datacenter, error) {
	args := p.Called(ctx, dcid, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.Datacenter), args.Error(1)
}

func (p *mockgtm) CreateDatacenter(ctx context.Context, dc *gtm.Datacenter, domain string) (*gtm.DatacenterResponse, error) {
	args := p.Called(ctx, dc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.DatacenterResponse), args.Error(1)
}

func (p *mockgtm) DeleteDatacenter(ctx context.Context, dc *gtm.Datacenter, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, dc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) UpdateDatacenter(ctx context.Context, dc *gtm.Datacenter, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, dc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) NewDatacenterResponse(ctx context.Context) *gtm.DatacenterResponse {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.DatacenterResponse)
}

func (p *mockgtm) NewDatacenter(ctx context.Context) *gtm.Datacenter {
	args := p.Called(ctx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.Datacenter)
}

func (p *mockgtm) ListDatacenters(ctx context.Context, domain string) ([]*gtm.Datacenter, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*gtm.Datacenter), args.Error(1)
}

func (p *mockgtm) CreateIPv4DefaultDatacenter(ctx context.Context, domain string) (*gtm.Datacenter, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.Datacenter), args.Error(1)
}

func (p *mockgtm) CreateIPv6DefaultDatacenter(ctx context.Context, domain string) (*gtm.Datacenter, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.Datacenter), args.Error(1)
}

func (p *mockgtm) CreateMapsDefaultDatacenter(ctx context.Context, domainName string) (*gtm.Datacenter, error) {
	args := p.Called(ctx, domainName)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.Datacenter), args.Error(1)
}

func (p *mockgtm) GetResource(ctx context.Context, rsrc string, domain string) (*gtm.Resource, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.Resource), args.Error(1)
}

func (p *mockgtm) CreateResource(ctx context.Context, rsrc *gtm.Resource, domain string) (*gtm.ResourceResponse, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResourceResponse), args.Error(1)
}

func (p *mockgtm) DeleteResource(ctx context.Context, rsrc *gtm.Resource, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) UpdateResource(ctx context.Context, rsrc *gtm.Resource, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, rsrc, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) NewResourceInstance(ctx context.Context, ri *gtm.Resource, a int) *gtm.ResourceInstance {
	args := p.Called(ctx, ri, a)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.ResourceInstance)
}

func (p *mockgtm) NewResource(ctx context.Context, rname string) *gtm.Resource {
	args := p.Called(ctx, rname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.Resource)
}

func (p *mockgtm) ListResources(ctx context.Context, domain string) ([]*gtm.Resource, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*gtm.Resource), args.Error(1)
}

func (p *mockgtm) GetAsMap(ctx context.Context, asmap string, domain string) (*gtm.AsMap, error) {
	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.AsMap), args.Error(1)
}

func (p *mockgtm) CreateAsMap(ctx context.Context, asmap *gtm.AsMap, domain string) (*gtm.AsMapResponse, error) {
	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.AsMapResponse), args.Error(1)
}

func (p *mockgtm) DeleteAsMap(ctx context.Context, asmap *gtm.AsMap, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) UpdateAsMap(ctx context.Context, asmap *gtm.AsMap, domain string) (*gtm.ResponseStatus, error) {

	args := p.Called(ctx, asmap, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) NewAsMap(ctx context.Context, mname string) *gtm.AsMap {
	args := p.Called(ctx, mname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.AsMap)
}

func (p *mockgtm) NewASAssignment(ctx context.Context, as *gtm.AsMap, a int, b string) *gtm.AsAssignment {
	args := p.Called(ctx, as, a, b)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.AsAssignment)
}

func (p *mockgtm) ListAsMaps(ctx context.Context, domain string) ([]*gtm.AsMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*gtm.AsMap), args.Error(1)
}

func (p *mockgtm) GetGeoMap(ctx context.Context, geo string, domain string) (*gtm.GeoMap, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.GeoMap), args.Error(1)
}

func (p *mockgtm) CreateGeoMap(ctx context.Context, geo *gtm.GeoMap, domain string) (*gtm.GeoMapResponse, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.GeoMapResponse), args.Error(1)
}

func (p *mockgtm) DeleteGeoMap(ctx context.Context, geo *gtm.GeoMap, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) UpdateGeoMap(ctx context.Context, geo *gtm.GeoMap, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, geo, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) NewGeoMap(ctx context.Context, mname string) *gtm.GeoMap {
	args := p.Called(ctx, mname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.GeoMap)
}

func (p *mockgtm) NewGeoAssignment(ctx context.Context, as *gtm.GeoMap, a int, b string) *gtm.GeoAssignment {
	args := p.Called(ctx, as, a, b)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.GeoAssignment)
}

func (p *mockgtm) ListGeoMaps(ctx context.Context, domain string) ([]*gtm.GeoMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*gtm.GeoMap), args.Error(1)
}

func (p *mockgtm) GetCidrMap(ctx context.Context, cidr string, domain string) (*gtm.CidrMap, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.CidrMap), args.Error(1)
}

func (p *mockgtm) CreateCidrMap(ctx context.Context, cidr *gtm.CidrMap, domain string) (*gtm.CidrMapResponse, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.CidrMapResponse), args.Error(1)
}

func (p *mockgtm) DeleteCidrMap(ctx context.Context, cidr *gtm.CidrMap, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) UpdateCidrMap(ctx context.Context, cidr *gtm.CidrMap, domain string) (*gtm.ResponseStatus, error) {
	args := p.Called(ctx, cidr, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gtm.ResponseStatus), args.Error(1)
}

func (p *mockgtm) NewCidrMap(ctx context.Context, mname string) *gtm.CidrMap {
	args := p.Called(ctx, mname)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.CidrMap)
}

func (p *mockgtm) NewCidrAssignment(ctx context.Context, as *gtm.CidrMap, a int, b string) *gtm.CidrAssignment {
	args := p.Called(ctx, as, a, b)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*gtm.CidrAssignment)
}

func (p *mockgtm) ListCidrMaps(ctx context.Context, domain string) ([]*gtm.CidrMap, error) {
	args := p.Called(ctx, domain)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*gtm.CidrMap), args.Error(1)
}
