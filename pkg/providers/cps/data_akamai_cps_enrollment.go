package cps

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	cpstools "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPSEnrollment() *schema.Resource {
	return &schema.Resource{
		Description: "Get an enrollment for given EnrollmentID",
		ReadContext: dataCPSEnrollmentRead,
		Schema: map[string]*schema.Schema{
			"enrollment_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"common_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sans": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"secure_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sni_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"admin_contact": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     contact,
			},
			"certificate_chain_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csr": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organizational_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_multi_stacked_certificates": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_configuration": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_mutual_authentication": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_ca_list_to_client": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"ocsp_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"set_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"disallowed_tls_versions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"clone_dns_names": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"geography": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"must_have_ciphers": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ocsp_stapling": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"preferred_ciphers": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quic_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tech_contact": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     contact,
			},
			"organization": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_line_one": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_line_two": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"postal_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"contract_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registration_authority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_challenges": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: cpstools.HashFromChallengesMap,
			},
			"http_challenges": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: cpstools.HashFromChallengesMap,
			},
		},
	}
}

func dataCPSEnrollmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "dataCPSEnrollmentRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching an enrollment")

	enrollmentID, err := tools.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	req := cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}
	enrollment, err := client.GetEnrollment(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not get an enrollment")
		return diag.FromErr(err)
	}

	attrs := createAttrs(enrollment, enrollmentID)

	if err = tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	challengesAttrs, err := getChallengesAttrs(ctx, enrollment, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = tools.SetAttrs(d, challengesAttrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(enrollmentID))
	return nil
}
