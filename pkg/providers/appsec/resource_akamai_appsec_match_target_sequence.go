package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"match_target_sequence": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
				Description:      "JSON-formatted definition of the processing sequence for all defined match targets ",
			},
		},
	}
}

func resourceMatchTargetSequenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceCreate")
	logger.Debugf("in resourceMatchTargetSequenceCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "matchTargetSequence", m)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonPayload := d.Get("match_target_sequence")

	createMatchTargetSequence := appsec.UpdateMatchTargetSequenceRequest{}
	if err := json.Unmarshal([]byte(jsonPayload.(string)), &createMatchTargetSequence); err != nil {
		return diag.FromErr(err)
	}
	createMatchTargetSequence.ConfigID = configID
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
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceRead")
	logger.Debugf("in resourceMatchTargetSequenceRead")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetType")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	matchTargetType := iDParts[1]

	getMatchTargetSequence := appsec.GetMatchTargetSequenceRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		Type:          matchTargetType,
	}

	matchTargetSequence, err := client.GetMatchTargetSequence(ctx, getMatchTargetSequence)
	if err != nil {
		logger.Errorf("calling 'getMatchTargetSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchTargetSequence)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("match_target_sequence", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceMatchTargetSequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetSequenceUpdate")
	logger.Debugf("in resourceMatchTargetSequenceUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetType")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "matchTargetSequence", m)
	if err != nil {
		return diag.FromErr(err)
	}
	matchTargetType := iDParts[1]

	jsonPayload := d.Get("match_target_sequence")

	updateMatchTargetSequence := appsec.UpdateMatchTargetSequenceRequest{}
	if err := json.Unmarshal([]byte(jsonPayload.(string)), &updateMatchTargetSequence); err != nil {
		return diag.FromErr(err)
	}
	updateMatchTargetSequence.ConfigID = configID
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

func resourceMatchTargetSequenceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(ctx, d, m)
}
