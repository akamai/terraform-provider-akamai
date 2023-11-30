package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/collections"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEdgeworkersActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeworkersActivationCreate,
		ReadContext:   resourceEdgeworkersActivationRead,
		UpdateContext: resourceEdgeworkersActivationUpdate,
		DeleteContext: resourceEdgeworkersActivationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEdgeworkersActivationImport,
		},
		Schema: resourceEdgeworkersActivationSchema(),
		Timeouts: &schema.ResourceTimeout{
			Delete:  &edgeworkersActivationResourceDeleteTimeout,
			Default: &edgeworkersActivationResourceDefaultTimeout,
		},
		CustomizeDiff: checkEdgeworkerExistsOnDiff,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourceEdgeworkersActivationV0().CoreConfigSchema().ImpliedType(),
			Upgrade: timeouts.MigrateToExplicit(),
		}},
	}
}

func resourceEdgeworkersActivationSchema() map[string]*schema.Schema {
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
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: tf.ValidateStringInSlice(validEdgeworkerActivationNetworks),
			Description:      "The network on which the version will be activated",
		},
		"activation_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "A unique identifier of the activation",
		},
		"note": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Assigns a log message to the activation request",
			DiffSuppressFunc: suppressNoteFieldForEdgeWorkersActivation,
		},
		"timeouts": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Enables to set timeout for processing",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: timeouts.ValidateDurationFormat,
					},
					"delete": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: timeouts.ValidateDurationFormat,
					},
				},
			},
		},
	}
}

const (
	stagingNetwork             = "STAGING"
	productionNetwork          = "PRODUCTION"
	activationStatusComplete   = "COMPLETE"
	activationStatusPresubmit  = "PRESUBMIT"
	activationStatusPending    = "PENDING"
	activationStatusInProgress = "IN_PROGRESS"
)

const (
	errorCodeVersionIsBeingDeactivated = "EW1031"
	errorCodeVersionAlreadyDeactivated = "EW1032"
)

var validEdgeworkerActivationNetworks = []string{stagingNetwork, productionNetwork}

var (
	activationPollMinimum                       = time.Minute
	activationPollInterval                      = activationPollMinimum
	edgeworkersActivationResourceDefaultTimeout = time.Minute * 30
	edgeworkersActivationResourceDeleteTimeout  = time.Minute * 60
)

const timeLayout = time.RFC3339

func resourceEdgeworkersActivationCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkersActivationCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Activating edgeworker")

	return upsertActivation(ctx, rd, m, client)
}

func resourceEdgeworkersActivationRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkersActivationRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading edgeworker activations")

	edgeworkerID, err := tf.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	activation, err := getCurrentActivation(ctx, client, edgeworkerID, network, false)
	if err != nil {
		if errors.Is(err, ErrEdgeworkerNoCurrentActivation) {
			return diag.Errorf(`%s read: no version active on network '%s' for edgeworker with id=%d`, ErrEdgeworkerActivation, network, edgeworkerID)
		}
		return diag.Errorf("%s read: %s", ErrEdgeworkerActivation, err)
	}

	if err := rd.Set("version", activation.Version); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	if err := rd.Set("activation_id", activation.ActivationID); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	if err := rd.Set("note", activation.Note); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceEdgeworkersActivationUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkersActivationUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Updating edgeworker activation")

	if !rd.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	return upsertActivation(ctx, rd, m, client)
}

func resourceEdgeworkersActivationDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkersActivationDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Deactivating edgeworker")

	edgeworkerID, err := tf.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tf.GetStringValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	note, err := tf.GetStringValue("note", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	findLatestDeactivation := false
	deactivation, err := client.DeactivateVersion(ctx, edgeworkers.DeactivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		DeactivateVersion: edgeworkers.DeactivateVersion{
			Version: version,
			Network: edgeworkers.ActivationNetwork(network),
			Note:    note,
		},
	})
	if err != nil {
		if errors.Is(err, edgeworkers.ErrVersionAlreadyDeactivated) {
			logger.Info(fmt.Sprintf("Version '%s' has already been deactivated on network '%s' for edgeworker with id=%d. Removing from state", version, network, edgeworkerID))
			return nil
		}
		if errors.Is(err, edgeworkers.ErrVersionBeingDeactivated) {
			findLatestDeactivation = true
		} else {
			return diag.Errorf("%s: %s", ErrEdgeworkerDeactivation, err)
		}
	}

	if findLatestDeactivation {
		deactivations, err := getDeactivationsByVersionAndNetwork(ctx, client, edgeworkerID, version, network)
		if err != nil {
			return diag.Errorf("%s: %s", ErrEdgeworkerDeactivation, err)
		}
		deactivation = &deactivations[0]
	}

	if _, err := waitForEdgeworkerDeactivation(ctx, client, edgeworkerID, deactivation.DeactivationID); err != nil {
		if errors.Is(err, ErrEdgeworkerDeactivationTimeout) {
			rd.SetId("")
			return append(tf.DiagWarningf("%s: %s", ErrEdgeworkerDeactivation, err), tf.DiagWarningf("Resource has been removed from the state, but deactivation is still ongoing on the server")...)
		}
		return diag.Errorf("%s: %s", ErrEdgeworkerDeactivation, err)
	}

	return nil
}

func resourceEdgeworkersActivationImport(_ context.Context, rd *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Edgeworkers", "resourceEdgeworkersActivationImport")

	logger.Debug("Importing edgeworker")

	parts := strings.Split(rd.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%s import: invalid import id '%s' - colon-separated list of edgeworker ID and network has to be supplied", ErrEdgeworkerActivation, rd.Id())
	}

	edgeworkerID, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%s import: edgeworker id must be an integer, got '%s'", ErrEdgeworkerActivation, parts[0])
	}

	network := parts[1]
	if !collections.StringInSlice(validEdgeworkerActivationNetworks, network) {
		return nil, fmt.Errorf("%s import: network must be 'STAGING' or 'PRODUCTION', got '%s'", ErrEdgeworkerActivation, network)
	}

	if err := rd.Set("edgeworker_id", edgeworkerID); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	if err := rd.Set("network", network); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	return []*schema.ResourceData{rd}, nil
}

func upsertActivation(ctx context.Context, rd *schema.ResourceData, m interface{}, client edgeworkers.Edgeworkers) diag.Diagnostics {
	edgeworkerID, err := tf.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tf.GetStringValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	versionsResp, err := client.ListEdgeWorkerVersions(ctx, edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeworkerID,
	})
	if err != nil {
		return diag.Errorf("%s: %s", ErrEdgeworkerActivation, err.Error())
	}
	if !versionExists(version, versionsResp.EdgeWorkerVersions) {
		return diag.Errorf(`%s: version '%s' is not valid for edgeworker with id=%d`, ErrEdgeworkerActivation, version, edgeworkerID)
	}

	currentActivation, err := getCurrentActivation(ctx, client, edgeworkerID, network, true)
	if err != nil && !errors.Is(err, ErrEdgeworkerNoCurrentActivation) {
		return diag.Errorf("%s: %s", ErrEdgeworkerActivation, err.Error())
	}

	if currentActivation != nil && currentActivation.Version == version {
		rd.SetId(fmt.Sprintf("%d:%s", edgeworkerID, network))
		return resourceEdgeworkersActivationRead(ctx, rd, m)
	}

	note, err := tf.GetStringValue("note", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	activation, err := client.ActivateVersion(ctx, edgeworkers.ActivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		ActivateVersion: edgeworkers.ActivateVersion{
			Network: edgeworkers.ActivationNetwork(network),
			Version: version,
			Note:    note,
		},
	})

	if err != nil {
		return diag.Errorf("%s: %s", ErrEdgeworkerActivation, err.Error())
	}

	if _, err := waitForEdgeworkerActivation(ctx, client, edgeworkerID, activation.ActivationID); err != nil {
		return diag.Errorf("%s: %s", ErrEdgeworkerActivation, err.Error())
	}

	rd.SetId(fmt.Sprintf("%d:%s", edgeworkerID, network))
	return resourceEdgeworkersActivationRead(ctx, rd, m)
}

