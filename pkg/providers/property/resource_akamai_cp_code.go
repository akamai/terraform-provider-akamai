package property

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"
)

// PAPI CP Code
//
// https://techdocs.akamai.com/property-mgr/reference/post-cpcodes
func resourceCPCode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPCodeCreate,
		ReadContext:   resourceCPCodeRead,
		UpdateContext: resourceCPCodeUpdate,
		// NB: CP Codes cannot be deleted https://techdocs.akamai.com/property-mgr/reference/post-cpcodes
		DeleteContext: schema.NoopContext,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCPCodeImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tf.IsNotBlank,
			},
			"contract": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("contract"),
				StateFunc:  addPrefixToState("ctr_"),
			},
			"contract_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract_id", "contract"},
				StateFunc:    addPrefixToState("ctr_"),
			},
			"group": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("group"),
				StateFunc:  addPrefixToState("grp_"),
			},
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group_id", "group"},
				StateFunc:    addPrefixToState("grp_"),
			},
			"product": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("product"),
				StateFunc:  addPrefixToState("prd_"),
			},
			"product_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"product"},
				StateFunc:    addPrefixToState("prd_"),
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Update: &cpCodeResourceUpdateTimeout,
		},
	}
}

var (
	updatePollMinimum           = time.Minute
	updatePollInterval          = updatePollMinimum
	cpCodeResourceUpdateTimeout = time.Minute * 30
)

const cpCodePrefix = "cpc_"

func resourceCPCodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("PAPI", "resourceCPCodeCreate")
	logger.Debugf("Creating CP Code")

	var name string
	if got, ok := d.GetOk("name"); ok {
		name = got.(string)
	}

	// Schema guarantees product_id/product are strings and one or the other is set
	var productID string
	if got, ok := d.GetOk("product_id"); ok {
		productID = got.(string)
	} else {
		productID = d.Get("product").(string)
	}
	productID = tools.AddPrefix(productID, "prd_")

	contractID, groupID := getContractIDAndGroupID(d)

	var cpCodeID string
	// Because CPCodes can't be deleted, we re-use an existing CPCode if it's there
	cpCode, err := findCPCode(ctx, client, name, contractID, groupID)
	if err != nil && !errors.Is(err, ErrCpCodeNotFound) {
		return diag.Errorf("%s: %s", ErrLookingUpCPCode, err)
	}

	if errors.Is(err, ErrCpCodeNotFound) {
		cpCodeID, err = createCPCode(ctx, client, name, productID, contractID, groupID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		cpCodeID = cpCode.ID
	}

	d.SetId(strings.TrimPrefix(cpCodeID, cpCodePrefix))
	return resourceCPCodeRead(ctx, d, m)
}

func resourceCPCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceCPCodeRead")
	client := inst.Client(meta)
	logger.Debugf("Read CP Code")

	contractID, groupID := getContractIDAndGroupID(d)

	if err := d.Set("group_id", groupID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("group", groupID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("contract_id", contractID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("contract", contractID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	cpCodeResp, err := client.GetCPCode(ctx, papi.GetCPCodeRequest{
		CPCodeID:   d.Id(),
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	cpCode := cpCodeResp.CPCode

	if err := d.Set("name", cpCode.Name); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	// we use the first value returned.  Most cpcodes have but a single product and we need to pick one for comparison.
	if len(cpCode.ProductIDs) == 0 {
		return diag.Errorf("Couldn't find product id on the CP Code")
	}
	if err := d.Set("product", cpCode.ProductIDs[0]); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("product_id", cpCode.ProductIDs[0]); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	logger.Debugf("Read CP Code: %+v", cpCode)
	return nil
}

func resourceCPCodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceCPCodeUpdate")
	client := inst.Client(meta)
	logger.Debugf("Update CP Code")

	if diags := checkImmutableChanged(d); diags != nil {
		d.Partial(true)
		return diags
	}

	contractID, groupID := getContractIDAndGroupID(d)

	// trimCPCodeID is needed here for backwards compatibility
	cpCodeID, err := strconv.Atoi(strings.TrimPrefix(d.Id(), cpCodePrefix))
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	cpCode, err := client.GetCPCodeDetail(ctx, cpCodeID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCPCode(ctx, papi.UpdateCPCodeRequest{
		ID:               cpCode.ID,
		Name:             name,
		Purgeable:        &cpCode.Purgeable,
		OverrideTimeZone: &cpCode.OverrideTimeZone,
		Contracts:        cpCode.Contracts,
		Products:         cpCode.Products,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// Because we use CPRG API for update, we need to ensure that changes are also present when fetching cpCode with PAPI
	if err := waitForCPCodeNameUpdate(ctx, client, contractID, groupID, d.Id(), name); err != nil {
		if errors.Is(err, ErrCPCodeUpdateTimeout) {
			return append(tools.DiagWarningf("%s", err), tools.DiagWarningf("Resource has been updated, but the change is still ongoing on the server")...)
		}
		return diag.FromErr(err)
	}

	return resourceCPCodeRead(ctx, d, m)
}

func resourceCPCodeImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceCPCodeImport")
	client := inst.Client(meta)
	logger.Debugf("Import CP Code")

	parts := strings.Split(d.Id(), ",")

	if len(parts) < 3 {
		return nil, fmt.Errorf("comma-separated list of CP code ID, contract ID and group ID has to be supplied in import: %s", d.Id())
	}
	if parts[0] == "" {
		return nil, errors.New("CP Code is a mandatory parameter")
	}
	cpCodeID := parts[0]
	contractID := tools.AddPrefix(parts[1], "ctr_")
	groupID := tools.AddPrefix(parts[2], "grp_")

	cpCodeResp, err := client.GetCPCode(ctx, papi.GetCPCodeRequest{
		CPCodeID:   cpCodeID,
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return nil, err
	}

	cpCode := cpCodeResp.CPCode

	if err := d.Set("name", cpCode.Name); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("contract_id", contractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("group_id", groupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if len(cpCode.ProductIDs) == 0 {
		return nil, fmt.Errorf("could not find product id on the CP Code")
	}
	if err := d.Set("product_id", cpCode.ProductIDs[0]); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("product", cpCode.ProductIDs[0]); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strings.TrimPrefix(cpCode.ID, cpCodePrefix))
	logger.Debugf("Import CP Code: %+v", cpCode)
	return []*schema.ResourceData{d}, nil
}

// createCPCode attempts to create a CP Code and returns the CP Code ID
func createCPCode(ctx context.Context, client papi.PAPI, name, productID, contractID, groupID string) (string, error) {
	r, err := client.CreateCPCode(ctx, papi.CreateCPCodeRequest{
		ContractID: contractID,
		GroupID:    groupID,
		CPCode: papi.CreateCPCode{
			ProductID:  productID,
			CPCodeName: name,
		},
	})
	if err != nil {
		return "", err
	}

	return r.CPCodeID, nil
}

func checkImmutableChanged(d *schema.ResourceData) diag.Diagnostics {
	immutables := []string{
		"contract",
		"contract_id",
		"group",
		"group_id",
		"product",
		"product_id",
	}

	var diags diag.Diagnostics
	for _, immutable := range immutables {
		if d.HasChange(immutable) {
			diags = append(diags, diag.Errorf("cp code attribute '%s' cannot be changed after creation (immutable)", immutable)...)
		}
	}
	return diags
}

func getContractIDAndGroupID(d *schema.ResourceData) (contractID, groupID string) {
	// Schema guarantees contract_id/contract are strings and one or the other is set
	if got, ok := d.GetOk("contract_id"); ok {
		contractID = got.(string)
	} else {
		contractID = d.Get("contract").(string)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")

	// Schema guarantees group_id/group are strings and one or the other is set
	if got, ok := d.GetOk("group_id"); ok {
		groupID = got.(string)
	} else {
		groupID = d.Get("group").(string)
	}
	groupID = tools.AddPrefix(groupID, "grp_")

	return
}

func waitForCPCodeNameUpdate(ctx context.Context, client papi.PAPI, contractID, groupID, CPCodeID, updatedName string) error {
	req := papi.GetCPCodeRequest{CPCodeID: CPCodeID, ContractID: contractID, GroupID: groupID}
	CPCodeResp, err := client.GetCPCode(ctx, req)
	if err != nil {
		return err
	}

	for CPCodeResp.CPCode.Name != updatedName {
		select {
		case <-time.After(tf.MaxDuration(updatePollInterval, updatePollMinimum)):
			CPCodeResp, err = client.GetCPCode(ctx, req)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return ErrCPCodeUpdateTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return fmt.Errorf("operation cancelled while waiting for CPCode update")
			}
			return fmt.Errorf("cp code update context terminated: %w", ctx.Err())
		}
	}

	return nil
}
