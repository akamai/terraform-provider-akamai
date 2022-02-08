package edgeworkers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const defaultBundle = "https://raw.githubusercontent.com/akamai/edgeworkers-examples/master/edgecompute/examples/getting-started/hello-world%20(EW)/helloworld.tgz"

func resourceEdgeWorker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeWorkerCreate,
		UpdateContext: resourceEdgeWorkerUpdate,
		ReadContext:   resourceEdgeWorkerRead,
		DeleteContext: resourceEdgeWorkerDelete,
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Defines the group association for the EdgeWorker",
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
				Default:     defaultBundle,
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
	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceTierID, err := tools.GetIntValue("resource_tier_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tools.GetIntValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	createEdgeWorkerIDReq := edgeworkers.CreateEdgeWorkerIDRequest{
		Name:           name,
		GroupID:        groupID,
		ResourceTierID: resourceTierID,
	}

	localBundlePath, err := tools.GetStringValue("local_bundle", d)
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
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
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
		return diag.Errorf("%s: %s", tools.ErrInvalidType, err.Error())
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
	attrs["group_id"] = edgeWorker.GroupID
	attrs["resource_tier_id"] = edgeWorker.ResourceTierID
	attrs["local_bundle_hash"] = bundleContentHash
	attrs["version"] = version
	if err = tools.SetAttrs(d, attrs); err != nil {
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
		return diag.Errorf("%s: %s", tools.ErrInvalidType, err.Error())
	}

	localBundlePath, err := tools.GetStringValue("local_bundle", d)
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
	hash, err := tools.GetStringValue("local_bundle_hash", d)
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
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
		_, err = client.CreateEdgeWorkerVersion(ctx, edgeworkers.CreateEdgeWorkerVersionRequest{
			EdgeWorkerID:  edgeWorkerIDReq,
			ContentBundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	groupID, err := tools.GetIntValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceTierID, err := tools.GetIntValue("resource_tier_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	edgeWorkerIDBodyReq := edgeworkers.EdgeWorkerIDBodyRequest{
		Name:           name,
		GroupID:        groupID,
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

func getSHAFromBundle(bundleContent *edgeworkers.Bundle) (string, error) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(bundleContent)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(buffer.Bytes())
	shaHash := hex.EncodeToString(h.Sum(nil))
	return shaHash, nil
}

func convertLocalBundleFileIntoBytes(localBundlePath string) ([]byte, error) {
	var filePath string
	if localBundlePath == defaultBundle {
		_, helloWorldFileName := path.Split(defaultBundle)
		helloWorldFilePath := filepath.Join("./", helloWorldFileName)
		if err := tools.DownloadFile(&tools.FileDownloader{}, helloWorldFilePath, defaultBundle); err != nil {
			return nil, fmt.Errorf("cannot download '%s' from %s: %s", helloWorldFileName, defaultBundle, err.Error())
		}
		filePath = helloWorldFilePath
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
