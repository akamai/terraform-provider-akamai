package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/apex/log"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resPropHostname() *schema.Resource {
	return &schema.Resource{
		CreateContext: resPropHostnameCreate,
		ReadContext:   resPropHostnameRead,
		UpdateContext: resPropHostnameUpdate,
		DeleteContext: schema.NoopContext, // To delete this resource is to simply remove it from TF's stewardship
		Importer: &schema.ResourceImporter{
			StateContext: resPropHostnameImport,
		},
		Schema: map[string]*schema.Schema{
			"property_id": {Type: schema.TypeString, Required: true, StateFunc: addPrefixToState("prp_")},
			"group_id":    {Type: schema.TypeString, Required: true, StateFunc: addPrefixToState("grp_")},
			"contract_id": {Type: schema.TypeString, Required: true, StateFunc: addPrefixToState("ctr_")},
			"property_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The current latest version number of the associated property",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the associated property version has been activated in any network",
			},
			"names": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Mapping of CNAMEs to Edge Hostname IDs",
			},
		},
	}
}

func resPropHostnameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger := akamai.Meta(m).Log("PAPI", "resPropHostnameCreate")
	client := inst.Client(akamai.Meta(m))
	ctx = log.NewContext(ctx, logger)

	// Schema guarantees property_id, contract_id, and group_id are strings
	PropertyID := tools.AddPrefix(d.Get("property_id").(string), "prp_")
	ContractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	GroupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")

	// Schema guarantees "names" is a map[string]interface{}
	confNames := d.Get("names").(map[string]interface{})

	// Fetch the property to get the latest version
	Property, err := fetchProperty(ctx, client, PropertyID, GroupID, ContractID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fetch current hostnames to determine whether we can skip the first update
	Hostnames, err := fetchPropertyHostnames(ctx, client, *Property)
	if err != nil {
		return diag.FromErr(err)
	}

	// If actual host names are already the same as configured, we silently adopt them regardless of activation status
	if cmp.Equal(confNames, hostnamesToRDMap(Hostnames), cmpopts.EquateEmpty()) {
		d.SetId(PropertyID)
		return resPropHostnameRead(ctx, d, m)
	}

	if latestVersionIsActive(*Property) {
		VersionID, err := createPropertyVersion(ctx, client, *Property)
		if err != nil {
			return diag.FromErr(err)
		}

		Property.LatestVersion = VersionID
	}

	if err := updatePropertyHostnames(ctx, client, *Property, rdMapToHostnames(confNames)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(PropertyID)
	return resPropHostnameRead(ctx, d, m)
}

func resPropHostnameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := inst.Client(akamai.Meta(m))
	logger := akamai.Meta(m).Log("PAPI", "resPropHostnameRead")
	ctx = log.NewContext(ctx, logger)

	// Schema guarantees property_id, group_id, and contract_id are strings
	PropertyID := tools.AddPrefix(d.Get("property_id").(string), "prp_")
	ContractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	GroupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")

	Property, err := fetchProperty(ctx, client, PropertyID, GroupID, ContractID)
	if err != nil {
		return diag.FromErr(err)
	}

	Hostnames, err := fetchPropertyHostnames(ctx, client, *Property)
	if err != nil {
		return diag.FromErr(err)
	}

	// Now that we have everything, fill out the ResourceData
	attrs := map[string]interface{}{
		"names":            hostnamesToRDMap(Hostnames),
		"property_version": Property.LatestVersion,
		"is_active":        latestVersionIsActive(*Property),
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resPropHostnameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resPropHostnameUpdate"))
	client := inst.Client(akamai.Meta(m))
	logger := log.FromContext(ctx)

	if d.HasChange("property_id") {
		err := fmt.Errorf(`attribute "property_id" cannot be changed`)
		logger.WithError(err).Error("could not update property")
		return diag.FromErr(err)
	}

	// Schema guarantees "names" is a map[string]interface{}
	names := d.Get("names").(map[string]interface{})

	// Schema guarantess "is_active" is bool
	versionIsActive := d.Get("is_active").(bool)

	// Schema guarantees the types of these
	Property := papi.Property{
		PropertyID:    d.Get("property_id").(string),
		GroupID:       d.Get("group_id").(string),
		ContractID:    d.Get("contract_id").(string),
		LatestVersion: d.Get("property_version").(int),
	}

	if versionIsActive {
		VersionID, err := createPropertyVersion(ctx, client, Property)
		if err != nil {
			return diag.FromErr(err)
		}

		Property.LatestVersion = VersionID
	}

	if err := updatePropertyHostnames(ctx, client, Property, rdMapToHostnames(names)); err != nil {
		return diag.FromErr(err)
	}

	return resPropHostnameRead(ctx, d, m)
}

func resPropHostnameImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	logger := akamai.Meta(m).Log("PAPI", "resPropHostnameImport")
	ctx = log.NewContext(ctx, logger)

	// User-supplied import ID is a comma-separated list of PropertyID[,GroupID[,ContractID]]
	// ContractID and GroupID are optional as long as the PropertyID is sufficient to fetch the property
	var PropertyID, GroupID, ContractID string
	parts := strings.Split(d.Id(), ",")

	switch len(parts) {
	case 1:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
	case 2:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
		GroupID = tools.AddPrefix(parts[1], "grp_")
	case 3:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
		GroupID = tools.AddPrefix(parts[1], "grp_")
		ContractID = tools.AddPrefix(parts[2], "ctr_")
	default:
		return nil, fmt.Errorf("invalid property identifier: %q", d.Id())
	}

	// Import only needs to set the resource ID and enough attributes that the read opertaion will function, so there's
	// no need to fetch anything if the user gave both GroupID and ContractID
	if GroupID != "" && ContractID != "" {
		attrs := map[string]interface{}{
			"property_id": PropertyID,
			"group_id":    GroupID,
			"contract_id": ContractID,
		}
		if err := rdSetAttrs(ctx, d, attrs); err != nil {
			return nil, err
		}

		d.SetId(PropertyID)
		return []*schema.ResourceData{d}, nil
	}

	// Missing GroupID, ContractID, or both -- Attempt a fetch to get them. If the PropertyID is not sufficient, PAPI
	// will return an error.
	Property, err := fetchProperty(ctx, inst.Client(akamai.Meta(m)), PropertyID, GroupID, ContractID)
	if err != nil {
		return nil, err
	}

	attrs := map[string]interface{}{
		"property_id": Property.PropertyID,
		"group_id":    Property.GroupID,
		"contract_id": Property.ContractID,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return nil, err
	}

	d.SetId(Property.PropertyID)
	return []*schema.ResourceData{d}, nil
}

