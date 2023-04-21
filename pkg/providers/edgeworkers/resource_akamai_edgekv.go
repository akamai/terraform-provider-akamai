package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	maxUpsertAttempts = 3
	maxInitDuration   = time.Duration(10) * time.Minute
	upsertWindow      = time.Duration(10) * time.Second
	initWindow        = time.Duration(10) * time.Second
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
				ValidateDiagFunc: tools.AggregateValidations(
					validation.ToDiagFunc(validation.IntAtLeast(0)),
					displayGroupIDWarning(),
				),
				// In the current API release, the value of group_id does not matter, so we suppress all but the first diff
				DiffSuppressFunc: func(_, old, _ string, d *schema.ResourceData) bool {
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
			"initial_data": {
				Type:        schema.TypeList,
				Optional:    true,
				Deprecated:  "The attribute 'initial_data' has been deprecated. To manage edgeKV items use 'akamai_edgekv_group_items' resource instead.",
				Description: "List of pairs to initialize the namespace. Just meaningful for creation, updates will be ignored.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"group": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "default",
						},
					},
				},
			},
		},
	}
}

func resourceEdgeKVCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Creating EdgeKV namespace configuration")

	retention64, err := tools.GetInt64Value("retention_in_seconds", tools.NewRawConfig(rd))
	if err != nil {
		return diag.FromErr(err)
	}
	retention := int(retention64)
	geoLocation, err := tools.GetStringValue("geo_location", rd)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	groupID, err := tools.GetIntValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tools.GetStringValue("namespace_name", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	// initialize edgekv
	logger.Debugf("Initializing EdgeKV...")
	initStart := time.Now()
	initStatus, err := client.InitializeEdgeKV(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	// wait for initialization
	for err == nil && initStatus.AccountStatus != "INITIALIZED" &&
		initStatus.ProductionStatus != "INITIALIZED" && initStatus.StagingStatus != "INITIALIZED" &&
		time.Since(initStart) < maxInitDuration {
		logger.Debugf("Still initializing EdgeKV...")
		time.Sleep(initWindow)
		initStatus, err = client.GetEdgeKVInitializationStatus(ctx)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if time.Since(initStart) >= maxInitDuration {
		return diag.Errorf("there was a timeout initializing the EdgeKV database: %s", time.Since(initStart).String())
	}

	// create namespace
	namespace, err := client.CreateEdgeKVNamespace(ctx, edgeworkers.CreateEdgeKVNamespaceRequest{
		Network: edgeworkers.NamespaceNetwork(network),
		Namespace: edgeworkers.Namespace{
			Name:        name,
			GeoLocation: geoLocation,
			Retention:   tools.IntPtr(retention),
			GroupID:     tools.IntPtr(groupID),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	data, err := tools.GetListValue("initial_data", rd)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	logger.Debugf("Writing initial set of data for the namespace '%s'", namespace.Name)
	if err := populateEKV(ctx, client, data, namespace, edgeworkers.ItemNetwork(network)); err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Written %d items to namespace '%s'", len(data), namespace.Name)

	rd.SetId(fmt.Sprintf("%s:%s", namespace.Name, network))

	return resourceEdgeKVRead(ctx, rd, m)
}

func resourceEdgeKVRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
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
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Updating EdgeKV namespace configuration")

	if rd.HasChanges("initial_data") {
		err := diag.Errorf("the field \"initial_data\" cannot be updated after resource creation")
		if diagnostics := diag.FromErr(tools.RestoreOldValues(rd, []string{"initial_data"})); diagnostics != nil {
			diagnostics = append(diagnostics, err[0])
			return diagnostics
		}
		return err
	}

	// at this point, just retention_in_seconds may be updated

	retention64, err := tools.GetInt64Value("retention_in_seconds", tools.NewRawConfig(rd))
	if err != nil {
		return diag.FromErr(err)
	}
	retention := int(retention64)
	// ignore group_id changes
	// changes on this field are not supported by current EdgeKV API version
	if diagnostics := diag.FromErr(tools.RestoreOldValues(rd, []string{"group_id"})); diagnostics != nil {
		return diagnostics
	}
	groupID, err := tools.GetIntValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tools.GetStringValue("namespace_name", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateEdgeKVNamespace(ctx, edgeworkers.UpdateEdgeKVNamespaceRequest{
		Network: edgeworkers.NamespaceNetwork(network),
		UpdateNamespace: edgeworkers.UpdateNamespace{
			Name:      name,
			Retention: tools.IntPtr(retention),
			GroupID:   tools.IntPtr(groupID),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEdgeKVRead(ctx, rd, m)
}

func resourceEdgeKVDelete(_ context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVDelete")
	logger.Debug("Deleting EdgeKV namespace configuration")
	logger.Info("EdgeKV namespace deletion is highly discouraged - resource will only be removed from local state")
	rd.SetId("")
	return nil
}

func populateEKV(ctx context.Context, client edgeworkers.Edgeworkers, data []interface{}, namespace *edgeworkers.Namespace, network edgeworkers.ItemNetwork) error {
	if len(data) == 0 {
		return nil
	}

	for i, rawItem := range data {
		item := rawItem.(map[string]interface{})
		upsertItemRequest := edgeworkers.UpsertItemRequest{
			ItemID:   getStringValue(item, "key"),
			ItemData: edgeworkers.Item(getStringValue(item, "value")),
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				NamespaceID: namespace.Name,
				Network:     network,
				GroupID:     getStringValue(item, "group"),
			},
		}
		for attempts := 0; attempts < maxUpsertAttempts; attempts++ {
			if i == 0 {
				time.Sleep(upsertWindow)
			}
			_, err := client.UpsertItem(ctx, upsertItemRequest)
			if err != nil {
				if strings.Contains(err.Error(), "The requested namespace does not exist or namespace type is not configured for") {
					// there might be some delay on namespace creation
					if maxUpsertAttempts > attempts+1 {
						continue
					}
				}
				return err
			}
			break
		}
	}

	return nil
}

func getStringValue(itemMap map[string]interface{}, name string) string {
	if value, ok := itemMap[name]; ok {
		return value.(string)
	}
	return ""
}

func displayGroupIDWarning() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       `Attribute "group_id" is required in order to support the next EdgeKV API release. Currently the value is not used.`,
				AttributePath: path,
			},
		}
	}
}
