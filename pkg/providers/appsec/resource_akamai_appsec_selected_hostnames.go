package appsec

import (
	"fmt"
	"strconv"
	"strings"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		Create: resourceSelectedHostnamesUpdate,
		Read:   resourceSelectedHostnamesRead,
		Update: resourceSelectedHostnamesUpdate,
		Delete: resourceSelectedHostnamesDelete,
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

func resourceSelectedHostnamesRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSelectedHostnamesRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SelectedHostnames")

	selectedhostnames := appsec.NewSelectedHostnamesResponse()

	var configid int
	var version int
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Read SelectedHostnames D.ID %v", d.Id()))
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")
		configid, _ = strconv.Atoi(s[0])
		version, _ = strconv.Atoi(s[1])
	} else {
		configid = d.Get("config_id").(int)
		version = d.Get("version").(int)
	}

	err := selectedhostnames.GetSelectedHostnames(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	newhdata := make([]string, 0, len(selectedhostnames.HostnameList))
	for _, hosts := range selectedhostnames.HostnameList {
		newhdata = append(newhdata, hosts.Hostname)
	}

	//edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  SET SelectedHostnames H %v", h))
	d.Set("hostnames", newhdata)
	d.Set("config_id", configid)
	d.Set("version", version)
	d.SetId(fmt.Sprintf("%d:%d", configid, version))
	return nil
}

func resourceSelectedHostnamesDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSelectedHostnamesDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting SelectedHostnames")

	return schema.Noop(d, meta)
}

func resourceSelectedHostnamesUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSelectedHostnamesUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating SelectedHostnames")

	selectedhostnames := appsec.NewSelectedHostnamesResponse()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	mode := d.Get("mode").(string)

	hn := &appsec.SelectedHostnamesResponse{}

	hostnamelist := d.Get("hostnames").([]interface{})

	for _, h := range hostnamelist {
		m := appsec.Hostname{}
		m.Hostname = h.(string)
		hn.HostnameList = append(hn.HostnameList, m)
	}

	//Fill in existing then decide what to do
	err := selectedhostnames.GetSelectedHostnames(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
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

	err = selectedhostnames.UpdateSelectedHostnames(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	return resourceSelectedHostnamesRead(d, meta)

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
