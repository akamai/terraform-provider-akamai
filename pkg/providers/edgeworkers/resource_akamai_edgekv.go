package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	initWindow = time.Duration(10) * time.Second
)

func resourceEdgeKV() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeKVCreate,
		ReadContext:   resourceEdgeKVRead,
		UpdateContext: resourceEdgeKVUpdate,
		DeleteContext: resourceEdgeKVDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name for the EKV namespace",
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(edgeworkers.NamespaceStagingNetwork), string(edgeworkers.NamespaceProductionNetwork),
				}, false)),
				Description: "The network on which the namespace will be activated",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Namespace ACC group ID. It will be used in EdgeKV API v2. Not updatable.",
				ValidateDiagFunc: tf.AggregateValidations(
					validation.ToDiagFunc(validation.IntAtLeast(0)),
					displayGroupIDWarning(),
				),
				// In the current API release, the value of group_id does not matter, so we suppress all but the first diff
				DiffSuppressFunc: func(_, old, _ string, _ *schema.ResourceData) bool {
					return old != ""
				},
				ForceNew: true,
			},
			"retention_in_seconds": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.Any(
					validation.IntInSlice([]int{0}),
					validation.All(validation.IntAtLeast(86400), validation.IntAtMost(315360000)),
				)),
				Description: "Retention period for data in this namespace. An update of this value will just affect new EKV items.",
			},
			"geo_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage location for data",
				ForceNew:    true,
			},
		},
	}
}

func waitForEdgeKVInitialization(ctx context.Context, client edgeworkers.Edgeworkers) error {
	status := &edgeworkers.EdgeKVInitializationStatus{}
	var err error

	for status.AccountStatus != "INITIALIZED" {
		select {
		case <-time.After(initWindow):
			status, err = client.GetEdgeKVInitializationStatus(ctx)
			if err != nil {
				return fmt.Errorf("could not get EdgeKV initialization status: %s", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("retry timeout reached: incorrect status of edgeKV: %s, %s", status.AccountStatus, ctx.Err())
		}
	}

	return nil
}

func resourceEdgeKVCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Creating EdgeKV namespace configuration")

	retention64, err := tf.GetInt64Value("retention_in_seconds", tf.NewRawConfig(rd))
	if err != nil {
		return diag.FromErr(err)
	}
	retention := int(retention64)
	geoLocation, err := tf.GetStringValue("geo_location", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	groupID, err := tf.GetIntValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tf.GetStringValue("namespace_name", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	status, err := client.GetEdgeKVInitializationStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	// If the status is "UNINITIALIZED", we have to send initialization request and wait for "INITIALIZED" status. If the
	// status is "PENDING" we have to wait. If the status is "INITIALIZED" we can proceed.
	if status.AccountStatus == "UNINITIALIZED" {
		// initialize edgekv
		logger.Debugf("Initializing EdgeKV...")
		_, err = client.InitializeEdgeKV(ctx)
		if err != nil {
			return diag.Errorf("could not initialize edgeKV: %s", err)
		}
		if err = waitForEdgeKVInitialization(ctx, client); err != nil {
			return diag.FromErr(err)
		}
	} else if status.AccountStatus == "PENDING" {
		if err = waitForEdgeKVInitialization(ctx, client); err != nil {
			return diag.FromErr(err)
		}
	}

	// create namespace
	namespace, err := client.CreateEdgeKVNamespace(ctx, edgeworkers.CreateEdgeKVNamespaceRequest{
		Network: edgeworkers.NamespaceNetwork(network),
		Namespace: edgeworkers.Namespace{
			Name:        name,
			GeoLocation: geoLocation,
			Retention:   ptr.To(retention),
			GroupID:     ptr.To(groupID),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(fmt.Sprintf("%s:%s", namespace.Name, network))

	return resourceEdgeKVRead(ctx, rd, m)
}

func resourceEdgeKVRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Reading EdgeKV namespace configuration")

	id := strings.Split(rd.Id(), ":")
	if len(id) < 2 {
		return diag.Errorf("invalid EdgeKV identifier: %s", rd.Id())
	}
	name := id[0]
	network := id[1]

	namespace, err := client.GetEdgeKVNamespace(ctx, edgeworkers.GetEdgeKVNamespaceRequest{
		Network: edgeworkers.NamespaceNetwork(network),
		Name:    name,
	})
	if err != nil {
		logger.Errorf("EdgeKV namespace '%s' not found in network '%s': %s", name, network, err.Error())
		return diag.FromErr(err)
	}

	if err := rd.Set("geo_location", namespace.GeoLocation); err != nil {
		return diag.FromErr(err)
	}

	if err := rd.Set("group_id", namespace.GroupID); err != nil {
		return diag.FromErr(err)
	}

	if err := rd.Set("namespace_name", namespace.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := rd.Set("network", network); err != nil {
		return diag.FromErr(err)
	}

	if err := rd.Set("retention_in_seconds", namespace.Retention); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEdgeKVUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Updating EdgeKV namespace configuration")

	// at this point, just retention_in_seconds may be updated
	retention64, err := tf.GetInt64Value("retention_in_seconds", tf.NewRawConfig(rd))
	if err != nil {
		return diag.FromErr(err)
	}
	retention := int(retention64)
	// ignore group_id changes, as changes on this field are not supported by current EdgeKV API version
	if diagnostics := diag.FromErr(tf.RestoreOldValues(rd, []string{"group_id"})); diagnostics != nil {
		return diagnostics
	}
	groupID, err := tf.GetIntValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tf.GetStringValue("namespace_name", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateEdgeKVNamespace(ctx, edgeworkers.UpdateEdgeKVNamespaceRequest{
		Network: edgeworkers.NamespaceNetwork(network),
		UpdateNamespace: edgeworkers.UpdateNamespace{
			Name:      name,
			Retention: ptr.To(retention),
			GroupID:   ptr.To(groupID),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEdgeKVRead(ctx, rd, m)
}

func resourceEdgeKVDelete(_ context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVDelete")
	logger.Debug("Deleting EdgeKV namespace configuration")
	logger.Info("EdgeKV namespace deletion is highly discouraged - resource will only be removed from local state")
	rd.SetId("")
	return nil
}

func displayGroupIDWarning() schema.SchemaValidateDiagFunc {
	return func(_ interface{}, path cty.Path) diag.Diagnostics {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       `Attribute "group_id" is required in order to support the next EdgeKV API release. Currently the value is not used.`,
				AttributePath: path,
			},
		}
	}
}
