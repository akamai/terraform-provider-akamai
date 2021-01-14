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
		p.mtx.Lock() // Serialize any requests which may impact injected dependencies
		defer p.mtx.Unlock()

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
			p.mtx.Lock() // Serialize any requests which may impact injected dependencies
			defer p.mtx.Unlock()

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

	logger := meta.Log("IAM", opName)
	logger = logger.WithFields(log.Fields{"operation_id": meta.OperationID()})

	p.SetIAM(iam.Client(meta.Session()))
	p.SetCache(metaCache{p, meta})

	return log.NewContext(ctx, logger)
}
