package appsec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceSelectedHostname() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSelectedHostnameCreate,
		ReadContext:   resourceSelectedHostnameRead,
		UpdateContext: resourceSelectedHostnameUpdate,
		DeleteContext: resourceSelectedHostnameDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"hostnames": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of hostnames to be added or removed from the protected hosts list",
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					Append,
					Replace,
					Remove,
				}, false)),
				Description: "How the hostnames are to be applied (APPEND, REMOVE or REPLACE)",
			},
		},
		DeprecationMessage: "This resource is deprecated with a scheduled end-of-life in v7.0.0 of our provider. Use the akamai_appsec_configuration instead.",
	}
}

func resourceSelectedHostnameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameCreate")
	logger.Debugf("in resourceSelectedHostnameCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames, err := tf.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}
	mode, err := tf.GetStringValue("mode", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// determine the actual hostname list to send to the API by combining the given hostnames & mode with the current hostnames
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
		ConfigID: configID,
		Version:  version,
	}
	getSelectdedHostnamesResponse, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'GetSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}
	currenthostnameset := schema.Set{F: schema.HashString}
	for _, h := range getSelectdedHostnamesResponse.HostnameList {
		currenthostnameset.Add(h.Hostname)
	}

	var desiredhostnameset *schema.Set
	switch mode {
	case Remove:
		desiredhostnameset = currenthostnameset.Difference(hostnames)
	case Append:
		desiredhostnameset = currenthostnameset.Union(hostnames)
	case Replace:
		desiredhostnameset = hostnames
	default:
		desiredhostnameset = hostnames
	}

	// convert to list of Hostname structs
	desiredhostnamelist := desiredhostnameset.List()
	newhostnames := make([]appsec.Hostname, 0, len(desiredhostnamelist))
	for _, h := range desiredhostnamelist {
		hostname := appsec.Hostname{
			Hostname: h.(string),
		}
		newhostnames = append(newhostnames, hostname)
	}

	version, err = getModifiableConfigVersion(ctx, configID, "selectedHostname", m)
	if err != nil {
		return diag.FromErr(err)
	}
	updateSelectedHostnames := appsec.UpdateSelectedHostnamesRequest{
		ConfigID:     configID,
		Version:      version,
		HostnameList: newhostnames,
	}

	_, err = client.UpdateSelectedHostnames(ctx, updateSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'UpdateSelectedHostnames': %s", err.Error())
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	// normally we don't set any attributes of the resource in Create, but for this resource we're not using the
	// supplied hostnames as is, rather we're combining them with the existing hostnames according to the value of mode
	if err := d.Set("hostnames", desiredhostnameset.List()); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", configID))

	return resourceSelectedHostnameRead(ctx, d, m)
}

func resourceSelectedHostnameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameRead")
	logger.Debugf("in resourceSelectedHostnameRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
		ConfigID: configID,
		Version:  version,
	}
	getSelectedHostnamesResponse, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'getSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	selectedhostnameset := schema.Set{F: schema.HashString}
	for _, hostname := range getSelectedHostnamesResponse.HostnameList {
		selectedhostnameset.Add(hostname.Hostname)
	}
	if err := d.Set("hostnames", selectedhostnameset.List()); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	// mode is not returned by API, so synthesize an appropriate value if we have none
	if _, ok := d.GetOk("mode"); !ok {
		if err := d.Set("mode", "REPLACE"); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	return nil
}

func resourceSelectedHostnameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameUpdate")
	logger.Debugf("in resourceSelectedHostnameUpdate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	mode, err := tf.GetStringValue("mode", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames, err := tf.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// determine the actual hostname list to send to the API by combining the given hostnames & mode with the current hostnames
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
		ConfigID: configID,
		Version:  version,
	}
	getSelectedHostnamesResponse, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'GetSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}
	currenthostnameset := schema.Set{F: schema.HashString}
	for _, h := range getSelectedHostnamesResponse.HostnameList {
		currenthostnameset.Add(h.Hostname)
	}

	var desiredhostnameset *schema.Set
	switch mode {
	case Remove:
		// implementing set difference manually here, as SDK's Set.Difference() doesn't seem to
		// give the correct result (elements of right-hand set not removed from left-hand set?)
		// desiredhostnameset = currenthostnameset.Difference(hostnames)
		desiredhostnameset = &schema.Set{F: currenthostnameset.F}
		hostnamelist := hostnames.List()
		for _, h := range currenthostnameset.List() {
			found := false
			for _, h2 := range hostnamelist {
				if h == h2 {
					found = true
					break
				}
			}
			if !found {
				desiredhostnameset.Add(h)
			}
		}
	case Append:
		desiredhostnameset = currenthostnameset.Union(hostnames)
	case Replace:
		desiredhostnameset = hostnames
	default:
		desiredhostnameset = hostnames
	}

	// convert to list of Hostname structs
	desiredhostnamelist := desiredhostnameset.List()
	newhostnames := make([]appsec.Hostname, 0, len(desiredhostnamelist))
	for _, h := range desiredhostnamelist {
		hostname := appsec.Hostname{
			Hostname: h.(string),
		}
		newhostnames = append(newhostnames, hostname)
	}

	version, err = getModifiableConfigVersion(ctx, configID, "selectedHostname", m)
	if err != nil {
		return diag.FromErr(err)
	}
	updateSelectedHostnames := appsec.UpdateSelectedHostnamesRequest{
		ConfigID:     configID,
		Version:      version,
		HostnameList: newhostnames,
	}

	_, err = client.UpdateSelectedHostnames(ctx, updateSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'UpdateSelectedHostnames': %s", err.Error())
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return resourceSelectedHostnameRead(ctx, d, m)
}

func resourceSelectedHostnameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(ctx, d, m)
}

// Append Replace Remove mode flags
const (
	Append  = "APPEND"
	Replace = "REPLACE"
	Remove  = "REMOVE"
)
