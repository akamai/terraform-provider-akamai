package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/providers/property/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPCode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCPCodeRead,
		Schema: map[string]*schema.Schema{
			"cp_code_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"cp_code_name", "cp_code_id"},
			},
			"cp_code_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"cp_code_name", "cp_code_id"},
				DiffSuppressFunc: tf.FieldPrefixSuppress("cpc_"),
			},
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataCPCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	log := meta.Log("PAPI", "dataCPCodeRead")
	log.Debug("Read CP Code")
	client := meta.Client().GetPAPI()
	var cpCodeName, cpCodeID, groupID, contractID string
	var err error

	if v, ok := d.GetOk("cp_code_name"); ok {
		cpCodeName = v.(string)
	}

	if v, ok := d.GetOk("cp_code_id"); ok {
		cpCodeID = v.(string)
	}

	if groupID, err = tf.GetStringValue("group_id", d); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", groupID); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	if contractID, err = tf.GetStringValue("contract_id", d); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	var cpCode papi.CPCode
	if cpCodeID != "" {
		cpCodeResp, err := client.GetCPCode(ctx, papi.GetCPCodeRequest{
			CPCodeID:   cpCodeID,
			ContractID: contractID,
			GroupID:    groupID,
		})
		if err != nil {
			if errors.Is(err, papi.ErrNotFound) {
				return diag.FromErr(fmt.Errorf("%w: %w", ErrLookingUpCPCodeByID, err))
			}
			return diag.FromErr(fmt.Errorf("could not load CP code: %w", err))
		}

		cpCode = cpCodeResp.CPCode
		if err := d.Set("cp_code_name", cpCodeResp.CPCode.Name); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		cpCode, err = tools.FindCPCodeByName(ctx, client, cpCodeName, contractID, groupID)
		if err != nil {
			if errors.Is(err, tools.ErrCPCodeNotFound) {
				return diag.FromErr(fmt.Errorf("%w: %w", ErrLookingUpCPCodeByName, err))
			}
			return diag.FromErr(fmt.Errorf("could not load CP codes: %w", err))
		}
		if err := d.Set("cp_code_id", cpCode.ID); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
	}

	if err := d.Set("product_ids", cpCode.ProductIDs); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("created_date", cpCode.CreatedDate); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strings.TrimPrefix(cpCode.ID, cpCodePrefix))
	log.Debugf("Read CP Code: %+v", cpCode)
	return nil
}
