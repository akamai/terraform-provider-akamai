package networklists

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/hash"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// network_lists v2
//
// https://techdocs.akamai.com/network-lists/reference/api
func resourceNetworkListSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkListSubscriptionUpdate,
		ReadContext:   resourceNetworkListSubscriptionRead,
		UpdateContext: resourceNetworkListSubscriptionUpdate,
		DeleteContext: resourceNetworkListSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"recipients": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkListSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListSubscriptionRead")

	getNetworkListSubscription := networklists.GetNetworkListSubscriptionRequest{}

	recipients := d.Get("recipients").([]interface{})
	nru := make([]string, 0, len(recipients))

	for _, h := range recipients {
		nru = append(nru, h.(string))
	}
	getNetworkListSubscription.Recipients = nru

	extractString := strings.Join(getNetworkListSubscription.Recipients, " ")
	recSHA := hash.GetSHAString(extractString)

	uniqueIDs := d.Get("network_list").([]interface{})
	IDs := make([]string, 0, len(uniqueIDs))

	for _, h := range uniqueIDs {
		IDs = append(IDs, h.(string))
	}

	getNetworkListSubscription.UniqueIDs = IDs

	extractStringUID := strings.Join(getNetworkListSubscription.UniqueIDs, " ")
	recSHAUID := hash.GetSHAString(extractStringUID)

	_, err := client.GetNetworkListSubscription(ctx, getNetworkListSubscription)
	if err != nil {
		logger.Errorf("calling 'getNetworkListSubscription': %s", err.Error())
	}

	d.SetId(fmt.Sprintf("%s:%s", recSHA, recSHAUID))

	return nil
}

func resourceNetworkListSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListSubscriptionDelete")

	removeNetworkListSubscription := networklists.RemoveNetworkListSubscriptionRequest{}
	recipients := d.Get("recipients").([]interface{})
	nru := make([]string, 0, len(recipients))

	for _, h := range recipients {
		nru = append(nru, h.(string))
	}
	removeNetworkListSubscription.Recipients = nru

	uniqueIDs := d.Get("network_list").([]interface{})
	IDs := make([]string, 0, len(uniqueIDs))

	for _, h := range uniqueIDs {
		IDs = append(IDs, h.(string))
	}

	removeNetworkListSubscription.UniqueIDs = IDs
	_, errd := client.RemoveNetworkListSubscription(ctx, removeNetworkListSubscription)
	if errd != nil {
		logger.Errorf("calling 'updateNetworkListSubscription': %s", errd.Error())
		return diag.FromErr(errd)
	}
	return nil
}

func resourceNetworkListSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListSubscriptionUpdate")

	updateNetworkListSubscription := networklists.UpdateNetworkListSubscriptionRequest{}

	recipients := d.Get("recipients").([]interface{})
	nru := make([]string, 0, len(recipients))

	for _, h := range recipients {
		nru = append(nru, h.(string))
	}
	updateNetworkListSubscription.Recipients = nru

	extractString := strings.Join(updateNetworkListSubscription.Recipients, " ")
	recSHA := hash.GetSHAString(extractString)

	uniqueIDs := d.Get("network_list").([]interface{})
	IDs := make([]string, 0, len(uniqueIDs))

	for _, h := range uniqueIDs {
		IDs = append(IDs, h.(string))
	}

	updateNetworkListSubscription.UniqueIDs = IDs

	extractStringUID := strings.Join(updateNetworkListSubscription.UniqueIDs, " ")
	recSHAUID := hash.GetSHAString(extractStringUID)

	_, err := client.UpdateNetworkListSubscription(ctx, updateNetworkListSubscription)
	if err != nil {
		logger.Errorf("calling 'updateNetworkListSubscription': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s:%s", recSHA, recSHAUID))

	return resourceNetworkListSubscriptionRead(ctx, d, m)
}
