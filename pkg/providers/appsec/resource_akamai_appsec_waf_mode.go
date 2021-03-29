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
func resourceWAFMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWAFModeUpdate,
		ReadContext:   resourceWAFModeRead,
		UpdateContext: resourceWAFModeUpdate,
		DeleteContext: resourceWAFModeDelete,
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
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					AAG,
					KRS,
				}, false),
			},
			"current_ruleset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eval_ruleset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eval_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eval_expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceWAFModeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFModeRead")

	getWAFMode := appsec.GetWAFModeRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getWAFMode.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getWAFMode.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			getWAFMode.Version = version
		}

		policyid := s[2]
		getWAFMode.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getWAFMode.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getWAFMode.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getWAFMode.PolicyID = policyid
	}
	wafmode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "wafModesDS", wafmode)

	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	if err := d.Set("current_ruleset", wafmode.Current); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("eval_status", wafmode.Eval); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("eval_ruleset", wafmode.Evaluating); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("eval_expiration_date", wafmode.Expires); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", getWAFMode.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getWAFMode.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getWAFMode.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", getWAFMode.ConfigID, getWAFMode.Version, getWAFMode.PolicyID))

	return nil
}

func resourceWAFModeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceWAFModeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFModeUpdate")

	updateWAFMode := appsec.UpdateWAFModeRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateWAFMode.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateWAFMode.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			updateWAFMode.Version = version
		}

		policyid := s[2]
		updateWAFMode.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateWAFMode.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateWAFMode.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateWAFMode.PolicyID = policyid
	}
	mode, err := tools.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateWAFMode.Mode = mode

	//if current mode = one to updare skip this call
	getWAFMode := appsec.GetWAFModeRequest{}
	getWAFMode.ConfigID = updateWAFMode.ConfigID
	getWAFMode.Version = updateWAFMode.Version
	getWAFMode.PolicyID = updateWAFMode.PolicyID

	wafmode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if wafmode.Mode != mode {
		_, erru := client.UpdateWAFMode(ctx, updateWAFMode)
		if erru != nil {
			logger.Errorf("calling 'updateWAFMode': %s", erru.Error())
			return diag.FromErr(erru)
		}
	}

	return resourceWAFModeRead(ctx, d, m)
}

const (
	AAG = "AAG"
	KRS = "KRS"
)
