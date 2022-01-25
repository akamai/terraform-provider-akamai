package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEdgeworkerActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeworkerActivationCreate,
		ReadContext:   resourceEdgeworkerActivationRead,
		UpdateContext: resourceEdgeworkerActivationUpdate,
		DeleteContext: resourceEdgeworkerActivationDelete,
		Schema:        resourceEdgeworkerActivationSchema(),
		Timeouts: &schema.ResourceTimeout{
			Default: &edgeworkerActivationResourceTimeout,
		},
	}
}

func resourceEdgeworkerActivationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"edgeworker_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Id of the EdgeWorker to activate",
		},
		"version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The version of EdgeWorker to activate",
		},
		"network": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				string(edgeworkers.ActivationNetworkStaging),
				string(edgeworkers.ActivationNetworkProduction),
			}, false)),
			Description: "The network on which the version will be activated",
		},
		"activation_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "A unique identifier of the activation",
		},
	}
}

var (
	activationStatusComplete            = "COMPLETE"
	activationStatusPresubmit           = "PRESUBMIT"
	activationStatusPending             = "PENDING"
	activationStatusInProgress          = "IN_PROGRESS"
	errorCodeVersionAlreadyActive       = "EW1021"
	activationPollMinimum               = time.Minute
	activationPollInterval              = activationPollMinimum
	edgeworkerActivationResourceTimeout = time.Minute * 30
)

func resourceEdgeworkerActivationCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkerActivationCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Activating edgeworker")

	edgeworkerID, err := tools.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tools.GetStringValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	activation, err := client.ActivateVersion(ctx, edgeworkers.ActivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		ActivateVersion: edgeworkers.ActivateVersion{
			Network: edgeworkers.ActivationNetwork(network),
			Version: version,
		},
	})

	if err != nil {
		e := err.(*edgeworkers.Error)
		if e.ErrorCode != errorCodeVersionAlreadyActive {
			return diag.Errorf("%s create: %s", ErrEdgeworkerActivation, err.Error())
		}

		resp, err := client.ListActivations(ctx, edgeworkers.ListActivationsRequest{
			EdgeWorkerID: edgeworkerID,
			Version:      version,
		})
		if err != nil {
			return diag.Errorf("%s create: %s", ErrEdgeworkerActivation, err.Error())
		}

		activations := filterActivationsByNetwork(resp.Activations, network)
		activation = &activations[0]
	} else {
		activation, err = waitForEdgeworkerActivation(ctx, client, edgeworkerID, activation.ActivationID)
		if err != nil {
			return diag.Errorf("%s create: %s", ErrEdgeworkerActivation, err.Error())
		}
	}

	if err := rd.Set("activation_id", activation.ActivationID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	rd.SetId(fmt.Sprintf("%d:%s", edgeworkerID, network))

	return resourceEdgeworkerActivationRead(ctx, rd, m)
}

func resourceEdgeworkerActivationRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkerActivationRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading edgeworker activations")

	edgeworkerID, err := tools.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListActivations(ctx, edgeworkers.ListActivationsRequest{
		EdgeWorkerID: edgeworkerID,
	})
	if err != nil {
		return diag.Errorf("%s read: %s", ErrEdgeworkerActivation, err.Error())
	}

	activations := filterActivationsByNetwork(resp.Activations, network)
	activation := &activations[0]

	if activation.Status != activationStatusComplete {
		return diag.Errorf("%s read: activation (%d) status is not '%s'", ErrEdgeworkerActivation, activation.ActivationID, activationStatusComplete)
	}

	if err := rd.Set("version", activation.Version); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	if err := rd.Set("activation_id", activation.ActivationID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceEdgeworkerActivationUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkerActivationUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Updating edgeworker activation")

	edgeworkerID, err := tools.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	newVersion, err := tools.GetStringValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	activation, err := client.ActivateVersion(ctx, edgeworkers.ActivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		ActivateVersion: edgeworkers.ActivateVersion{
			Network: edgeworkers.ActivationNetwork(network),
			Version: newVersion,
		},
	})
	if err != nil {
		return diag.Errorf("%s update: %s", ErrEdgeworkerActivation, err.Error())
	}

	activation, err = waitForEdgeworkerActivation(ctx, client, edgeworkerID, activation.ActivationID)
	if err != nil {
		return diag.Errorf("%s update: %s", ErrEdgeworkerActivation, err.Error())
	}

	if err := rd.Set("activation_id", activation.ActivationID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	rd.SetId(fmt.Sprintf("%d:%s", edgeworkerID, network))

	return resourceEdgeworkerActivationRead(ctx, rd, m)
}

func resourceEdgeworkerActivationDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkerActivationDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Deactivating edgeworker")

	edgeworkerID, err := tools.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tools.GetStringValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	deactivation, err := client.DeactivateVersion(ctx, edgeworkers.DeactivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		DeactivateVersion: edgeworkers.DeactivateVersion{
			Version: version,
			Network: edgeworkers.ActivationNetwork(network),
		},
	})
	if err != nil {
		return diag.Errorf("%s: %s", ErrEdgeworkerDeactivation, err.Error())
	}

	if err := waitForEdgeworkerDeactivation(ctx, client, edgeworkerID, deactivation.DeactivationID); err != nil {
		return diag.Errorf("%s: %s", ErrEdgeworkerDeactivation, err.Error())
	}

	rd.SetId("")
	return nil
}

func waitForEdgeworkerActivation(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID, activationID int) (*edgeworkers.Activation, error) {
	activation, err := client.GetActivation(ctx, edgeworkers.GetActivationRequest{
		EdgeWorkerID: edgeworkerID,
		ActivationID: activationID,
	})
	if err != nil {
		return nil, err
	}
	for activation.Status != activationStatusComplete {
		if activation.Status != activationStatusPresubmit && activation.Status != activationStatusPending && activation.Status != activationStatusInProgress {
			return nil, ErrEdgeworkerActivationFailure
		}
		select {
		case <-time.After(tools.MaxDuration(activationPollInterval, activationPollMinimum)):
			activation, err = client.GetActivation(ctx, edgeworkers.GetActivationRequest{
				EdgeWorkerID: edgeworkerID,
				ActivationID: activationID,
			})
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, ErrEdgeworkerActivationTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil, ErrEdgeworkerActivationCanceled
			}
			return nil, fmt.Errorf("%v: %w", ErrEdgeworkerActivationContextTerminated, ctx.Err())
		}
	}
	return activation, nil
}

func waitForEdgeworkerDeactivation(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID, deactivationID int) error {
	deactivation, err := client.GetDeactivation(ctx, edgeworkers.GetDeactivationRequest{
		EdgeWorkerID:   edgeworkerID,
		DeactivationID: deactivationID,
	})
	if err != nil {
		return err
	}
	for deactivation.Status != activationStatusComplete {
		if deactivation.Status != activationStatusPresubmit && deactivation.Status != activationStatusPending && deactivation.Status != activationStatusInProgress {
			return ErrEdgeworkerDeactivationFailure
		}
		select {
		case <-time.After(tools.MaxDuration(activationPollInterval, activationPollMinimum)):
			deactivation, err = client.GetDeactivation(ctx, edgeworkers.GetDeactivationRequest{
				EdgeWorkerID:   edgeworkerID,
				DeactivationID: deactivationID,
			})
			if err != nil {
				return err
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return ErrEdgeworkerDeactivationTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return ErrEdgeworkerDeactivationCanceled
			}
			return fmt.Errorf("%v: %w", ErrEdgeworkerDeactivationContextTerminated, ctx.Err())
		}
	}
	return nil
}

func filterActivationsByNetwork(acts []edgeworkers.Activation, net string) (activations []edgeworkers.Activation) {
	for _, act := range acts {
		if act.Network == string(net) {
			activations = append(activations, act)
		}
	}
	return activations
}
