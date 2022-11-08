package cps

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	cpstools "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPSEnrollments() *schema.Resource {
	return &schema.Resource{
		Description: "Get enrollments for given ContractID",
		ReadContext: dataCPSEnrollmentsRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enrollments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enrollment_id": {
							Type:     schema.TypeInt,
							Computed: true,
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
						"pending_changes": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataCPSEnrollmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "dataCPSEnrollmentsRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching enrollments")

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	req := cps.ListEnrollmentsRequest{
		ContractID: contractID,
	}
	enrollments, err := client.ListEnrollments(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not get enrollments")
		return diag.FromErr(err)
	}

	enrollmentsAttrs := make([]interface{}, 0)

	for _, enrollment := range enrollments.Enrollments {
		enID, err := cpstools.GetEnrollmentID(enrollment.Location)
		if err != nil {
			return diag.FromErr(err)
		}
		attrs := createAttrs(&enrollment, enID)

		if len(enrollment.PendingChanges) > 0 {
			attrs["pending_changes"] = true
		} else {
			attrs["pending_changes"] = false
		}

		enrollmentsAttrs = append(enrollmentsAttrs, attrs)
	}

	attrs := map[string]interface{}{"enrollments": enrollmentsAttrs}

	if err = tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("akamai_cps_enrollments")
	return nil
}
