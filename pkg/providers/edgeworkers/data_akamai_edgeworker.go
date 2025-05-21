package edgeworkers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEdgeWorker() *schema.Resource {
	return &schema.Resource{
		Description: "Get an edgeworker for given EdgeWorkerID",
		ReadContext: dataEdgeWorkerRead,
		Schema: map[string]*schema.Schema{
			"edgeworker_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the EdgeWorker",
			},
			"local_bundle": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path where the EdgeWorkers tgz code bundle will be stored",
				Default:     filepath.Join(".", "/default_name.tgz"),
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The EdgeWorker name",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the group association for the EdgeWorker",
			},
			"resource_tier_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of a resource tier",
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

func dataEdgeWorkerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeWorkers", "dataEdgeWorkerRead")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	client := inst.Client(meta)
	logger.Debug("Reading EdgeWorker")

	edgeWorkerID, err := tf.GetIntValue("edgeworker_id", d)
	if err != nil {
		return diag.Errorf("could not get value: %s", err)
	}

	edgeWorker, err := client.GetEdgeWorkerID(ctx, edgeworkers.GetEdgeWorkerIDRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		return diag.Errorf("could not get the edgeworker: %s", err)
	}

	versions, err := client.ListEdgeWorkerVersions(ctx, edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		return diag.Errorf("could not list edgeworker versions: %s", err)
	}

	var (
		warnings                   []string
		content                    []byte
		version, bundleContentHash string
		allAttrs                   bool
	)

	if len(versions.EdgeWorkerVersions) > 0 {
		version, err = getLatestEdgeWorkerIDBundleVersion(versions)
		if err != nil {
			return diag.Errorf("could not indicate the latest version for edgeworker bundle: %s", err)
		}

		bundleContent, err := client.GetEdgeWorkerVersionContent(ctx, edgeworkers.GetEdgeWorkerVersionContentRequest{
			EdgeWorkerID: edgeWorkerID,
			Version:      version,
		})
		if err != nil {
			return diag.Errorf("coult not get version content: %s", err)
		}

		content, err = io.ReadAll(bundleContent)
		if err != nil {
			return diag.Errorf("could not read content of a bundle: %s", err)
		}
		bundleContentHash, err = getSHAFromBundle(&edgeworkers.Bundle{Reader: bytes.NewReader(content)})
		if err != nil {
			return diag.Errorf("could not calculate bundle hash:  %s", err)
		}

		localBundlePath, err := tf.GetStringValue("local_bundle", d)
		if err != nil {
			return diag.Errorf("could not get local bundle: %s", err)
		}

		file, err := createFileFromPath(localBundlePath)
		if err != nil {
			return diag.Errorf("could not create a file with specified path (%s): %s", localBundlePath, err)
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				logger.Debug(fmt.Sprintf("Failed closing a file: %s", err))
			}
		}(file)

		_, err = io.Copy(file, bytes.NewReader(content))
		if err != nil {
			return diag.Errorf("could not copy bundle content to file: (%s): %s", file.Name(), err)
		}

		validateBundleResponse, err := client.ValidateBundle(ctx, edgeworkers.ValidateBundleRequest{
			Bundle: edgeworkers.Bundle{Reader: bytes.NewReader(content)},
		})
		if err != nil {
			return diag.Errorf("could not validate a bundle: %s", err)
		}

		if len(validateBundleResponse.Errors) > 0 {
			return diag.Errorf("local bundle is not valid: %s", validateBundleResponse.Errors)
		}
		warnings, err = convertWarningsToListOfStrings(validateBundleResponse)
		if err != nil {
			return diag.Errorf("cannot marshal json: %s", err)
		}
		allAttrs = true
	} else {
		allAttrs = false
	}

	attrs := createAttrs(edgeWorker, version, bundleContentHash, warnings, allAttrs)
	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.Errorf("could not set attributes: %s:", err)
	}

	d.SetId(strconv.Itoa(edgeWorkerID))
	return nil
}

func createAttrs(edgeworker *edgeworkers.EdgeWorkerID, version, bundleContentHash string, warnings []string, createAllAttrs bool) map[string]interface{} {
	if createAllAttrs {
		return map[string]interface{}{
			"name":              edgeworker.Name,
			"group_id":          strconv.FormatInt(edgeworker.GroupID, 10),
			"resource_tier_id":  edgeworker.ResourceTierID,
			"local_bundle_hash": bundleContentHash,
			"version":           version,
			"warnings":          warnings,
		}
	}
	return map[string]interface{}{
		"name":             edgeworker.Name,
		"group_id":         strconv.FormatInt(edgeworker.GroupID, 10),
		"resource_tier_id": edgeworker.ResourceTierID,
	}
}

func createFileFromPath(filePath string) (*os.File, error) {
	path := strings.Split(filePath, string(os.PathSeparator))
	if len(path) > 1 {
		dirPath := filepath.Join(path[:len(path)-1]...)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}
