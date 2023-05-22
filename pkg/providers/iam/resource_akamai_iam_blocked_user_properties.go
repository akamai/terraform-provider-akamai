package iam

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIAMBlockedUserProperties() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a user in your account",
		CreateContext: resourceIAMBlockedUserPropertiesCreate,
		ReadContext:   resourceIAMBlockedUserPropertiesRead,
		UpdateContext: resourceIAMBlockedUserPropertiesUpdate,
		DeleteContext: resourceIAMBlockedUserPropertiesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIAMBlockedUserPropertiesImport,
		},
		Schema: map[string]*schema.Schema{
			"identity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique identifier for a user's profile, which corresponds to a user's actual profile or client ID",
			},
			"group_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "A unique identifier for a group",
			},
			"blocked_properties": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of properties to block for a user",
			},
		},
	}
}

func resourceIAMBlockedUserPropertiesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMBlockedUserPropertiesCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Creating blocked user properties")

	identityID, err := tf.GetStringValue("identity_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tf.GetIntValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	listRequest := iam.ListBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    int64(groupID),
	}

	existingBlockedProperties, err := client.ListBlockedProperties(ctx, listRequest)
	if err != nil {
		logger.WithError(err).Errorf("failed to fetch blocked user properties")
		return diag.Errorf("failed to fetch blocked user properties: %s", err)
	}
	if len(existingBlockedProperties) > 0 {
		logger.Errorf("there are already blocked properties on server, please import resource first")
		return diag.Errorf("there are already blocked properties on server, please import resource first")
	}

	blockedProperties, err := tf.GetInterfaceArrayValue("blocked_properties", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := iam.UpdateBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    int64(groupID),
		Properties: tf.ConvertListOfIntToInt64(blockedProperties),
	}

	_, err = client.UpdateBlockedProperties(ctx, request)
	if err != nil {
		logger.WithError(err).Errorf("failed to create blocked user properties")
		return diag.Errorf("failed to create blocked user properties: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%d", request.IdentityID, request.GroupID))
	return resourceIAMBlockedUserPropertiesRead(ctx, d, m)
}

func resourceIAMBlockedUserPropertiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMBlockedUserPropertiesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading blocked user properties")

	idParts, err := splitID(d.Id(), 2, "identityID:groupID")
	if err != nil {
		return diag.FromErr(err)
	}

	identityID := idParts[0]
	groupID, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	request := iam.ListBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    groupID,
	}

	blockedProperties, err := client.ListBlockedProperties(ctx, request)
	if err != nil {
		logger.WithError(err).Errorf("failed to fetch blocked user properties")
		return diag.Errorf("failed to fetch blocked user properties: %s", err)
	}

	if err = d.Set("identity_id", identityID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err = d.Set("group_id", groupID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err = d.Set("blocked_properties", blockedProperties); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceIAMBlockedUserPropertiesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMBlockedUserPropertiesUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Updating blocked user properties")

	identityID, err := tf.GetStringValue("identity_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tf.GetIntValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	blockedProperties, err := tf.GetInterfaceArrayValue("blocked_properties", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := iam.UpdateBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    int64(groupID),
		Properties: tf.ConvertListOfIntToInt64(blockedProperties),
	}

	_, err = client.UpdateBlockedProperties(ctx, request)
	if err != nil {
		logger.WithError(err).Errorf("failed to update blocked user properties")
		return diag.Errorf("failed to update blocked user properties: %s", err)
	}

	return resourceIAMBlockedUserPropertiesRead(ctx, d, m)
}

func resourceIAMBlockedUserPropertiesDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMBlockedUserPropertiesDelete")

	logger.Debug("Deleting blocked user properties")

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "IAM API does not support deletion of blocked properties - resource will only be removed from state.",
		},
	}
}

func resourceIAMBlockedUserPropertiesImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMBlockedUserPropertiesImport")
	logger.Debug("Importing blocked user properties")

	if d.Id() == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	return schema.ImportStatePassthroughContext(ctx, d, m)
}

func splitID(id string, expectedNum int, example string) ([]string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != expectedNum {
		return nil, fmt.Errorf("id '%s' is incorrectly formatted: should be of form '%s'", id, example)
	}
	return parts, nil
}
