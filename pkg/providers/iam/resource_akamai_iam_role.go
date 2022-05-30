package iam

import (
	"context"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIAMRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMRoleCreate,
		ReadContext:   resourceIAMRoleRead,
		UpdateContext: resourceIAMRoleUpdate,
		DeleteContext: resourceIAMRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name you supply for a role",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description for a role",
			},
			"granted_roles": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
				Description: "The list of existing unique identifiers for the granted roles",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The role type which indicates whether it's a standard role provided by Akamai or a custom role for the account",
			},
		},
	}
}

func resourceIAMRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMRoleCreate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating Role")

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}
	grantedRoles, err := tools.GetSetValue("granted_roles", d)

	grantedRolesIDs := getSortedGrantedRolesIDs(tools.SetToIntSlice(grantedRoles))

	role, err := client.CreateRole(ctx, iam.CreateRoleRequest{
		Name:         name,
		Description:  description,
		GrantedRoles: grantedRolesIDs,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(role.RoleID, 10))
	return resourceIAMRoleRead(ctx, d, m)
}

func resourceIAMRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMRoleRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Reading Role")

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("%s: %s", tools.ErrInvalidType, err.Error())
	}
	role, err := client.GetRole(ctx, iam.GetRoleRequest{
		ID:           int64(roleID),
		GrantedRoles: true,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	grantedRolesIDs := make([]int, 0, len(role.GrantedRoles))
	for _, r := range role.GrantedRoles {
		grantedRolesIDs = append(grantedRolesIDs, int(r.RoleID))
	}

	attrs := make(map[string]interface{})
	attrs["name"] = role.RoleName
	attrs["description"] = role.RoleDescription
	attrs["type"] = role.RoleType
	attrs["granted_roles"] = grantedRolesIDs
	if err = tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIAMRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMRoleUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Updating Role")
	roleID := d.Id()
	roleIDReq, err := strconv.Atoi(roleID)
	if err != nil {
		return diag.Errorf("%s: %s", tools.ErrInvalidType, err.Error())
	}

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}
	grantedRoles, err := tools.GetSetValue("granted_roles", d)
	if err != nil {
		return diag.FromErr(err)
	}

	grantedRolesIDs := getSortedGrantedRolesIDs(tools.SetToIntSlice(grantedRoles))

	_, err = client.UpdateRole(ctx, iam.UpdateRoleRequest{
		ID: int64(roleIDReq),
		RoleRequest: iam.RoleRequest{
			Name:         name,
			Description:  description,
			GrantedRoles: grantedRolesIDs,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIAMRoleRead(ctx, d, m)
}

func resourceIAMRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMRoleDelete")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Deleting Role")
	roleID := d.Id()
	roleIDReq, err := strconv.Atoi(roleID)
	if err != nil {
		return diag.Errorf("%s: %s", tools.ErrInvalidType, err.Error())
	}

	err = client.DeleteRole(ctx, iam.DeleteRoleRequest{
		ID: int64(roleIDReq),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getSortedGrantedRolesIDs(grantedRoles []int) []iam.GrantedRoleID {
	sort.Ints(grantedRoles)
	grantedRolesIDs := make([]iam.GrantedRoleID, 0, len(grantedRoles))
	for _, role := range grantedRoles {
		grantedRolesIDs = append(grantedRolesIDs, iam.GrantedRoleID{ID: int64(role)})
	}
	return grantedRolesIDs
}
