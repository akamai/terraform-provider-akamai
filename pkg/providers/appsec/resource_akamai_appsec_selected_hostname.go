package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSelectedHostname() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSelectedHostnameUpdate,
		ReadContext:   resourceSelectedHostnameRead,
		UpdateContext: resourceSelectedHostnameUpdate,
		DeleteContext: resourceSelectedHostnameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hostnames": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Append,
					Replace,
					Remove,
				}, false),
			},
		},
	}
}

func resourceSelectedHostnameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameRead")

	getSelectedHostname := appsec.GetSelectedHostnameRequest{}

	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")
		getSelectedHostname.ConfigID, _ = strconv.Atoi(s[0])
		getSelectedHostname.Version, _ = strconv.Atoi(s[1])
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSelectedHostname.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSelectedHostname.Version = version

	}

	selectedhostname, err := client.GetSelectedHostname(ctx, getSelectedHostname)
	if err != nil {
		logger.Errorf("calling 'getSelectedHostname': %s", err.Error())
		return diag.FromErr(err)
	}

	newhdata := make([]string, 0, len(selectedhostname.HostnameList))
	for _, hosts := range selectedhostname.HostnameList {
		newhdata = append(newhdata, hosts.Hostname)
	}

	d.Set("hostnames", newhdata)
	d.Set("config_id", getSelectedHostname.ConfigID)
	d.Set("version", getSelectedHostname.Version)
	d.SetId(fmt.Sprintf("%d:%d", getSelectedHostname.ConfigID, getSelectedHostname.Version))

	return nil
}

func resourceSelectedHostnameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceSelectedHostnameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnameUpdate")

	updateSelectedHostname := appsec.UpdateSelectedHostnameRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSelectedHostname.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSelectedHostname.Version = version

	mode := d.Get("mode").(string)

	hn := appsec.GetSelectedHostnamesRequest{}

	hostnamelist := d.Get("hostnames").([]interface{})

	for _, h := range hostnamelist {
		h1 := appsec.Hostname{}
		h1.Hostname = h.(string)
		hn.HostnameList = append(hn.HostnameList, h1)
	}

	getSelectedHostnames := appsec.GetSelectedHostnamesRequest{}
	getSelectedHostnames.ConfigID = configid
	getSelectedHostnames.Version = version

	selectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'getSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}

	switch mode {
	case Remove:
		for idx, h := range selectedhostnames.HostnameList {

			for _, hl := range hostnamelist {
				if h.Hostname == hl.(string) {
					RemoveIndex(selectedhostnames.HostnameList, idx)
				}
			}
		}
	case Append:
		for _, h := range selectedhostnames.HostnameList {
			m := appsec.Hostname{}
			m.Hostname = h.Hostname
			hn.HostnameList = append(hn.HostnameList, m)
		}
		selectedhostnames.HostnameList = hn.HostnameList
	case Replace:
		selectedhostnames.HostnameList = hn.HostnameList
	default:
		selectedhostnames.HostnameList = hn.HostnameList
	}

	updateSelectedHostname.HostnameList = selectedhostnames.HostnameList

	_, erru := client.UpdateSelectedHostname(ctx, updateSelectedHostname)
	if erru != nil {
		logger.Errorf("calling 'updateSelectedHostname': %s", erru.Error())
	}

	return resourceSelectedHostnameRead(ctx, d, m)
}

//RemoveIndex reemove host from list
func RemoveIndex(hl []appsec.Hostname, index int) []appsec.Hostname {
	return append(hl[:index], hl[index+1:]...)
}

// Append Replace Remove mode flags
const (
	Append  = "APPEND"
	Replace = "REPLACE"
	Remove  = "REMOVE"
)
