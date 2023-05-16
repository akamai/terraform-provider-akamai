package edgeworkers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/sync/errgroup"
)

var (
	defaultBundleURL          = "https://raw.githubusercontent.com/akamai/edgeworkers-examples/master/edgecompute/examples/getting-started/hello-world%20(EW)/helloworld.tgz"
	defaultBundleDownloadPath = "./bundle/helloworld.tgz"
)

func resourceEdgeWorker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeWorkerCreate,
		UpdateContext: resourceEdgeWorkerUpdate,
		ReadContext:   resourceEdgeWorkerRead,
		DeleteContext: resourceEdgeWorkerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEdgeWorkerImport,
		},
		CustomizeDiff: customdiff.All(
			bundleHashCustomDiff,
		),
		Schema: map[string]*schema.Schema{
			"edgeworker_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the EdgeWorker",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The EdgeWorker name",
			},
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Defines the group association for the EdgeWorker",
				DiffSuppressFunc: tf.FieldPrefixSuppress("grp_"),
			},
			"resource_tier_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of a resource tier",
			},
			"local_bundle": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				DefaultFunc: schema.EnvDefaultFunc("EW_DEFAULT_BUNDLE_URL", defaultBundleURL),
				Description: "The path to the EdgeWorkers tgz code bundle",
			},
			"local_bundle_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The local bundle hash for the EdgeWorker",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bundle version",
			},
			"warnings": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The list of warnings returned by EdgeWorker validation",
			},
		},
	}
}

func resourceEdgeWorkerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeWorkers", "resourceEdgeWorkerCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating EdgeWorker")
	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceTierID, err := tf.GetIntValue("resource_tier_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupIDNum, err := tools.GetIntID(groupID, "grp_")
	if err != nil {
		return diag.Errorf("invalid group_id provided: %s", err)
	}
	createEdgeWorkerIDReq := edgeworkers.CreateEdgeWorkerIDRequest{
		Name:           name,
		GroupID:        groupIDNum,
		ResourceTierID: resourceTierID,
	}

	localBundlePath, err := tf.GetStringValue("local_bundle", d)
	if err != nil {
		return diag.FromErr(err)
	}
	bytesArray, err := convertLocalBundleFileIntoBytes(localBundlePath)
	if err != nil {
		return diag.FromErr(err)
	}
	validateBundleResponse, err := client.ValidateBundle(ctx, edgeworkers.ValidateBundleRequest{
		Bundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(validateBundleResponse.Errors) > 0 {
		return diag.Errorf("local bundle is not valid: %s", validateBundleResponse.Errors)
	}
	warnings, err := convertWarningsToListOfStrings(validateBundleResponse)
	if err != nil {
		return diag.Errorf("cannot marshal json %s", err)
	}
	if err = d.Set("warnings", warnings); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	edgeWorkerID, err := client.CreateEdgeWorkerID(ctx, createEdgeWorkerIDReq)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(edgeWorkerID.EdgeWorkerID))

	_, err = client.CreateEdgeWorkerVersion(ctx, edgeworkers.CreateEdgeWorkerVersionRequest{
		EdgeWorkerID:  edgeWorkerID.EdgeWorkerID,
		ContentBundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEdgeWorkerRead(ctx, d, m)
}

func resourceEdgeWorkerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeWorkers", "resourceEdgeWorkerRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Reading EdgeWorker")

	edgeWorkerID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("%s: %s", tf.ErrInvalidType, err.Error())
	}
	edgeWorker, err := client.GetEdgeWorkerID(ctx, edgeworkers.GetEdgeWorkerIDRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	versions, err := client.ListEdgeWorkerVersions(ctx, edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	var version, bundleContentHash string
	if len(versions.EdgeWorkerVersions) > 0 {
		version, err = getLatestEdgeWorkerIDBundleVersion(versions)
		if err != nil {
			return diag.Errorf("cannot indicate the latest version for edgeworker bundle")
		}
		bundleContent, err := client.GetEdgeWorkerVersionContent(ctx, edgeworkers.GetEdgeWorkerVersionContentRequest{
			EdgeWorkerID: edgeWorkerID,
			Version:      version,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		bundleContentHash, err = getSHAFromBundle(bundleContent)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	attrs := make(map[string]interface{})
	attrs["edgeworker_id"] = edgeWorkerID
	attrs["name"] = edgeWorker.Name
	attrs["group_id"] = strconv.FormatInt(edgeWorker.GroupID, 10)
	attrs["resource_tier_id"] = edgeWorker.ResourceTierID
	attrs["local_bundle_hash"] = bundleContentHash
	attrs["version"] = version
	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEdgeWorkerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeWorkers", "resourceEdgeWorkerUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Updating EdgeWorker version")
	edgeWorkerID := d.Id()
	edgeWorkerIDReq, err := strconv.Atoi(edgeWorkerID)
	if err != nil {
		return diag.Errorf("%s: %s", tf.ErrInvalidType, err.Error())
	}

	localBundlePath, err := tf.GetStringValue("local_bundle", d)
	if err != nil {
		return diag.FromErr(err)
	}
	bytesArray, err := convertLocalBundleFileIntoBytes(localBundlePath)
	if err != nil {
		return diag.FromErr(err)
	}
	bundleContentHash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)})
	if err != nil {
		return diag.FromErr(err)
	}
	hash, err := tf.GetStringValue("local_bundle_hash", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if bundleContentHash != hash {
		validateBundleResponse, err := client.ValidateBundle(ctx, edgeworkers.ValidateBundleRequest{
			Bundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(validateBundleResponse.Errors) > 0 {
			return diag.Errorf("local bundle is not valid: %s", validateBundleResponse.Errors)
		}
		warnings, err := convertWarningsToListOfStrings(validateBundleResponse)
		if err != nil {
			return diag.Errorf("cannot marshal json %s", err)
		}
		if err = d.Set("warnings", warnings); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
		_, err = client.CreateEdgeWorkerVersion(ctx, edgeworkers.CreateEdgeWorkerVersionRequest{
			EdgeWorkerID:  edgeWorkerIDReq,
			ContentBundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupIDNum, err := tools.GetIntID(groupID, "grp_")
	if err != nil {
		return diag.FromErr(err)
	}
	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceTierID, err := tf.GetIntValue("resource_tier_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	edgeWorkerIDBodyReq := edgeworkers.EdgeWorkerIDBodyRequest{
		Name:           name,
		GroupID:        groupIDNum,
		ResourceTierID: resourceTierID,
	}
	_, err = client.UpdateEdgeWorkerID(ctx, edgeworkers.UpdateEdgeWorkerIDRequest{
		EdgeWorkerIDBodyRequest: edgeWorkerIDBodyReq,
		EdgeWorkerID:            edgeWorkerIDReq,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceEdgeWorkerRead(ctx, d, m)
}

func resourceEdgeWorkerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeWorkers", "resourceEdgeWorkerDelete")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Deleting EdgeWorker and its bundle version")
	edgeWorkerIDReq, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	activations, err := checkEdgeWorkerActivations(ctx, client, edgeWorkerIDReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = deactivateEdgeWorkerVersions(ctx, client, activations); err != nil {
		return diag.FromErr(err)
	}

	versions, err := client.ListEdgeWorkerVersions(ctx, edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeWorkerIDReq,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	for _, v := range versions.EdgeWorkerVersions {
		err := client.DeleteEdgeWorkerVersion(ctx, edgeworkers.DeleteEdgeWorkerVersionRequest{
			EdgeWorkerID: edgeWorkerIDReq,
			Version:      v.Version,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = client.DeleteEdgeWorkerID(ctx, edgeworkers.DeleteEdgeWorkerIDRequest{
		EdgeWorkerID: edgeWorkerIDReq,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// checkEdgeWorkerActivations checks if there are any completed activations on staging or production networks
func checkEdgeWorkerActivations(ctx context.Context, client edgeworkers.Edgeworkers, edgeWorkerID int) ([]*edgeworkers.Activation, error) {
	var activations []*edgeworkers.Activation

	for _, network := range validEdgeworkerActivationNetworks {
		act, err := getCurrentActivation(ctx, client, edgeWorkerID, network, false)
		if err != nil {
			return nil, err
		}
		if act != nil {
			activations = append(activations, act)
		}
	}
	return activations, nil
}

// deactivateEdgeWorkerVersions loops through activations and deactivates versions in order to delete the edgeworker
func deactivateEdgeWorkerVersions(ctx context.Context, client edgeworkers.Edgeworkers, activations []*edgeworkers.Activation) error {
	g, ctxGroup := errgroup.WithContext(ctx)
	for _, act := range activations {
		act := act
		g.Go(func() error {
			return deactivateEdgeWorkerVersion(ctxGroup, client, act.EdgeWorkerID, act.Network, act.Version)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error deactivating edgeworker version: %s", err)
	}

	return nil
}

// deactivateEdgeWorkerVersion deactivates edgeworker version and waits for its completion
func deactivateEdgeWorkerVersion(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID int, network, version string) error {
	deactivation, err := client.DeactivateVersion(ctx, edgeworkers.DeactivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		DeactivateVersion: edgeworkers.DeactivateVersion{
			Network: edgeworkers.ActivationNetwork(network),
			Version: version,
		},
	})
	if err != nil {
		return err
	}
	_, err = waitForEdgeworkerDeactivation(ctx, client, edgeworkerID, deactivation.DeactivationID)
	if err != nil {
		return err
	}

	return nil
}

// since version of EdgeWorkerID bundle has type string and can be any unique value,
// this function get the latest version of EdgeWorkerID bundle according to time creation
func getLatestEdgeWorkerIDBundleVersion(versions *edgeworkers.ListEdgeWorkerVersionsResponse) (string, error) {
	var version string
	var createdTime time.Time
	for _, v := range versions.EdgeWorkerVersions {
		parsedCreatedTime, err := time.Parse(time.RFC3339, v.CreatedTime)
		if err != nil {
			return "", err
		}
		if parsedCreatedTime.After(createdTime) {
			createdTime = parsedCreatedTime
			version = v.Version
		}
	}
	return version, nil
}

func convertLocalBundleFileIntoBytes(localBundlePath string) ([]byte, error) {
	var filePath string
	if localBundlePath == defaultBundleURL {
		if err := downloadFile(defaultBundleDownloadPath, defaultBundleURL); err != nil {
			return nil, fmt.Errorf("cannot download '%s' from %s: %s", defaultBundleDownloadPath, defaultBundleURL, err.Error())
		}
		filePath = defaultBundleDownloadPath
	} else {
		filePath = localBundlePath
	}
	return ioutil.ReadFile(filePath)
}

func convertWarningsToListOfStrings(res *edgeworkers.ValidateBundleResponse) ([]string, error) {
	var warnings []string
	if len(res.Warnings) > 0 {
		for _, w := range res.Warnings {
			warning, err := json.Marshal(w)
			if err != nil {
				return nil, err
			}
			warnings = append(warnings, string(warning))
		}
	}
	return warnings, nil
}

func resourceEdgeWorkerImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeWorkers", "resourceEdgeWorkerImport")

	logger.Debug("Importing EdgeWorker version")
	client := inst.Client(meta)

	edgeWorkerID, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid edgeworker ID format: %s", err)
	}

	versions, err := client.ListEdgeWorkerVersions(ctx, edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		return nil, err
	}
	if len(versions.EdgeWorkerVersions) > 0 {
		version, err := getLatestEdgeWorkerIDBundleVersion(versions)
		if err != nil {
			return nil, fmt.Errorf("unable to determine the latest edgeworker bundle version")
		}
		bundleContent, err := client.GetEdgeWorkerVersionContent(ctx, edgeworkers.GetEdgeWorkerVersionContentRequest{
			EdgeWorkerID: edgeWorkerID,
			Version:      version,
		})
		if err != nil {
			return nil, err
		}

		validateBundleResponse, err := client.ValidateBundle(ctx, edgeworkers.ValidateBundleRequest{
			Bundle: *bundleContent,
		})
		if err != nil {
			return nil, err
		}
		warnings, err := convertWarningsToListOfStrings(validateBundleResponse)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal json %s", err)
		}
		if err = d.Set("warnings", warnings); err != nil {
			return nil, fmt.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	return []*schema.ResourceData{d}, nil
}

func bundleHashCustomDiff(_ context.Context, diff *schema.ResourceDiff, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("EdgeWorkers", "bundleHashCustomDiff")

	allSetComputed := func(fields ...string) error {
		for _, f := range fields {
			if err := diff.SetNewComputed(f); err != nil {
				return fmt.Errorf("cannot set new computed for '%s': %s", f, err)
			}
		}
		return nil
	}

	localBundleHash, err := tf.GetStringValue("local_bundle_hash", diff)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return fmt.Errorf("cannot get 'local_bundle_hash' value: %s", err)
	}
	if err != nil && errors.Is(err, tf.ErrNotFound) { // hash may be empty when resource was not created yet
		return allSetComputed("local_bundle_hash", "version", "warnings")
	}

	localBundleFileName, err := tf.GetStringValue("local_bundle", diff)
	if err != nil {
		return fmt.Errorf("cannot get 'local_bundle' value: %s", err)
	}

	f, err := openBundleFile(localBundleFileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Debugf("error closing bundle file in defer: %s", err)
		}
	}()

	hash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: f})
	if err != nil {
		return fmt.Errorf("error calculating bundle hash: %s", err)
	}

	if hash != localBundleHash {
		return allSetComputed("local_bundle_hash", "version", "warnings")
	}

	return nil
}

func openBundleFile(localBundleFileName string) (io.ReadCloser, error) {
	if localBundleFileName == defaultBundleURL {
		resp, err := http.Get(defaultBundleURL)
		if err != nil {
			return nil, fmt.Errorf("cannot dowload default bundle: %s", err)
		}
		return resp.Body, nil
	}

	f, err := os.Open(localBundleFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open bundle file (%s): %s", localBundleFileName, err)
	}
	return f, nil
}

// downloadFile downloads a file from the given URL and saves it under the given path
func downloadFile(path, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}
	defer func() {
		err = out.Close()
	}()
	return err
}