func getCurrentActivation(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID int, network string, waitForDeactivation bool) (*edgeworkers.Activation, error) {
	activationsResp, err := client.ListActivations(ctx, edgeworkers.ListActivationsRequest{
		EdgeWorkerID: edgeworkerID,
	})
	if err != nil {
		return nil, err
	}

	activations := sortActivationsByDate(filterActivationsByNetwork(activationsResp.Activations, network))
	if len(activations) == 0 {
		return nil, ErrEdgeworkerNoCurrentActivation
	}

	latestActivation := &activations[0]
	if !statusOngoingOrReady(latestActivation.Status) {
		return nil, ErrEdgeworkerNoCurrentActivation
	}

	if statusOngoing(latestActivation.Status) {
		latestActivation, err = waitForEdgeworkerActivation(ctx, client, edgeworkerID, latestActivation.ActivationID)
		if err != nil {
			return nil, err
		}
		return latestActivation, nil
	}

	latestDeactivation, err := getLatestCompletedDeactivation(ctx, client, edgeworkerID, latestActivation.Version, network, waitForDeactivation)
	if err != nil {
		if errors.Is(err, ErrEdgeworkerNoLatestDeactivation) {
			return latestActivation, nil
		}
		return nil, err
	}

	isDeactivated, err := wasDeactivationLater(latestActivation, latestDeactivation)
	if err != nil {
		return nil, err
	}

	if isDeactivated {
		return nil, ErrEdgeworkerNoCurrentActivation
	}

	return latestActivation, nil
}

func wasDeactivationLater(activation *edgeworkers.Activation, deactivation *edgeworkers.Deactivation) (bool, error) {
	activationTime, err := time.Parse(timeLayout, activation.CreatedTime)
	if err != nil {
		return false, fmt.Errorf("failed to parse activation time")
	}
	deactivationTime, err := time.Parse(timeLayout, deactivation.CreatedTime)
	if err != nil {
		return false, fmt.Errorf("failed to parse deactivation time")
	}

	return deactivationTime.After(activationTime), nil
}

func getDeactivationsByVersionAndNetwork(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID int, version, network string) ([]edgeworkers.Deactivation, error) {
	deactivationsResp, err := client.ListDeactivations(ctx, edgeworkers.ListDeactivationsRequest{
		EdgeWorkerID: edgeworkerID,
		Version:      version,
	})
	if err != nil {
		return nil, err
	}

	return sortDeactivationsByDate(filterDeactivationsByNetwork(deactivationsResp.Deactivations, network)), nil
}

func getLatestCompletedDeactivation(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID int, version, network string, wait bool) (*edgeworkers.Deactivation, error) {
	deactivations, err := getDeactivationsByVersionAndNetwork(ctx, client, edgeworkerID, version, network)
	if err != nil {
		return nil, err
	}
	if len(deactivations) == 0 {
		return nil, ErrEdgeworkerNoLatestDeactivation
	}

	for i := range deactivations {
		d := &deactivations[i]
		if wait && statusOngoing(d.Status) {
			d, err = waitForEdgeworkerDeactivation(ctx, client, edgeworkerID, d.DeactivationID)
			if err != nil {
				return nil, err
			}
		}
		if d.Status == activationStatusComplete {
			return d, nil
		}
	}
	return nil, ErrEdgeworkerNoLatestDeactivation
}

func versionExists(version string, versions []edgeworkers.EdgeWorkerVersion) bool {
	for _, v := range versions {
		if v.Version == version {
			return true
		}
	}
	return false
}

func waitForEdgeworkerActivation(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID, activationID int) (*edgeworkers.Activation, error) {
	activation, err := client.GetActivation(ctx, edgeworkers.GetActivationRequest{
		EdgeWorkerID: edgeworkerID,
		ActivationID: activationID,
	})
	if err != nil {
		return nil, err
	}
	for activation != nil && activation.Status != activationStatusComplete {
		if !statusOngoing(activation.Status) {
			return nil, ErrEdgeworkerActivationFailure
		}
		select {
		case <-time.After(tf.MaxDuration(activationPollInterval, activationPollMinimum)):
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
				return nil, ErrEdgeworkerActivationCancelled
			}
			return nil, fmt.Errorf("%v: %w", ErrEdgeworkerActivationContextTerminated, ctx.Err())
		}
	}
	return activation, nil
}

