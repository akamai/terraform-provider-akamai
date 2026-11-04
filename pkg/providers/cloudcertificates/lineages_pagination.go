package cloudcertificates

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v12/internal/pagination"
)

func lineagesPaginationOptions(
	client cloudcertificates.CloudCertificates,
	request cloudcertificates.ListLineagesRequest,
	limit *int64,
) pagination.Options[cloudcertificates.Lineage, cloudcertificates.ListLineagesResponse] {
	return pagination.Options[cloudcertificates.Lineage, cloudcertificates.ListLineagesResponse]{
		Limit:    limit,
		PageSize: cloudcertificates.MaxListLineagesPageSize,
		Fetch: func(ctx context.Context, cursor string, pageSize int) (*cloudcertificates.ListLineagesResponse, error) {
			pageRequest := request
			pageRequest.After = cursor
			pageRequest.PageSize = pageSize
			return client.ListLineages(ctx, pageRequest)
		},
		Items:      func(page *cloudcertificates.ListLineagesResponse) []cloudcertificates.Lineage { return page.Lineages },
		NextCursor: func(page *cloudcertificates.ListLineagesResponse) *string { return page.NextCursor },
	}
}
