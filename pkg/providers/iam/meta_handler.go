package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Alias for any TF CRUD operation function having this common signature
type tfCRUDFunc = func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics

// Compose a TF CRUD entry point function that processes the meta and invokes the impl with no meta
func (p *provider) tfCRUD(opName string, impl tfCRUDFunc) tfCRUDFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		ctx = p.handleMeta(ctx, m, opName)

		p.log(ctx).Debugf("Start of Terraform action")
		defer p.log(ctx).Debugf("End of Terraform action")

		return impl(ctx, d, nil)
	}
}

// Compose a schema.ResourceImporter that processes the meta and invokes the impl with no meta
func (p *provider) tfImporter(opName string, impl schema.StateContextFunc) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
			ctx = p.handleMeta(ctx, m, opName)

			p.log(ctx).Debugf("Start of Terraform action")
			defer p.log(ctx).Debugf("End of Terraform action")

			return impl(ctx, d, nil)
		},
	}
}

// Accept dependencies from Meta and setup the context. Does nothing when meta is nil
func (p *provider) handleMeta(ctx context.Context, m interface{}, opName string) context.Context {
	if m == nil {
		return ctx
	}

	meta := akamai.Meta(m)

	if p.assertMeta == nil {
		p.assertMeta = mkAssertMeta(meta)
	}

	p.assertMeta(meta)

	logger := meta.Log("IAM", opName)
	logger = logger.WithFields(log.Fields{"operation_id": meta.OperationID()})

	p.SetIAM(iam.Client(meta.Session()))
	p.SetCache(metaCache{p, meta})

	return log.NewContext(ctx, logger)
}

// Build a function that verifies the assumption that we receive exactly one meta value
func mkAssertMeta(originalMeta akamai.OperationMeta) func(akamai.OperationMeta) {
	if originalMeta == nil {
		panic("BUG: originalMeta can't be nil")
	}

	return func(newMeta akamai.OperationMeta) {
		if newMeta != originalMeta {
			panic("BUG: Received a new and different meta (invariant broken)")
		}
	}
}
