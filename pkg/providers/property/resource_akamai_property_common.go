package property

// This file contains functions removed from resource_akamai_property.go that are still referenced elsewhere

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func getGroup(ctx context.Context, d *schema.ResourceData, meta akamai.OperationMeta) (*papi.Group, error) {
	logger := meta.Log("PAPI", "getGroup")
	client := inst.Client(meta)
	logger.Debugf("Fetching groups")
	groupID, err := tools.GetStringValue("group", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, ErrNoGroupProvided
	}
	res, err := client.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFetchingGroups, err.Error())
	}
	groupID = tools.AddPrefix(groupID, "grp_")

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
	logger.Debugf("Group found: %s", group.GroupID)
	return group, nil
}

func getContract(ctx context.Context, d *schema.ResourceData, meta akamai.OperationMeta) (*papi.Contract, error) {
	logger := meta.Log("PAPI", "getContract")
	client := inst.Client(meta)
	logger.Debugf("Fetching contract")
	contractID, err := tools.GetStringValue("contract", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, ErrNoContractProvided
	}
	res, err := client.GetContracts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFetchingContracts, err.Error())
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
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

	logger.Debugf("Contract found: %s", contract.ContractID)
	return contract, nil
}

func getProduct(ctx context.Context, d *schema.ResourceData, contractID string, meta akamai.OperationMeta) (*papi.ProductItem, error) {
	logger := meta.Log("PAPI", "getProduct")
	client := inst.Client(meta)
	if contractID == "" {
		return nil, ErrNoContractProvided
	}
	logger.Debugf("Fetching product")
	productID, err := tools.GetStringValue("product", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, ErrNoProductProvided
	}
	res, err := client.GetProducts(ctx, papi.GetProductsRequest{ContractID: contractID})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrProductFetch, err.Error())
	}
	productID = tools.AddPrefix(productID, "prd_")
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

	logger.Debugf("Product found: %s", product.ProductID)
	return &product, nil
}

func convertString(v string) interface{} {
	if f1, err := strconv.ParseFloat(v, 64); err == nil {
		return f1
	}
	// FIXME: execution will never reach this as every int representation will be captured by ParseFloat() above
	// this should either be moved above ParseFloat block or removed
	if f2, err := strconv.ParseInt(v, 10, 64); err == nil {
		return f2
	}
	if f3, err := strconv.ParseBool(v); err == nil {
		return f3
	}
	return v
}

func findProperty(ctx context.Context, name string, meta akamai.OperationMeta) (*papi.Property, error) {
	client := inst.Client(meta)
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