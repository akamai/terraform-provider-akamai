package cloudcertificates

import (
	"slices"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudcertificates"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// defaultPageSize is the default page size used for paginated API requests.
const defaultPageSize int64 = 100

// subjectModel represents X.509 certificate subject fields.
type subjectModel struct {
	CommonName   types.String `tfsdk:"common_name"`
	Organization types.String `tfsdk:"organization"`
	Country      types.String `tfsdk:"country"`
	State        types.String `tfsdk:"state"`
	Locality     types.String `tfsdk:"locality"`
}

// bindingModel represents a certificate-hostname binding.
type bindingModel struct {
	CertificateID types.String `tfsdk:"certificate_id"`
	Hostname      types.String `tfsdk:"hostname"`
	Network       types.String `tfsdk:"network"`
	ResourceType  types.String `tfsdk:"resource_type"`
}

// validKeyCombinations maps each supported key type to its valid key sizes.
var validKeyCombinations = map[cloudcertificates.CryptographicAlgorithm][]cloudcertificates.KeySize{
	cloudcertificates.CryptographicAlgorithmECDSA: {cloudcertificates.KeySizeP256, cloudcertificates.KeySizeP384},
	cloudcertificates.CryptographicAlgorithmRSA:   {cloudcertificates.KeySize2048},
}

// validGeoClassesByNetwork maps each supported secure network type to its valid geographic network classes.
var validGeoClassesByNetwork = map[cloudcertificates.SecureNetwork][]cloudcertificates.GeoClass{
	cloudcertificates.SecureNetworkEnhancedTLS: {
		cloudcertificates.GeoClassStandardWorldwide,
		cloudcertificates.GeoClassContiguousUS,
		cloudcertificates.GeoClassReservedGlobal,
	},
	cloudcertificates.SecureNetworkStandardTLS: {
		cloudcertificates.GeoClassStandardWorldwide,
	},
}

// validGeoClasses returns all unique valid geographic network classes across all secure network types.
func validGeoClasses() []cloudcertificates.GeoClass {
	var result []cloudcertificates.GeoClass
	for _, classes := range validGeoClassesByNetwork {
		result = append(result, classes...)
	}
	slices.Sort(result)
	return slices.Compact(result)
}

// validKeyTypes returns sorted valid key type names derived from ValidKeyCombinations.
func validKeyTypes() []cloudcertificates.CryptographicAlgorithm {
	var keys []cloudcertificates.CryptographicAlgorithm
	for k := range validKeyCombinations {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// validKeySizes returns sorted key sizes derived from ValidKeyCombinations.
func validKeySizes() []cloudcertificates.KeySize {
	var sizes []cloudcertificates.KeySize
	for _, t := range validKeyTypes() {
		sizes = append(sizes, validKeyCombinations[t]...)
	}
	slices.Sort(sizes)
	return sizes
}

// isEmptySubject returns true if all fields of a Subject are empty strings.
func isEmptySubject(subject cloudcertificates.Subject) bool {
	return subject.CommonName == "" && subject.Organization == "" && subject.Country == "" &&
		subject.State == "" && subject.Locality == ""
}
