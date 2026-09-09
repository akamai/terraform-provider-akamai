package property

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/providers/property/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cpCodePrefix = "cpc_"

type cpCodeResourceConfig struct {
	updatePollMinimum           time.Duration
	updatePollInterval          time.Duration
	cpCodeResourceUpdateTimeout time.Duration
}

func defaultCPCodeResourceConfig() cpCodeResourceConfig {
	return cpCodeResourceConfig{
		updatePollMinimum:           time.Minute,
		updatePollInterval:          time.Minute,
		cpCodeResourceUpdateTimeout: time.Minute * 30,
	}
}

// PAPI CP Code
//
// https://techdocs.akamai.com/property-mgr/reference/post-cpcodes
func resourceCPCode(config cpCodeResourceConfig) *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPCodeCreate,
		ReadContext:   resourceCPCodeRead,
		UpdateContext: resourceCPCodeUpdate(config),
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
			"contract_id": {
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: addPrefixToState("ctr_"),
			},
			"group_id": {
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: addPrefixToState("grp_"),
			},
			"product_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: addPrefixToState("prd_"),
			},
			"timeouts": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enables to set timeout for processing",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"update": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: timeouts.ValidateDurationFormat,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Update: &config.cpCodeResourceUpdateTimeout,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourceCPCodeV0(config).CoreConfigSchema().ImpliedType(),
			Upgrade: timeouts.MigrateToExplicit(),
		}},
	}
}

func resourceCPCodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceCPCodeCreate")
	logger.Debugf("Creating CP Code")

	var name string
	if got, ok := d.GetOk("name"); ok {
		name = got.(string)
	}

	// Schema no longer guarantees that product_id is set, this field is required only for creation
	productID, err := tf.GetStringValue("product_id", d)
	if err != nil {
		return diag.Errorf("`product_id` must be specified for creation")
	}
	productID = str.AddPrefix(productID, "prd_")

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = str.AddPrefix(contractID, "ctr_")

	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID = str.AddPrefix(groupID, "grp_")

	client := meta.Client().GetPAPI()

	products, err := client.GetProducts(ctx, papi.GetProductsRequest{
		ContractID: contractID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var productIDExists bool
	for _, v := range products.Products.Items {
		if productID == str.AddPrefix(v.ProductID, "prd_") {
			productIDExists = true
			break
		}
	}

	if !productIDExists {
		return diag.Errorf("`product_id` `%s` does not exist under contract `%s`, you need to provide a valid `product_id`", productID, contractID)
	}

	var cpCodeID string
	// Because CPCodes can't be deleted, we re-use an existing CPCode if it's there
	cpCode, err := tools.FindCPCodeByName(ctx, client, name, contractID, groupID)
	if err != nil && !errors.Is(err, tools.ErrCPCodeNotFound) {
		return diag.Errorf("%s: %s", ErrLookingUpCPCodeByName, err)
	}

	var diags diag.Diagnostics
	if errors.Is(err, tools.ErrCPCodeNotFound) {
		cpCodeID, err = createCPCode(ctx, client, name, productID, contractID, groupID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		// we use the first value returned. Most cpcodes have but a single product and we need to pick one for comparison.
		if len(cpCode.ProductIDs) == 0 {
			return diag.Errorf(errCPCodeNoProductID, name)
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "CP Code already exists",
			Detail:   fmt.Sprintf("CP Code %q has associated product IDs; skipping creation of a new CP Code.", name),
		})
		cpCodeID = cpCode.ID
	}

	d.SetId(strings.TrimPrefix(cpCodeID, cpCodePrefix))

	readDiags := resourceCPCodeRead(ctx, d, m)
	diags = append(diags, readDiags...)

	return diags
}

func resourceCPCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceCPCodeRead")
	logger.Debugf("Read CP Code")

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = str.AddPrefix(contractID, "ctr_")

	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID = str.AddPrefix(groupID, "grp_")

	if err := d.Set("group_id", groupID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("contract_id", contractID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	cpCodeResp, err := meta.Client().GetPAPI().GetCPCode(ctx, papi.GetCPCodeRequest{
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
		return diag.Errorf(errCPCodeNoProductID, cpCode.Name)
	}

	if err := d.Set("product_id", cpCode.ProductIDs[0]); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	logger.Debugf("Read CP Code: %+v", cpCode)
	return nil
}

func resourceCPCodeUpdate(config cpCodeResourceConfig) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		meta := meta.Must(m)
		logger := meta.Log("PAPI", "resourceCPCodeUpdate")
		logger.Debugf("Update CP Code")

		if !d.HasChangeExcept("timeouts") {
			logger.Debug("Only timeouts were updated, skipping")
			return nil
		}

		if diags := checkImmutableChanged(d); diags != nil {
			d.Partial(true)
			return diags
		}

		contractID, err := tf.GetStringValue("contract_id", d)
		if err != nil {
			return diag.FromErr(err)
		}
		contractID = str.AddPrefix(contractID, "ctr_")
		groupID, err := tf.GetStringValue("group_id", d)
		if err != nil {
			return diag.FromErr(err)
		}
		groupID = str.AddPrefix(groupID, "grp_")

		// trimCPCodeID is needed here for backwards compatibility
		cpCodeID, err := strconv.Atoi(strings.TrimPrefix(d.Id(), cpCodePrefix))
		if err != nil {
			return diag.FromErr(err)
		}

		name, err := tf.GetStringValue("name", d)
		if err != nil {
			return diag.FromErr(err)
		}

		cpCode, err := meta.Client().GetReportingGroups().GetCPCode(ctx, reportinggroups.GetCPCodeRequest{CPCodeID: int64(cpCodeID)})
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = meta.Client().GetReportingGroups().UpdateCPCode(ctx, reportinggroups.UpdateCPCodeRequest{
			CPCodeID:         cpCode.CPCodeID,
			CPCodeName:       name,
			Purgeable:        &cpCode.Purgeable,
			OverrideTimeZone: ptr.To(cpCode.OverrideTimeZone),
			Contracts:        cpCode.Contracts,
			Products:         cpCode.Products,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		// Because we use CPRG API for update, we need to ensure that changes are also present when fetching cpCode with PAPI
		if err := waitForCPCodeNameUpdate(ctx, waitForCPCodeNameUpdateOpts{
			client:      meta.Client().GetPAPI(),
			contractID:  contractID,
			groupID:     groupID,
			cpCodeID:    d.Id(),
			updatedName: name,
			config:      config,
		}); err != nil {
			if errors.Is(err, ErrCPCodeUpdateTimeout) {
				return append(tf.DiagWarningf("%s", err), tf.DiagWarningf("Resource has been updated, but the change is still ongoing on the server")...)
			}
			return diag.FromErr(err)
		}

		return resourceCPCodeRead(ctx, d, m)
	}
}

func resourceCPCodeImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceCPCodeImport")
	logger.Debugf("Import CP Code")

	parts := strings.Split(d.Id(), ",")

	if len(parts) < 3 {
		return nil, fmt.Errorf("comma-separated list of CP code ID, contract ID and group ID has to be supplied in import: %s", d.Id())
	}
	if parts[0] == "" {
		return nil, errors.New("CP Code is a mandatory parameter")
	}
	cpCodeID := parts[0]
	contractID := str.AddPrefix(parts[1], "ctr_")
	groupID := str.AddPrefix(parts[2], "grp_")

	cpCodeResp, err := meta.Client().GetPAPI().GetCPCode(ctx, papi.GetCPCodeRequest{
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
		return nil, fmt.Errorf(errCPCodeNoProductID, cpCode.Name)
	}
	if err := d.Set("product_id", cpCode.ProductIDs[0]); err != nil {
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
		"contract_id",
		"group_id",
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

type waitForCPCodeNameUpdateOpts struct {
	client      papi.PAPI
	contractID  string
	groupID     string
	cpCodeID    string
	updatedName string
	config      cpCodeResourceConfig
}

func waitForCPCodeNameUpdate(ctx context.Context, opts waitForCPCodeNameUpdateOpts) error {
	req := papi.GetCPCodeRequest{CPCodeID: opts.cpCodeID, ContractID: opts.contractID, GroupID: opts.groupID}
	CPCodeResp, err := opts.client.GetCPCode(ctx, req)
	if err != nil {
		return err
	}

	for CPCodeResp.CPCode.Name != opts.updatedName {
		select {
		case <-time.After(tf.MaxDuration(opts.config.updatePollInterval, opts.config.updatePollMinimum)):
			CPCodeResp, err = opts.client.GetCPCode(ctx, req)
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

var errCPCodeNoProductID = "the CP code named `%s` already exists, but does not have a PAPI-supported product ID, so it cannot be managed by Terraform"
