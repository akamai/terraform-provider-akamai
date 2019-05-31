package akamai

import (
	cps "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceEnrollment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnrollmentCreate,
		Read:   resourceEnrollmentStub,
		Update: resourceEnrollmentStub,
		Delete: resourceEnrollmentStub,
		Exists: resourceEnrollmentExists,
		Schema: akamaiEnrollmentSchema,
	}
}

var akamaiEnrollmentSchema = map[string]*schema.Schema{
	"contract": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"deploy_not_after": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"deploy_not_before": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"admin_contact": {
		Type:     schema.TypeSet,
		Required: true,
		Elem:     resourceCPSContact(),
	},
	"certificate_chain_type": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"certificate_type": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"change_management": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false, // VERIFY
	},
	"csr": {
		Type:     schema.TypeSet,
		Required: true, // VERIFY
		Elem:     resourceCPSCSR(),
	},
	"enable_multi_stacked_certificates": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false, // VERIFY
	},
	"max_allowed_san_names": &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	},
	"max_allowed_wildcard_san_names": &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	},
	"network_configuration": {
		Type:     schema.TypeSet,
		Required: true,
		Elem:     resourceCPSNetworkConfiguration(),
	},
	"org": {
		Type:     schema.TypeSet,
		Required: true,
		Elem:     resourceCPSOrganization(),
	},
	"pending_changes": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"ra": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"signature_algorithm": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"tech_contact": &schema.Schema{
		Type:     schema.TypeSet,
		Required: true, // VERIFY
		Elem:     resourceCPSContact(),
	},
	"third_party": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true, // VERIFY
		Elem:     resourceCPSThirdParty(),
	},
	"validation_type": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"location": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
}

func resourceEnrollmentCreate(d *schema.ResourceData, meta interface{}) error {
	params := unmarshalCreateEnrollmentParams(d)
	enrollment := unmarshalEnrollment(d)

	response, err := enrollment.Create(*params)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] enrollmentCreate response %+v", response)

	d.SetId(response.Location)

	return nil
}

func resourceEnrollmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	enrollment, err := cps.GetEnrollment(d.Id())

	if err != nil || enrollment == nil {
		log.Printf("[DEBUG] Enrollment with location doesn't exist: %+v", d.Id())
		return false, err
	}

	log.Printf("[DEBUG] Enrollment found: %+v", enrollment)

	return true, nil
}

func resourceEnrollmentRead(d *schema.ResourceData, meta interface{}) error {
	enrollment, err := cps.GetEnrollment(d.Get("location").(string))

	if err != nil || enrollment == nil {
		log.Printf("[DEBUG] Enrollment with location doesn't exist: %+v", d.Id())
		return  err
	}

	log.Printf("[DEBUG] Enrollment found: %+v", enrollment)
	d.SetId(*enrollment.Location)
	return nil
}

func resourceEnrollmentStub(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func unmarshalEnrollment(d *schema.ResourceData) *cps.Enrollment {
	enrollment := &cps.Enrollment{
		AdminContact: unmarshalCPSContact(
			getSingleSchemaSetItem(d.Get("admin_contact")),
		),
		CertificateChainType: readNullableString(d.Get("certificate_chain_type")),
		CertificateType:      cps.CertificateType(d.Get("certificate_type").(string)),
		CertificateSigningRequest: unmarshalCPSCSR(
			getSingleSchemaSetItem(d.Get("csr")),
		),
		NetworkConfiguration: unmarshalCPSNetworkConfiguration(
			getSingleSchemaSetItem(d.Get("network_configuration")),
		),
		Organization: unmarshalCPSOrganization(
			getSingleSchemaSetItem(d.Get("org")),
		),
		RegistrationAuthority: cps.RegistrationAuthority(d.Get("ra").(string)),
		TechContact: unmarshalCPSContact(
			getSingleSchemaSetItem(d.Get("tech_contact")),
		),
		ValidationType: cps.ValidationType(d.Get("validation_type").(string)),
	}

	if changeManagement, ok := d.Get("change_management").(bool); ok {
		enrollment.ChangeManagement = changeManagement
	}

	if enableMultiStacked, ok := d.Get("enable_multi_stacked_certificates").(bool); ok {
		enrollment.EnableMultiStacked = enableMultiStacked
	}

	/*
		if maxSans, ok := d.Get("max_allowed_san_names").(int); ok {
			enrollment.MaxAllowedSans = maxSans
		}

		if maxWildcardSans, ok := d.Get("max_allowed_wildcard_san_names").(int); ok {
			enrollment.MaxAllowedWildcardSans = maxWildcardSans
		}
	*/
	if pendingChanges, ok := unmarshalSetString(d.Get("pending_changes")); ok {
		enrollment.PendingChanges = &pendingChanges
	}

	if shaString, ok := d.Get("signature_algorithm").(string); ok {
		sha := cps.SHA(shaString)
		enrollment.SignatureAuthority = &sha
	}

	if thirdParty, ok := d.GetOkExists("third_party"); ok {
		enrollment.ThirdParty = unmarshalCPSThirdParty(
			getSingleSchemaSetItem(thirdParty),
		)
	}

	return enrollment
}

func unmarshalCreateEnrollmentParams(d *schema.ResourceData) *cps.CreateEnrollmentQueryParams {
	params := &cps.CreateEnrollmentQueryParams{
		ContractID: d.Get("contract").(string),
	}

	if deployNotBefore, ok := d.Get("deploy_not_before").(string); ok {
		params.DeployNotBefore = &deployNotBefore
	}

	if deployNotAfter, ok := d.Get("deploy_not_after").(string); ok {
		params.DeployNotAfter = &deployNotAfter
	}

	return params
}

func unmarshalListEnrollmentsParams(d *schema.ResourceData) *cps.ListEnrollmentsQueryParams {
	return &cps.ListEnrollmentsQueryParams{
		ContractID: d.Get("contract").(string),
	}
}

func getCPSV2Service(d *schema.ResourceData) (*edgegrid.Config, error) {
	edgerc := d.Get("edgerc").(string)
	section := d.Get("cps_section").(string)

	config, err := edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	cps.Init(config)

	return &config, nil
}
