package iam

import (
	"context"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
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
				Type:             schema.TypeList,
				Elem:             &schema.Schema{Type: schema.TypeInt},
				Required:         true,
				DiffSuppressFunc: suppressDiffInGrantedRoles,
				Description:      "The list of existing unique identifiers for the granted roles",
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
	meta := meta.Must(m)
	logger := meta.Log("IAM", "resourceIAMRoleCreate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating Role")

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tf.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}
	grantedRoles, err := tf.GetListValue("granted_roles", d)
	if err != nil {
		return diag.FromErr(err)
	}

	grantedRolesIDs := getGrantedRolesIDs(grantedRoles)

	role, err := client.CreateRole(ctx, iam.CreateRoleRequest{
		Name:         name,
		Description:  description,
		GrantedRoles: grantedRolesIDs,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(int(role.RoleID)))
	return resourceIAMRoleRead(ctx, d, m)
}

func resourceIAMRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "resourceIAMRoleRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Reading Role")

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("%s: %s", tf.ErrInvalidType, err.Error())
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
	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIAMRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
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
		return diag.Errorf("%s: %s", tf.ErrInvalidType, err.Error())
	}

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tf.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}
	grantedRoles, err := tf.GetListValue("granted_roles", d)
	if err != nil {
		return diag.FromErr(err)
	}

	grantedRolesIDs := getGrantedRolesIDs(grantedRoles)

	_, err = client.UpdateRole(ctx, iam.UpdateRoleRequest{
		ID: int64(roleIDReq),
		RoleRequest: iam.RoleRequest{
			Name:         name,
			Description:  description,
			GrantedRoles: grantedRolesIDs,
		},
	})

	diags := diag.Diagnostics{}
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return append(resourceIAMRoleRead(ctx, d, m), diags...)
}

func resourceIAMRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
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
		return diag.Errorf("%s: %s", tf.ErrInvalidType, err.Error())
	}

	err = client.DeleteRole(ctx, iam.DeleteRoleRequest{
		ID: int64(roleIDReq),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getGrantedRolesIDs(grantedRoles []interface{}) []iam.GrantedRoleID {
	grantedRolesIDs := make([]iam.GrantedRoleID, 0, len(grantedRoles))
	for _, role := range grantedRoles {
		grantedRolesIDs = append(grantedRolesIDs, iam.GrantedRoleID{ID: int64(role.(int))})
	}
	return grantedRolesIDs
}

func suppressDiffInGrantedRoles(_, o, n string, d *schema.ResourceData) bool {
	key := "granted_roles"

	oldValue, newValue := d.GetChange(key)
	oldGrantedRoles := oldValue.([]interface{})
	newGrantedRoles := newValue.([]interface{})
	if len(oldGrantedRoles) != len(newGrantedRoles) {
		return o == n
	}

	oldGrantedRolesIDs := make([]int, 0, len(oldGrantedRoles))
	for _, v := range oldGrantedRoles {
		oldGrantedRolesIDs = append(oldGrantedRolesIDs, v.(int))
	}

	newGrantedRolesIDs := make([]int, 0, len(newGrantedRoles))
	for _, v := range newGrantedRoles {
		newGrantedRolesIDs = append(newGrantedRolesIDs, v.(int))
	}

	sort.Ints(oldGrantedRolesIDs)
	sort.Ints(newGrantedRolesIDs)

	for i := range oldGrantedRoles {
		if oldGrantedRolesIDs[i] != newGrantedRolesIDs[i] {
			return false
		}
	}
	return true
}