func waitForEdgeworkerDeactivation(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID, deactivationID int) (*edgeworkers.Deactivation, error) {
	deactivation, err := client.GetDeactivation(ctx, edgeworkers.GetDeactivationRequest{
		EdgeWorkerID:   edgeworkerID,
		DeactivationID: deactivationID,
	})
	if err != nil {
		return nil, err
	}
	for deactivation != nil && deactivation.Status != activationStatusComplete {
		if !statusOngoing(deactivation.Status) {
			return nil, ErrEdgeworkerDeactivationFailure
		}
		select {
		case <-time.After(tf.MaxDuration(activationPollInterval, activationPollMinimum)):
			deactivation, err = client.GetDeactivation(ctx, edgeworkers.GetDeactivationRequest{
				EdgeWorkerID:   edgeworkerID,
				DeactivationID: deactivationID,
			})
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, ErrEdgeworkerDeactivationTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil, ErrEdgeworkerDeactivationCancelled
			}
			return nil, fmt.Errorf("%v: %w", ErrEdgeworkerDeactivationContextTerminated, ctx.Err())
		}
	}
	return deactivation, nil
}

func filterActivationsByNetwork(acts []edgeworkers.Activation, net string) (activations []edgeworkers.Activation) {
	for _, act := range acts {
		if act.Network == net {
			activations = append(activations, act)
		}
	}
	return activations
}

func filterDeactivationsByNetwork(deacts []edgeworkers.Deactivation, net string) (deactivations []edgeworkers.Deactivation) {
	for _, deact := range deacts {
		if deact.Network == edgeworkers.ActivationNetwork(net) {
			deactivations = append(deactivations, deact)
		}
	}
	return deactivations
}

func sortActivationsByDate(activations []edgeworkers.Activation) []edgeworkers.Activation {
	sort.Slice(activations, func(i, j int) bool {
		t1, err := time.Parse(timeLayout, activations[i].CreatedTime)
		if err != nil {
			panic(err)
		}
		t2, err := time.Parse(timeLayout, activations[j].CreatedTime)
		if err != nil {
			panic(err)
		}
		return t1.After(t2)
	})
	return activations
}

func sortDeactivationsByDate(deactivations []edgeworkers.Deactivation) []edgeworkers.Deactivation {
	sort.Slice(deactivations, func(i, j int) bool {
		t1, err := time.Parse(timeLayout, deactivations[i].CreatedTime)
		if err != nil {
			panic(err)
		}
		t2, err := time.Parse(timeLayout, deactivations[j].CreatedTime)
		if err != nil {
			panic(err)
		}
		return t1.After(t2)
	})
	return deactivations
}

// checkEdgeworkerExistsOnDiff is used as CustomizeDiff function
// it checks if edgeworker with provided edgeworker_id exists on ForceNew
// to avoid deactivating and then failing to activate
func checkEdgeworkerExistsOnDiff(ctx context.Context, rd *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Edgeworkers", "checkEdgeworkerExistsOnDiff")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading edgeworker activations")

	if !rd.HasChange("edgeworker_id") {
		return nil
	}

	resp, err := client.ListEdgeWorkersID(ctx, edgeworkers.ListEdgeWorkersIDRequest{})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrEdgeworkerActivation, err)
	}

	edgeworkerID, err := tf.GetIntValue("edgeworker_id", rd)
	if err != nil {
		return err
	}

	for _, e := range resp.EdgeWorkers {
		if e.EdgeWorkerID == edgeworkerID {
			return nil
		}
	}

	return fmt.Errorf("%w: edgeworker with id=%d was not found", ErrEdgeworkerActivation, edgeworkerID)
}

func suppressNoteFieldForEdgeWorkersActivation(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue != newValue && d.HasChanges("version", "network") {
		return false
	}
	return true
}

func statusOngoing(status string) bool {
	return status == activationStatusInProgress || status == activationStatusPending || status == activationStatusPresubmit
}

func statusOngoingOrReady(status string) bool {
	return status == activationStatusComplete || statusOngoing(status)
}
