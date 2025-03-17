package property

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var listActivePropertyHostnamesResultsPerPage = 50

func dataSourcePropertyHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyHostnamesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tf.IsNotBlank,
			},
			"contract_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tf.IsNotBlank,
			},
			"property_id": {
				Type:             schema.TypeString,
				Required:         true,
				StateFunc:        addPrefixToState("prp_"),
				ValidateDiagFunc: tf.IsNotBlank,
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Version of property to fetch hostnames for. If not provided, 'latest' is used",
			},
			"hostnames": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of hostnames",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cname_type":             {Type: schema.TypeString, Computed: true},
						"edge_hostname_id":       {Type: schema.TypeString, Computed: true},
						"cname_from":             {Type: schema.TypeString, Computed: true},
						"cname_to":               {Type: schema.TypeString, Computed: true},
						"cert_provisioning_type": {Type: schema.TypeString, Computed: true},
						"cert_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     certStatus,
						},
					},
				},
			},
			"hostname_bucket": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of hostnames for property of type HOSTNAME_BUCKET",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cname_from": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hostname that your end users see, indicated by the Host header in end user requests.",
						},
						"cname_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Only one supported EDGE_HOSTNAME value.",
						},
						"staging_edge_hostname_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies each edge hostname.",
						},
						"staging_cert_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the certificate's provisioning type. Either CPS_MANAGED type for the certificates you create with the Certificate Provisioning System API (CPS), or DEFAULT for the Default Domain Validation (DV) certificates created automatically. Note that you can't specify the DEFAULT value if your property hostname uses the akamaized.net domain suffix.",
						},
						"staging_cname_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. This member corresponds to the edge hostname object's edgeHostnameDomain member.",
						},
						"production_edge_hostname_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies each edge hostname.",
						},
						"production_cert_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the certificate's provisioning type. Either CPS_MANAGED type for the certificates you create with the Certificate Provisioning System API (CPS), or DEFAULT for the Default Domain Validation (DV) certificates created automatically. Note that you can't specify the DEFAULT value if your property hostname uses the akamaized.net domain suffix.",
						},
						"production_cname_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. This member corresponds to the edge hostname object's edgeHostnameDomain member.",
						},
						"cert_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     certStatus,
						},
					},
				},
			},
			"filter_pending_default_certs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow to include `DEFAULT` cert types that have staging or production in a `PENDING` state. Default is false.",
			},
		},
	}
}

func dataPropertyHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := Client(meta)
	log := meta.Log("PAPI", "dataPropertyHostnamesRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	log.Debug("Listing Property Hostnames")

	// groupID / contractID is string as per schema.
	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID = str.AddPrefix(groupID, "grp_")
	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = str.AddPrefix(contractID, "ctr_")

	propertyID, err := tf.GetStringValue("property_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	propertyID = str.AddPrefix(propertyID, "prp_")

	version, err := tf.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	filterCerts, err := tf.GetBoolValue("filter_pending_default_certs", d)
	if err != nil {
		return diag.FromErr(err)
	}

	property, err := client.GetProperty(ctx, papi.GetPropertyRequest{ContractID: contractID, GroupID: groupID, PropertyID: propertyID})
	if err != nil {
		log.Error("could not fetch property", "error", err)
		return diag.FromErr(err)
	}

	if property.Property.PropertyType != nil && *property.Property.PropertyType == "HOSTNAME_BUCKET" {
		var diags diag.Diagnostics
		if version != 0 {
			diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: "provided `version` for HOSTNAME_BUCKET property, ignoring provided value"})
		}
		hostnames, err := getAllActivePropertyHostnames(ctx, client, contractID, groupID, propertyID, filterCerts)
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}

		d.SetId(propertyID)

		if err := d.Set("hostname_bucket", flattenBucketHostnames(hostnames)); err != nil {
			return append(diags, diag.Errorf("error setting hostnames: %s", err)...)
		}

		return diags
	}

	var prpVersion *papi.GetPropertyVersionsResponse
	if version == 0 {
		prpVersion, err = client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
			PropertyID: propertyID,
			ContractID: contractID,
			GroupID:    groupID,
		})
	} else {
		prpVersion, err = client.GetPropertyVersion(ctx, papi.GetPropertyVersionRequest{
			PropertyID:      propertyID,
			PropertyVersion: version,
			ContractID:      contractID,
			GroupID:         groupID,
		})
	}
	if err != nil {
		return diag.FromErr(err)
	}

	version = prpVersion.Version.PropertyVersion
	contractID = prpVersion.ContractID
	groupID = prpVersion.GroupID

	if err = d.Set("version", version); err != nil {
		return diag.FromErr(err)
	}

	hostNamesReq := papi.GetPropertyVersionHostnamesRequest{
		PropertyID:        propertyID,
		GroupID:           groupID,
		ContractID:        contractID,
		PropertyVersion:   version,
		IncludeCertStatus: true,
	}

	log.Debug("fetching property hostnames")
	hostnamesResponse, err := client.GetPropertyVersionHostnames(ctx, hostNamesReq)
	if err != nil {
		log.Error("could not fetch property hostnames", "error", err)
		return diag.FromErr(err)
	}

	log.Debug("fetched property hostnames", logFields(*hostnamesResponse))

	// setting concatenated id to uniquely identify data
	d.SetId(propertyID + strconv.Itoa(version))

	if err := d.Set("hostnames", flattenHostnames(hostnamesResponse.Hostnames.Items)); err != nil {
		return diag.Errorf("error setting hostnames: %s", err)
	}

	return nil

}

func getAllActivePropertyHostnames(ctx context.Context, client papi.PAPI, contractID, groupID, propertyID string, filterCerts bool) ([]papi.HostnameItem, error) {
	offset := 0
	returnedResults := listActivePropertyHostnamesResultsPerPage
	var allHostnames []papi.HostnameItem

	for returnedResults == listActivePropertyHostnamesResultsPerPage {
		hostnames, err := client.ListActivePropertyHostnames(ctx, papi.ListActivePropertyHostnamesRequest{
			ContractID:        contractID,
			GroupID:           groupID,
			PropertyID:        propertyID,
			Limit:             listActivePropertyHostnamesResultsPerPage,
			Offset:            offset,
			IncludeCertStatus: true,
		})
		if err != nil {
			return nil, fmt.Errorf("error fetching property hostnames: %w", err)
		}

		if filterCerts {
			for _, item := range hostnames.Hostnames.Items {
				prodPending := item.ProductionCertType == "DEFAULT" && item.CertStatus.Production[0].Status == "PENDING"
				stagingPending := item.StagingCertType == "DEFAULT" && item.CertStatus.Staging[0].Status == "PENDING"

				if prodPending || stagingPending {
					allHostnames = append(allHostnames, item)
				}
			}
		} else {
			allHostnames = append(allHostnames, hostnames.Hostnames.Items...)
		}

		returnedResults = len(hostnames.Hostnames.Items)
		offset += listActivePropertyHostnamesResultsPerPage
	}

	return allHostnames, nil
}
