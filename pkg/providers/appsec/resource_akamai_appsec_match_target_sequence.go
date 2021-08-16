package appsec

import (
	"context"
	"encoding/json"
	"errors"
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
func resourceMatchTargetSequence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchTargetSequenceCreate,
		ReadContext:   resourceMatchTargetSequenceRead,
		UpdateContext: resourceMatchTargetSequenceUpdate,
		DeleteContext: resourceMatchTargetSequenceDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_target_sequence": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceMatchTargetSequenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceCreate")
	logger.Debugf("in resourceMatchTargetSequenceCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "matchTargetSequence", m)
	jsonPayload := d.Get("match_target_sequence")

	createMatchTargetSequence := appsec.UpdateMatchTargetSequenceRequest{}
	if err := json.Unmarshal([]byte(jsonPayload.(string)), &createMatchTargetSequence); err != nil {
		return diag.FromErr(err)
	}
	createMatchTargetSequence.ConfigID = configid
	createMatchTargetSequence.ConfigVersion = version

	_, err = client.UpdateMatchTargetSequence(ctx, createMatchTargetSequence)
	if err != nil {
		logger.Errorf("calling 'updateMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createMatchTargetSequence.ConfigID, createMatchTargetSequence.Type))
	return resourceMatchTargetSequenceRead(ctx, d, m)
}

func resourceMatchTargetSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceRead")
	logger.Debugf("in resourceMatchTargetSequenceRead")

	idParts, err := splitID(d.Id(), 2, "configid:matchtargettype")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	matchTargetType := idParts[1]

	getMatchTargetSequence := appsec.GetMatchTargetSequenceRequest{}
	getMatchTargetSequence.ConfigID = configid
	getMatchTargetSequence.ConfigVersion = version
	getMatchTargetSequence.Type = matchTargetType

	matchTargetSequence, err := client.GetMatchTargetSequence(ctx, getMatchTargetSequence)
	if err != nil {
		logger.Errorf("calling 'getMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchTargetSequence)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("match_target_sequence", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceMatchTargetSequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceUpdate")
	logger.Debugf("in resourceMatchTargetSequenceUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:matchtargettype")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "matchTargetSequence", m)
	matchTargetType := idParts[1]

	jsonPayload := d.Get("match_target_sequence")

	updateMatchTargetSequence := appsec.UpdateMatchTargetSequenceRequest{}
	if err := json.Unmarshal([]byte(jsonPayload.(string)), &updateMatchTargetSequence); err != nil {
		return diag.FromErr(err)
	}
	updateMatchTargetSequence.ConfigID = configid
	updateMatchTargetSequence.ConfigVersion = version

	if matchTargetType != updateMatchTargetSequence.Type {
		err = fmt.Errorf("match target type %s cannot be changed to %s", matchTargetType, updateMatchTargetSequence.Type)
		return diag.FromErr(err)
	}

	_, err = client.UpdateMatchTargetSequence(ctx, updateMatchTargetSequence)
	if err != nil {
		logger.Errorf("calling 'updateMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceMatchTargetSequenceRead(ctx, d, m)
}

func resourceMatchTargetSequenceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(context.TODO(), d, m)
}
