package cps

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	cpstools "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/cps/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPSEnrollments() *schema.Resource {
	return &schema.Resource{
		Description: "Get enrollments for given ContractID",
		ReadContext: dataCPSEnrollmentsRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contract ID for which enrollments are retrieved",
			},
			"enrollments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enrollment_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of enrollment",
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Common name used for enrollment",
						},
						"sans": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "List of SANs",
						},
						"secure_network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of TLS deployment network",
						},
						"sni_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Server Name Indication is used for enrollment",
						},
						"admin_contact": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        contact,
							Description: "Contact information for the certificate administrator used at organization",
						},
						"certificate_chain_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate trust chain type",
						},
						"csr": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Data used for generation of Certificate Signing Request",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"country_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The code of the country where organization is located",
									},
									"city": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "City where organization is located",
									},
									"organization": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of organization used in all legal documents",
									},
									"organizational_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Organizational unit of organization",
									},
									"preferred_trust_chain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "For the Let's Encrypt Domain Validated (DV) SAN certificates, the preferred trust chain will be included by CPS with the leaf certificate in the TLS handshake. If the field does not have a value, whichever trust chain Akamai chooses will be used by default",
									},
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "State or province of organization location",
									},
								},
							},
						},
						"enable_multi_stacked_certificates": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable Dual-Stacked certificate deployment for enrollment",
						},
						"network_configuration": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Settings containing network information and TLS metadata used by CPS",
							Elem: &schema.Resource{

								Schema: map[string]*schema.Schema{
									"client_mutual_authentication": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "The trust chain configuration used for client mutual authentication",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"send_ca_list_to_client": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Enable the server to send the certificate authority (CA) list to the client",
												},
												"ocsp_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Enable the OCSP stapling",
												},
												"set_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The identifier of the set of trust chains, created in the Trust Chain Manager",
												},
											},
										},
									},
									"disallowed_tls_versions": {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "TLS versions which are disallowed",
									},
									"clone_dns_names": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Enable CPS to direct traffic using all the SANs listed in the SANs parameter when enrollment is created",
									},
									"geography": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Geography type used for enrollment",
									},
									"must_have_ciphers": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mandatory Ciphers which are included for enrollment",
									},
									"ocsp_stapling": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enable OCSP stapling",
									},
									"preferred_ciphers": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Preferred Ciphers which are included for enrollment",
									},
									"quic_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Enable QUIC protocol",
									},
								},
							},
						},
						"signature_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SHA algorithm type",
						},
						"tech_contact": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        contact,
							Description: "Contact information for an administrator at Akamai",
						},
						"organization": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Organization information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of organization",
									},
									"phone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The phone number of the administrator who is organization contact",
									},
									"address_line_one": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The address of organization",
									},
									"address_line_two": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The address of organization",
									},
									"city": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "City where organization is located",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The region where organization resides",
									},
									"postal_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The postal code of organization",
									},
									"country_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Country code of country where organization is located",
									},
								},
							},
						},
						"certificate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate type of enrollment",
						},
						"validation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enrolment validation type",
						},
						"registration_authority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The registration authority or certificate authority (CA) used to obtain a certificate",
						},
						"pending_changes": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether some changes are pending",
						},
					},
				},
			},
		},
	}
}

func dataCPSEnrollmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "dataCPSEnrollmentsRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching enrollments")

	contractID, err := tf.GetStringValue("contract_id", d)
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

	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("akamai_cps_enrollments")
	return nil
}
