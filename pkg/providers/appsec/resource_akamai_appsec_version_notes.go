package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceVersionNotes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVersionNotesCreate,
		ReadContext:   resourceVersionNotesRead,
		UpdateContext: resourceVersionNotesUpdate,
		DeleteContext: resourceVersionNotesDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version_notes": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceVersionNotesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceVersionNotesCreate")
	logger.Debugf("in resourceVersionNotesCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "editVersionNotes", m)
	notes, err := tools.GetStringValue("version_notes", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createVersionNotes := appsec.UpdateVersionNotesRequest{
		ConfigID: configid,
		Version:  version,
		Notes:    notes,
	}

	_, erru := client.UpdateVersionNotes(ctx, createVersionNotes)
	if erru != nil {
		logger.Errorf("calling 'createVersionNotes': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d", createVersionNotes.ConfigID))

	return resourceVersionNotesRead(ctx, d, m)
}

func resourceVersionNotesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceVersionNotesRead")
	logger.Debugf("resourceVersionNotesRead")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)

	getVersionNotes := appsec.GetVersionNotesRequest{
		ConfigID: configid,
		Version:  version,
	}

	versionnotes, err := client.GetVersionNotes(ctx, getVersionNotes)
	if err != nil {
		logger.Errorf("calling 'getVersionNotes': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getVersionNotes.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("version_notes", versionnotes.Notes); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "versionNotesDS", versionnotes)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	return nil
}

func resourceVersionNotesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceVersionNotesUpdate")
	logger.Debugf("resourceVersionNotesUpdate")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "editVersionNotes", m)
	notes, err := tools.GetStringValue("version_notes", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateVersionNotes := appsec.UpdateVersionNotesRequest{
		ConfigID: configid,
		Version:  version,
		Notes:    notes,
	}

	_, erru := client.UpdateVersionNotes(ctx, updateVersionNotes)
	if erru != nil {
		logger.Errorf("calling 'updateVersionNotes': %s", erru.Error())
		return diag.FromErr(erru)
	}
	return resourceVersionNotesRead(ctx, d, m)
}

func resourceVersionNotesDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(context.TODO(), d, m)
}
