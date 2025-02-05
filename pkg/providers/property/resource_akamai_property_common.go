package property

// This file contains functions removed from resource_akamai_property.go that are still referenced elsewhere

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
)

func getGroup(ctx context.Context, client papi.PAPI, groupID string) (*papi.Group, error) {
	log := log.FromContext(ctx)
	log.Debugf("Fetching groups")

	res, err := client.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFetchingGroups, err.Error())
	}
	groupID = str.AddPrefix(groupID, "grp_")

	var group *papi.Group
	var groupFound bool
	for _, g := range res.Groups.Items {
		if g.GroupID == groupID {
			group = g
			groupFound = true
			break
		}
	}
	if !groupFound {
		return nil, fmt.Errorf("%w: %s", ErrGroupNotFound, groupID)
	}
	log.Debugf("Group found: %s", group.GroupID)
	return group, nil
}

func getContract(ctx context.Context, client papi.PAPI, contractID string) (*papi.Contract, error) {
	log := log.FromContext(ctx)
	log.Debugf("Fetching contract")

	res, err := client.GetContracts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFetchingContracts, err.Error())
	}
	var contract *papi.Contract
	var contractFound bool
	for _, c := range res.Contracts.Items {
		if c.ContractID == contractID {
			contract = c
			contractFound = true
			break
		}
	}
	if !contractFound {
		return nil, fmt.Errorf("%w: %s", ErrContractNotFound, contractID)
	}

	log.Debugf("Contract found: %s", contract.ContractID)
	return contract, nil
}

func getProduct(ctx context.Context, client papi.PAPI, productID, contractID string) (*papi.ProductItem, error) {
	if contractID == "" {
		return nil, ErrNoContractProvided
	}

	log := log.FromContext(ctx)
	log.Debugf("Fetching product")

	res, err := client.GetProducts(ctx, papi.GetProductsRequest{ContractID: contractID})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrProductFetch, err.Error())
	}
	var productFound bool
	var product papi.ProductItem
	for _, p := range res.Products.Items {
		if p.ProductID == productID {
			product = p
			productFound = true
			break
		}
	}
	if !productFound {
		return nil, fmt.Errorf("%w: %s", ErrProductNotFound, productID)
	}

	log.Debugf("Product found: %s", product.ProductID)
	return &product, nil
}

func findProperty(ctx context.Context, name string, meta meta.Meta) (*papi.Property, error) {
	client := Client(meta)
	results, err := client.SearchProperties(ctx, papi.SearchRequest{Key: papi.SearchKeyPropertyName, Value: name})
	if err != nil {
		return nil, err
	}
	if len(results.Versions.Items) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPropertyNotFound, name)
	}

	property, err := client.GetProperty(ctx, papi.GetPropertyRequest{
		ContractID: results.Versions.Items[0].ContractID,
		GroupID:    results.Versions.Items[0].GroupID,
		PropertyID: results.Versions.Items[0].PropertyID,
	})
	if err != nil {
		return nil, err
	}
	if len(property.Properties.Items) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPropertyNotFound, name)
	}
	return property.Properties.Items[0], nil
}
