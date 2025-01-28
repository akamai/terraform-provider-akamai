package property

import (
	"context"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
