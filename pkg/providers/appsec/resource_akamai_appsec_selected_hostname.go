package appsec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
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
				Type:     schema.TypeInt,
				Required: true,
			},
			"hostnames": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					Append,
					Replace,
					Remove,
				}, false)),
			},
		},
	}
}

func resourceSelectedHostnameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameCreate")
	logger.Debugf("resourceSelectedHostnameCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames, err := tools.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}
	mode, err := tools.GetStringValue("mode", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// determine the actual hostname list to send to the API by combining the given hostnames & mode with the current hostnames

	getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
		ConfigID: configID,
		Version:  getLatestConfigVersion(ctx, configID, m),
	}

	currentselectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'GetSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}
	currenthostnameset := schema.Set{F: schema.HashString}
	for _, h := range currentselectedhostnames.HostnameList {
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

	updateSelectedHostnames := appsec.UpdateSelectedHostnamesRequest{
		ConfigID:     configID,
		Version:      getModifiableConfigVersion(ctx, configID, "selectedHostname", m),
		HostnameList: newhostnames,
	}

	_, err = client.UpdateSelectedHostnames(ctx, updateSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'UpdateSelectedHostnames': %s", err.Error())
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	// normally we don't set any attributes of the resource in Create, but for this resource we're not using the
	// supplied hostnames as is, rather we're combining them with the existing hostnames according to the value of mode
	if err := d.Set("hostnames", desiredhostnameset.List()); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", configID))

	return resourceSelectedHostnameRead(ctx, d, m)
}

func resourceSelectedHostnameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameRead")
	logger.Debugf("resourceSelectedHostnameRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
		ConfigID: configID,
		Version:  getLatestConfigVersion(ctx, configID, m),
	}

	selectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'getSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	selectedhostnameset := schema.Set{F: schema.HashString}
	for _, hostname := range selectedhostnames.HostnameList {
		selectedhostnameset.Add(hostname.Hostname)
	}
	if err := d.Set("hostnames", selectedhostnameset.List()); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	// mode is not returned by API, so synthesize an appropriate value if we have none
	if _, ok := d.GetOk("mode"); !ok {
		if err := d.Set("mode", "REPLACE"); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	return nil
}

func resourceSelectedHostnameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameUpdate")
	logger.Debugf("resourceSelectedHostnameUpdate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	mode, err := tools.GetStringValue("mode", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames, err := tools.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// determine the actual hostname list to send to the API by combining the given hostnames & mode with the current hostnames

	getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
		ConfigID: configID,
		Version:  getLatestConfigVersion(ctx, configID, m),
	}

	selectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'GetSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}
	currenthostnameset := schema.Set{F: schema.HashString}
	for _, h := range selectedhostnames.HostnameList {
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

	updateSelectedHostnames := appsec.UpdateSelectedHostnamesRequest{
		ConfigID:     configID,
		Version:      getModifiableConfigVersion(ctx, configID, "selectedHostname", m),
		HostnameList: newhostnames,
	}

	_, err = client.UpdateSelectedHostnames(ctx, updateSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'UpdateSelectedHostnames': %s", err.Error())
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
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
