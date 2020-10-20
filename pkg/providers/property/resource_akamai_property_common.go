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
	groupID, err = tools.AddPrefix(groupID, "grp_")
	if err != nil {
		return nil, err
	}

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
	contractID, err = tools.AddPrefix(contractID, "ctr_")
	if err != nil {
		return nil, err
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
	productID, err = tools.AddPrefix(productID, "prd_")
	if err != nil {
		return nil, err
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

	logger.Debugf("Product found: %s", product.ProductID)
	return &product, nil
}

func extractOptions(options *schema.Set) (map[string]interface{}, error) {
	optv := make(map[string]interface{})
	for _, option := range options.List() {
		optionMap, ok := option.(map[string]interface{})
		if !ok {
			continue
		}
		if val, ok := optionMap["value"].(string); ok && val != "" {
			optv[optionMap["key"].(string)] = convertString(val)
			continue
		}
		vals, ok := optionMap["values"]
		if !ok {
			continue
		}
		valsSet, ok := vals.(*schema.Set)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "values", "*schema.Set")
		}
		if valsSet.Len() == 0 {
			optv[optionMap["key"].(string)] = convertString(optionMap["value"].(string))
			continue
		}
		if valsSet.Len() > 0 {
			op := make([]interface{}, 0)
			for _, v := range vals.(*schema.Set).List() {
				op = append(op, convertString(v.(string)))
			}

			optv[optionMap["key"].(string)] = op
		}
	}
	return optv, nil
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

func extractRules(drules *schema.Set) ([]papi.Rules, error) {

	var rules []papi.Rules
	for _, v := range drules.List() {
		rule := papi.Rules{Name: "default"}
		vv, ok := v.(map[string]interface{})
		if ok {
			rule.Name = vv["name"].(string)
			rule.Comment = vv["comment"].(string)

			criteriaMustSatisfy, ok := vv["criteria_match"]
			if ok {
				if criteriaMustSatisfy.(string) == "all" {
					rule.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAll
				}

				if criteriaMustSatisfy.(string) == "any" {
					rule.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAny
				}
			}
			behaviors, ok := vv["behavior"]
			if ok {
				for _, behavior := range behaviors.(*schema.Set).List() {
					behaviorMap, ok := behavior.(map[string]interface{})
					if ok {
						newBehavior := papi.RuleBehavior{}
						newBehavior.Name = behaviorMap["name"].(string)
						behaviorOptions, ok := behaviorMap["option"]
						if ok {
							opts, err := extractOptions(behaviorOptions.(*schema.Set))
							if err != nil {
								return nil, err
							}
							newBehavior.Options = opts
						}
						rule.Behaviors = mergeBehaviors(rule.Behaviors, newBehavior)
					}
				}
			}

			criterias, ok := vv["criteria"]
			if ok {
				for _, criteria := range criterias.(*schema.Set).List() {
					criteriaMap, ok := criteria.(map[string]interface{})
					if ok {
						newCriteria := papi.RuleBehavior{}
						newCriteria.Name = criteriaMap["name"].(string)
						criteriaOptions, ok := criteriaMap["option"]
						if ok {
							crit, err := extractOptions(criteriaOptions.(*schema.Set))
							if err != nil {
								return nil, err
							}
							newCriteria.Options = crit
						}
						rule.Criteria = append(rule.Criteria, newCriteria)
					}
				}
			}

			variables, ok := vv["variable"]
			if ok {
				for _, variable := range variables.(*schema.Set).List() {
					variableMap, ok := variable.(map[string]interface{})
					if ok {
						newVariable := papi.RuleVariable{}
						newVariable.Name = variableMap["name"].(string)
						newVariable.Description = variableMap["description"].(string)
						newVariable.Value = variableMap["value"].(string)
						newVariable.Hidden = variableMap["hidden"].(bool)
						newVariable.Sensitive = variableMap["sensitive"].(bool)
						rule.Variables = addVariable(rule.Variables, newVariable)
					}
				}
			}

			childRules, ok := vv["rule"]
			if ok && childRules.(*schema.Set).Len() > 0 {
				rules, err := extractRules(childRules.(*schema.Set))
				if err != nil {
					return nil, err
				}
				for _, newRule := range rules {
					rule.Children = append(rule.Children, newRule)
				}
			}
		}
		rules = append(rules, rule)
	}

	return rules, nil
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

func mergeBehaviors(old []papi.RuleBehavior, new papi.RuleBehavior) []papi.RuleBehavior {
	for i := range old {
		if new.Name == "cpCode" || new.Name == "origin" {
			if old[i].Name == new.Name {
				old[i].Options = mergeOptions(old[i].Options, new.Options)
				return old
			}
		}
	}

	return append(old, new)
}

// MergeOptions merges the given options with the existing options
func mergeOptions(old, new papi.RuleOptionsMap) papi.RuleOptionsMap {
	options := make(papi.RuleOptionsMap)
	for k, v := range old {
		options[k] = v
	}
	for k, v := range new {
		options[k] = v
	}
	return options
}

func addVariable(old []papi.RuleVariable, new papi.RuleVariable) []papi.RuleVariable {
	for i := range old {
		if old[i].Name == new.Name {
			old[i] = new
			return old
		}
	}

	return append(old, new)
}
