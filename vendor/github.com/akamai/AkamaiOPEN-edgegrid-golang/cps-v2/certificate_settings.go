package cps

type RegistrationAuthority string
type CertificateType string
type AkamaiCipher string
type NetworkType string
type OCSPSetting string
type TLSType string
type SHA string
type ValidationType string

const (
	LetsEncryptRA          RegistrationAuthority = "letsencrypt"
	SymantecRA             RegistrationAuthority = "symantec"
	ThirdPartyRA           RegistrationAuthority = "third-party"
	SanCertificate         CertificateType       = "san"
	SymantecCertificate    CertificateType       = "single"
	WildCardCertificate    CertificateType       = "wildcard"
	WildCardSanCertificate CertificateType       = "wildcard-san"
	ThirdPartyCertificate  CertificateType       = "third-party"
	AK2018Q3               AkamaiCipher          = "ak-akamai-2018q3"
	AK2017Q3               AkamaiCipher          = "ak-akamai-default-2017q3"
	AK2016Q3               AkamaiCipher          = "ak-akamai-default-2016q3"
	AKPCIDSS               AkamaiCipher          = "ak-pci-dss-3.2"
	AKDefault              AkamaiCipher          = "ak-akamai-default"
	AK2016Q1               AkamaiCipher          = "ak-akamai-default-2016q1"
	AKPFSSupported         AkamaiCipher          = "ak-akamai-pfs-supported"
	AKPFS                  AkamaiCipher          = "ak-akamai-pfs"
	AKRecommended          AkamaiCipher          = "ak-akamai-recommended"
	AKSoftErrors           AkamaiCipher          = "ak-soft-errors"
	AKSoftErrorsWithExport AkamaiCipher          = "ak-soft-errors-with-export"
	AKTLS                  AkamaiCipher          = "ak-akamai-tls-1.2"
	AKPCIDSSDefault        AkamaiCipher          = "ak-pci-dss"
	AKPCIDSS3              AkamaiCipher          = "ak-pci-dss-3.1"
	StandardWorldWide      NetworkType           = "standard-worldwide"
	WorldWideRussia        NetworkType           = "worldwide-russia"
	WorldWide              NetworkType           = "worldwide"
	StandardTLS            TLSType               = "standard-tls"
	EnhancedTLS            TLSType               = "enhanced-tls"
	SHA1                   SHA                   = "SHA-1"
	SHA256                 SHA                   = "SHA-256"
	DomainValidation       ValidationType        = "dv"
	OrganizationValidation ValidationType        = "ov"
	ExtendedValidation     ValidationType        = "ev"
	ThirdPartyValidation   ValidationType        = "third-party"
)

type Contact struct {
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Title          *string `json:"title"`
	Organization   *string `json:"organizationName"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	AddressLineOne *string `json:"addressLineOne"`
	AddressLineTwo *string `json:"addressLineTwo"`
	City           *string `json:"city"`
	Region         *string `json:"region"`
	PostalCode     *string `json:"postalCode"`
	Country        *string `json:"country"`
}

type Organization struct {
	Name           *string `json:"name"`
	Phone          *string `json:"phone"`
	AddressLineOne *string `json:"addressLineOne"`
	AddressLineTwo *string `json:"addressLineTwo"`
	City           *string `json:"city"`
	Region         *string `json:"region"`
	PostalCode     *string `json:"postalCode"`
	Country        *string `json:"country"`
}

type CSR struct {
	CommonName         string    `json:"cn"`
	AlternativeNames   *[]string `json:"sans"`
	City               *string   `json:"l"`
	State              *string   `json:"st"`
	CountryCode        *string   `json:"c"`
	Organization       *string   `json:"o"`
	OrganizationalUnit *string   `json:"ou"`
}

type DomainNameSettings struct {
	CloneDomainNames bool      `json:"cloneDnsNames"`
	DomainNames      *[]string `json:"dnsNames"`
}

type NetworkConfiguration struct {
	DisallowedTLSVersions *[]string           `json:"disallowedTlsVersions"`
	DomainNameSettings    *DomainNameSettings `json:"dnsNameSettings"`
	Geography             string              `json:"geography"`
	MustHaveCiphers       AkamaiCipher        `json:"mustHaveCiphers"`
	// NetworkType           *NetworkType        `json:"networkType"`
	OCSPStapling     *OCSPSetting `json:"ocspStapling"`
	PreferredCiphers AkamaiCipher `json:"preferredCiphers"`
	QUICEnabled      bool         `json:"quicEnabled"`
	SecureNetwork    TLSType      `json:"secureNetwork"`
	// ServerNameIndication  *DomainNameSettings `json:"sni"`
	SNIOnly bool `json:"sniOnly"`
}

type ThirdParty struct {
	ExcludeSANS bool `json:"excludeSans"`
}
