package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
				ConflictsWith: []string{"sequence_map"},
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressJsonProvided,
			},
			"sequence_map": {
				Type:             schema.TypeMap,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
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

		matchtargetseqtype, err := tools.GetStringValue("type", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateMatchTargetSequence.Type = matchtargetseqtype

		d.Set("type", updateMatchTargetSequence.Type)

		sequenceMap, ok := d.Get("sequence_map").(map[string]interface{})
		if !ok {
			logger.Warnf("get map  'updateMatchTargetSequence': %s", err.Error())
			return diag.FromErr(err)
		}
		logger.Warnf("calling 'getMatchTargetSequence SEQ MAP': %s", sequenceMap)

		for target, sequence := range sequenceMap {
			logger.Warnf("calling 'getMatchTargetSequence SEQ MAP LOOP': %s", target)
			targetsequence.TargetID, _ = strconv.Atoi(target)

			targetsequence.Sequence, _ = strconv.Atoi(sequence.(string))

			updateMatchTargetSequence.TargetSequence = append(updateMatchTargetSequence.TargetSequence, targetsequence)

		}
		logger.Warnf("calling 'getMatchTargetSequence SEQ MAP LOOP EXIT ': %s", updateMatchTargetSequence)
	}

	updatematchtargetsequence, err := client.UpdateMatchTargetSequence(ctx, updateMatchTargetSequence)
	if err != nil {
		logger.Warnf("calling 'updateMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	targetsequence = v2.TargetSequence{}
	sequencemap := []v2.TargetSequence{}

	for _, targets := range updatematchtargetsequence.TargetSequence {
		logger.Warnf("calling 'getMatchTargetSequence SEQ MAP LOOP': %s", targets)
		targetsequence.TargetID = targets.TargetID
		targetsequence.Sequence = targets.Sequence
		sequencemap = append(sequencemap, targetsequence)

	}

	d.Set("sequence_map", sequenceToMap(sequencemap))

	d.SetId(fmt.Sprintf("%d:%d", updateMatchTargetSequence.ConfigID, updateMatchTargetSequence.ConfigVersion))
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
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")
		getMatchTargetSequences.ConfigID, _ = strconv.Atoi(s[0])
		getMatchTargetSequences.ConfigVersion, _ = strconv.Atoi(s[1])

		matchtargetseqtype, err := tools.GetStringValue("type", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getMatchTargetSequences.Type = matchtargetseqtype
		d.Set("type", getMatchTargetSequences.Type)

	} else {
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

		matchtargetseqtype, err := tools.GetStringValue("type", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getMatchTargetSequences.Type = matchtargetseqtype
		d.Set("type", getMatchTargetSequences.Type)
	}

	matchtargetsequences, err := client.GetMatchTargetSequences(ctx, getMatchTargetSequences)
	if err != nil {
		logger.Warnf("calling 'getMatchTargetSequence': %s", err.Error())
	}

	logger.Warnf("calling 'getMatchTargetSequence': %s", matchtargetsequences.MatchTargets.WebsiteTargets)
	targetsequence := v2.TargetSequence{}
	sequencemap := []v2.TargetSequence{}

	if getMatchTargetSequences.Type == "website" {
		for _, targets := range matchtargetsequences.MatchTargets.WebsiteTargets {
			logger.Warnf("calling 'getMatchTargetSequence SEQ MAP LOOP': %s", targets)
			targetsequence.TargetID = targets.TargetID
			targetsequence.Sequence = targets.Sequence
			sequencemap = append(sequencemap, targetsequence)

		}
	}

	if getMatchTargetSequences.Type == "api" {
		for _, targets := range matchtargetsequences.MatchTargets.APITargets {
			logger.Warnf("calling 'getMatchTargetSequence SEQ MAP LOOP': %s", targets)
			targetsequence.TargetID = targets.TargetID
			targetsequence.Sequence = targets.Sequence
			sequencemap = append(sequencemap, targetsequence)

		}
	}

	logger.Warnf("calling 'getMatchTargetSequence SEQ MAP LOOP EXIT ': %s", sequencemap)
	d.Set("sequence_map", sequenceToMap(sequencemap))

	d.Set("type", getMatchTargetSequences.Type)

	d.SetId(fmt.Sprintf("%d:%d", getMatchTargetSequences.ConfigID, getMatchTargetSequences.ConfigVersion))

	return nil
}

func sequenceToMap(sequenceMap []v2.TargetSequence) map[string]interface{} {
	var sequencemap = make(map[string]interface{})
	if len(sequenceMap) > 0 {
		for _, seqs := range sequenceMap {
			target := strconv.Itoa(seqs.TargetID)
			seq := strconv.Itoa(seqs.Sequence)
			sequencemap[target] = seq

		}
	}
	return sequencemap
}
