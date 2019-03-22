package akamai

import (
	cps "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCPSNetworkConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: cpsNetworkConfigurationSchema,
	}
}

func unmarshalCPSNetworkConfiguration(d map[string]interface{}) *cps.NetworkConfiguration {
	networkConfiguration := &cps.NetworkConfiguration{
		Geography: d["geography"].(string),
		// NetworkType:   cps.NetworkType(d["network_type"].(string)),
		SecureNetwork: cps.TLSType(d["secure_network"].(string)),
		DomainNameSettings: unmarshalCPSDomainNameSettings(
			getSingleSchemaSetItem(d["dns_name_settings"]),
		),
	}

	if disallowedTLS, ok := unmarshalSetString(d["disallowed_tls_version"]); ok {
		networkConfiguration.DisallowedTLSVersions = &disallowedTLS
	}

	if mustHaveCiphers, ok := d["must_have_ciphers"].(string); ok {
		networkConfiguration.MustHaveCiphers = cps.AkamaiCipher(mustHaveCiphers)
	}

	if ocspString := readNullableString(d["ocsp_stapling"]); ocspString != nil {
		ocspStapling := cps.OCSPSetting(*ocspString)
		networkConfiguration.OCSPStapling = &ocspStapling
	}

	if preferredCiphers, ok := d["preferred_ciphers"].(string); ok {
		networkConfiguration.PreferredCiphers = cps.AkamaiCipher(preferredCiphers)
	}

	if quicEnabled, ok := d["quic_enabled"].(bool); ok {
		networkConfiguration.QUICEnabled = quicEnabled
	}

	if sniOnly, ok := d["sni_only"].(bool); ok {
		networkConfiguration.SNIOnly = sniOnly
	}

	return networkConfiguration
}

var cpsNetworkConfigurationSchema = map[string]*schema.Schema{
	"disallowed_tls_version": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"dns_name_settings": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     resourceCPSDomainNameSettings(),
	},
	"geography": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"must_have_ciphers": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"network_type": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"ocsp_stapling": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"preferred_ciphers": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"quic_enabled": &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	},
	"secure_network": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}, /*
		"sni": {
			Type:     schema.TypeSet,
			O: true,
			Elem:     resourceCPSDomainNameSettings(),
		},*/
	"sni_only": &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
		ForceNew: true,
	},
}
