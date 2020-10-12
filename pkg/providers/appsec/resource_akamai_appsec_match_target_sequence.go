package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"json": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"type", "sequence_map"},
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"sequence_map": {
				Type:             schema.TypeMap,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeInt},
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
		},
	}
}

func resourceMatchTargetSequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceUpdate")

	updateMatchTargetSequence := v2.UpdateMatchTargetSequenceRequest{}
	targetsequence := v2.TargetSequence{}

	jsonpostpayload, ok := d.GetOk("json")
	if ok {

		json.Unmarshal([]byte(jsonpostpayload.(string)), &updateMatchTargetSequence)
		//updateMatchTargetSequence.TargetSequence.TargetID, _ = strconv.Atoi(d.Id())
	} else {
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

		//var ts []v2.TargetSequence
		sequenceMap, ok := d.Get("sequence_map").(map[int]interface{})
		if !ok {
			return diag.FromErr(err)
		}

		//sequenceMaps := make(map[int]int, len(sequenceMap))
		for target, sequence := range sequenceMap {
			targetsequence.TargetID = target
			targetsequence.Sequence = sequence.(int)
			//ts = append(ts, targetsequence)
			updateMatchTargetSequence.TargetSequence = append(updateMatchTargetSequence.TargetSequence, targetsequence)

		}

		//updateMatchTargetSequence.TargetID, _ = strconv.Atoi(d.Id())
		updateMatchTargetSequence.Type = d.Get("type").(string)

	}

	_, err := client.UpdateMatchTargetSequence(ctx, updateMatchTargetSequence)
	if err != nil {
		logger.Warnf("calling 'updateMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

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

	getMatchTargetSequences := v2.GetMatchTargetSequencesRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getMatchTargetSequences.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getMatchTargetSequences.ConfigVersion = version

	getMatchTargetSequences.ConfigID, _ = strconv.Atoi(d.Id())

	matchtargetsequences, err := client.GetMatchTargetSequences(ctx, getMatchTargetSequences)
	if err != nil {
		logger.Warnf("calling 'getMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	d.Set("type", matchtargetsequences.MatchTargets.APITargets[0].Type)

	d.Set("target_id", matchtargetsequences.MatchTargets.WebsiteTargets[0].TargetID)
	d.SetId(strconv.Itoa(matchtargetsequences.MatchTargets.WebsiteTargets[0].TargetID))

	return nil
}
