package property

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
)

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
				Computed:    true,
				Description: "This is a computed value - provider will always use 'latest' version, providing own version number is not supported",
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
	groupID = tools.AddPrefix(groupID, "grp_")
	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")

	propertyID, err := tf.GetStringValue("property_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	propertyID = tools.AddPrefix(propertyID, "prp_")

	latestVersion, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
		PropertyID: propertyID,
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	version := latestVersion.Version.PropertyVersion
	contractID = latestVersion.ContractID
	groupID = latestVersion.GroupID

	if err := d.Set("version", version); err != nil {
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
		log.WithError(err).Error("could not fetch property hostnames")
		return diag.FromErr(err)
	}

	log.WithFields(logFields(*hostnamesResponse)).Debug("fetched property hostnames")

	// setting concatenated id to uniquely identify data
	d.SetId(propertyID + strconv.Itoa(version))

	if err := d.Set("hostnames", flattenHostnames(hostnamesResponse.Hostnames.Items)); err != nil {
		return diag.Errorf("error setting hostnames: %s", err)
	}

	return nil
}