// Retrieves basic info for a Property
func fetchProperty(ctx context.Context, client papi.PAPI, PropertyID, GroupID, ContractID string) (*papi.Property, error) {
	req := papi.GetPropertyRequest{
		PropertyID: PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("fetching property")
	res, err := client.GetProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not read property")
		return nil, err
	}

	logger = logger.WithFields(logFields(*res))

	if res.Property == nil {
		err := fmt.Errorf("PAPI::GetProperty() response did not contain a property")
		logger.WithError(err).Error("could not look up property")
		return nil, err
	}

	logger.Debug("property fetched")
	return res.Property, nil
}

// Fetch hostnames for latest version of given property
func fetchPropertyHostnames(ctx context.Context, client papi.PAPI, Property papi.Property) ([]papi.Hostname, error) {
	req := papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      Property.PropertyID,
		GroupID:         Property.GroupID,
		ContractID:      Property.ContractID,
		PropertyVersion: Property.LatestVersion,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("fetching property hostnames")
	res, err := client.GetPropertyVersionHostnames(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not fetch property hostnames")
		return nil, err
	}

	logger.WithFields(logFields(*res)).Debug("fetched property hostnames")
	return res.Hostnames.Items, nil
}

// Create a new property version based on the latest version of the given property
func createPropertyVersion(ctx context.Context, client papi.PAPI, Property papi.Property) (NewVersion int, err error) {
	req := papi.CreatePropertyVersionRequest{
		PropertyID: Property.PropertyID,
		ContractID: Property.ContractID,
		GroupID:    Property.GroupID,
		Version: papi.PropertyVersionCreate{
			CreateFromVersion: Property.LatestVersion,
		},
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("creating new property version")
	res, err := client.CreatePropertyVersion(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not create new property version")
		return
	}

	logger.WithFields(logFields(*res)).Info("property version created")
	NewVersion = res.PropertyVersion
	return
}

// Set hostnames of the latest version of the given property
func updatePropertyHostnames(ctx context.Context, client papi.PAPI, Property papi.Property, Hostnames []papi.Hostname) error {
	req := papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      Property.PropertyID,
		GroupID:         Property.GroupID,
		ContractID:      Property.ContractID,
		PropertyVersion: Property.LatestVersion,
		Hostnames:       Hostnames,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("updating property hostnames")
	res, err := client.UpdatePropertyVersionHostnames(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not create new property version")
		return err
	}

	logger.WithFields(logFields(*res)).Info("property hostnames updated")
	return nil
}

// Convert given hostnames to the map form that can be stored in a schema.ResourceData
func hostnamesToRDMap(Hostnames []papi.Hostname) map[string]interface{} {
	m := map[string]interface{}{}
	for _, hn := range Hostnames {
		m[hn.CnameFrom] = hn.EdgeHostnameID
	}

	return m
}

// Convert the given map from a schema.ResourceData to a slice of papi.Hostnames
func rdMapToHostnames(given map[string]interface{}) []papi.Hostname {
	var Hostnames []papi.Hostname

	for cname, ehnID := range given {
		Hostnames = append(Hostnames, papi.Hostname{
			CnameType:      "EDGE_HOSTNAME",
			EdgeHostnameID: ehnID.(string), // guaranteed by schema to be a string
			CnameFrom:      cname,
		})
	}

	return Hostnames
}

// Tests true when the given property's latest version is active in any environment
func latestVersionIsActive(Property papi.Property) bool {
	// NB: The property version could possibly be in a transitional state where PAPI will reject a request to update
	//     hostnames or rules. It is the responsibility of the user to ensure the resource dependencies are such that
	//     changes happen in the correct order.
	ActiveInProd := Property.ProductionVersion != nil && *Property.ProductionVersion == Property.LatestVersion
	ActiveInStaging := Property.StagingVersion != nil && *Property.StagingVersion == Property.LatestVersion

	return ActiveInProd || ActiveInStaging
}

// Set many attributes of a schema.ResourceData in one call
func rdSetAttrs(ctx context.Context, d *schema.ResourceData, AttributeValues map[string]interface{}) error {
	logger := log.FromContext(ctx)

	for attr, value := range AttributeValues {
		if err := d.Set(attr, value); err != nil {
			logger.WithError(err).Errorf("could not set %q", attr)
			return err
		}
	}

	return nil
}
