package property

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyIncludeActivation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyIncludeActivationRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the contract under which the include is activated",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the group under which the include is activated",
			},
			"include_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies include targeted with activation",
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The network for which the activation is to be found",
				ValidateDiagFunc: tools.ValidateStringInSlice([]string{string(papi.ActivationNetworkProduction), string(papi.ActivationNetworkStaging)}),
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The include version targeted with activation",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the include targeted with activation",
			},
			"note": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Log message assigned to the activation",
			},
			"notify_emails": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of email addresses to notify when the activation status changes",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

type includeActivationAttrs struct {
	contractID string
	groupID    string
	includeID  string
	network    string
}

func dataPropertyIncludeActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	log := meta.Log("PAPI", "dataPropertyIncludeActivationRead")
	log.Debug("Reading Property Include Activation")

	attrs, err := getIncludeActivationAttrs(d)
	if err != nil {
		return diag.Errorf("getIncludeActivationAttrs error: %s", err)
	}

	activations, err := client.ListIncludeActivations(ctx, papi.ListIncludeActivationsRequest{
		IncludeID:  attrs.includeID,
		ContractID: attrs.contractID,
		GroupID:    attrs.groupID,
	})
	if err != nil {
		return diag.Errorf("could not list include activations: %s", err)
	}

	filteredActivations := filterIncludeActivationsByNetwork(activations.Activations.Items, attrs.network)
	latestActivation, err := findLatestIncludeActivation(filteredActivations)
	if err != nil {
		log.Info(fmt.Sprintf("%s: there is no active version on %s network", err, attrs.network))
	}

	attributes := createIncludeActivationAttrs(latestActivation)

	if err = tools.SetAttrs(d, attributes); err != nil {
		return diag.Errorf("could not set attributes: %s", err)
	}
	d.SetId(attrs.includeID + ":" + attrs.network)

	return nil
}

func createIncludeActivationAttrs(latestActivation *papi.IncludeActivation) map[string]interface{} {
	var version, name, note string
	var notifyEmails []string

	if latestActivation != nil {
		version = strconv.Itoa(latestActivation.IncludeVersion)
		name = latestActivation.IncludeName
		note = latestActivation.Note
		notifyEmails = latestActivation.NotifyEmails
	}

	return map[string]interface{}{
		"version":       version,
		"name":          name,
		"note":          note,
		"notify_emails": notifyEmails,
	}
}

// findLatestIncludeActivation finds the latest activation of type `ACTIVATE` with status `ACTIVE` or `PENDING`.
// If it encounters activation of type `DEACTIVATE` with status `ACTIVE` first or does not find any activation of type
// `ACTIVATE` with `ACTIVE` status, it returns nil
func findLatestIncludeActivation(activations []papi.IncludeActivation) (*papi.IncludeActivation, error) {
	if len(activations) == 0 {
		return nil, ErrNoLatestIncludeActivation
	}

	sort.Slice(activations, func(i, j int) bool {
		return activations[i].UpdateDate > activations[j].UpdateDate
	})

	for _, activation := range activations {
		if activation.ActivationType == papi.ActivationTypeActivate &&
			(activation.Status == papi.ActivationStatusActive || activation.Status == papi.ActivationStatusPending) {
			return &activation, nil
		}
		if activation.ActivationType == papi.ActivationTypeDeactivate &&
			(activation.Status == papi.ActivationStatusActive || activation.Status == papi.ActivationStatusPending) {
			return nil, ErrNoLatestIncludeActivation
		}
	}

	return nil, ErrNoLatestIncludeActivation
}

func filterIncludeActivationsByNetwork(activations []papi.IncludeActivation, network string) []papi.IncludeActivation {
	var filteredActivations []papi.IncludeActivation
	for _, activation := range activations {
		if string(activation.Network) == network {
			filteredActivations = append(filteredActivations, activation)
		}
	}

	return filteredActivations
}

func getIncludeActivationAttrs(d *schema.ResourceData) (*includeActivationAttrs, error) {
	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return nil, fmt.Errorf("could not get `contract_id` attribute: %s", err)
	}

	groupID, err := tools.GetStringValue("group_id", d)
	if err != nil {
		return nil, fmt.Errorf("could not get `group_id` attribute: %s", err)
	}

	includeID, err := tools.GetStringValue("include_id", d)
	if err != nil {
		return nil, fmt.Errorf("could not get `include_id` attribute: %s", err)
	}

	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return nil, fmt.Errorf("could not get `network` attribute: %s", err)
	}

	return &includeActivationAttrs{
		contractID: contractID,
		groupID:    groupID,
		includeID:  includeID,
		network:    network,
	}, nil
}
