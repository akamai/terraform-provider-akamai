package appsec

import (
	"context"
	"encoding/json"
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
func resourceMatchTargetSequence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchTargetSequenceUpdate,
		ReadContext:   resourceMatchTargetSequenceRead,
		UpdateContext: resourceMatchTargetSequenceUpdate,
		DeleteContext: resourceMatchTargetSequenceDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_target_sequence": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceMatchTargetSequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceUpdate")

	updateMatchTargetSequence := appsec.UpdateMatchTargetSequenceRequest{}

	jsonpostpayload, ok := d.GetOk("match_target_sequence")
	if ok {

		json.Unmarshal([]byte(jsonpostpayload.(string)), &updateMatchTargetSequence)

		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateMatchTargetSequence.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateMatchTargetSequence.ConfigVersion = version

		d.Set("type", updateMatchTargetSequence.Type)

	}

	_, err := client.UpdateMatchTargetSequence(ctx, updateMatchTargetSequence)
	if err != nil {
		logger.Errorf("calling 'updateMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", updateMatchTargetSequence.ConfigID, updateMatchTargetSequence.ConfigVersion, updateMatchTargetSequence.Type))
	return resourceMatchTargetSequenceRead(ctx, d, m)
}

func resourceMatchTargetSequenceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	d.SetId("")

	return nil
}

func resourceMatchTargetSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceRead")

	getMatchTargetSequence := appsec.GetMatchTargetSequenceRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")
		getMatchTargetSequence.ConfigID, _ = strconv.Atoi(s[0])
		getMatchTargetSequence.ConfigVersion, _ = strconv.Atoi(s[1])

		getMatchTargetSequence.Type = s[2]

	}
	_, err := client.GetMatchTargetSequence(ctx, getMatchTargetSequence)
	if err != nil {
		logger.Errorf("calling 'getMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", getMatchTargetSequence.ConfigID, getMatchTargetSequence.ConfigVersion, getMatchTargetSequence.Type))

	return nil
}
