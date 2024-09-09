package edgeworkers

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

var (
	bundlePathForCreate = "testdata/TestResEdgeWorkersEdgeWorker/bundles/bundleForCreate.tgz"
	bundleHashForCreate = "ba1ca447bdfebf06dee5be85eb17745b9f5dd6c718a3020409a5848f341d510f"
	bundlePathForUpdate = "testdata/TestResEdgeWorkersEdgeWorker/bundles/bundleForUpdate.tgz"
	bundleHashForUpdate = "ec177aef45a71354febdc58d0130af48c087a735e022fa53afa9b8f1e7afc245"
	defaultBundleHash   = "ba1ca447bdfebf06dee5be85eb17745b9f5dd6c718a3020409a5848f341d510f"
)

func TestResourceEdgeWorkersEdgeWorker(t *testing.T) {
	type edgeWorkerAttributes struct {
		groupID, name, localBundle, localBundleHash, version string
		resourceTierID                                       int
		timeouts                                             string
	}

	var (
		expectReadEdgeWorkerWithOneVersion = func(t *testing.T, client *edgeworkers.Mock, name, localBundlePath, version, timeForCreation string, groupID, resourceTierID, edgeWorkerID, numberOfTimes int) {
			edgeWorkerGetReq := edgeworkers.GetEdgeWorkerIDRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerGetRes := edgeworkers.EdgeWorkerID{
				EdgeWorkerID:   edgeWorkerID,
				Name:           name,
				GroupID:        int64(groupID),
				ResourceTierID: resourceTierID,
			}
			edgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
				CreatedTime:  timeForCreation,
			}
			edgeWorkerListVersionsReq := edgeworkers.ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerVersionResp := edgeworkers.ListEdgeWorkerVersionsResponse{
				EdgeWorkerVersions: []edgeworkers.EdgeWorkerVersion{
					edgeWorkerVersion,
				},
			}
			edgeWorkerVersionContentGetReq := edgeworkers.GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      version,
			}

			bytesArray, err := convertLocalBundleFileIntoBytes(localBundlePath)
			require.NoError(t, err)

			client.On("GetEdgeWorkerID", mock.Anything, edgeWorkerGetReq).Return(&edgeWorkerGetRes, nil).Times(numberOfTimes)
			client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerListVersionsReq).Return(&edgeWorkerVersionResp, nil).Times(numberOfTimes)
			for i := 0; i < numberOfTimes; i++ {
				edgeWorkerVersionContentGetRes := edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)}
				client.On("GetEdgeWorkerVersionContent", mock.Anything, edgeWorkerVersionContentGetReq).Return(&edgeWorkerVersionContentGetRes, nil).Once()
			}
		}

		expectReadEdgeWorkerWithTwoVersions = func(t *testing.T, client *edgeworkers.Mock, name, localBundlePath, version, timeForCreation, timeForUpdate string, groupID, resourceTierID, edgeWorkerID, numberOfTimes int) {
			edgeWorkerGetReq := edgeworkers.GetEdgeWorkerIDRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerGetRes := edgeworkers.EdgeWorkerID{
				EdgeWorkerID:   edgeWorkerID,
				Name:           name,
				GroupID:        int64(groupID),
				ResourceTierID: resourceTierID,
			}
			firstEdgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
				CreatedTime:  timeForCreation,
			}
			secondEdgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "2.0",
				CreatedTime:  timeForUpdate,
			}
			edgeWorkerListVersionsReq := edgeworkers.ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerVersionResp := edgeworkers.ListEdgeWorkerVersionsResponse{
				EdgeWorkerVersions: []edgeworkers.EdgeWorkerVersion{
					firstEdgeWorkerVersion,
					secondEdgeWorkerVersion,
				},
			}
			edgeWorkerVersionContentGetReq := edgeworkers.GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      version,
			}

			bytesArray, err := convertLocalBundleFileIntoBytes(localBundlePath)
			require.NoError(t, err)

			client.On("GetEdgeWorkerID", mock.Anything, edgeWorkerGetReq).Return(&edgeWorkerGetRes, nil).Times(numberOfTimes)
			client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerListVersionsReq).Return(&edgeWorkerVersionResp, nil).Times(numberOfTimes)
			for i := 0; i < numberOfTimes; i++ {
				edgeWorkerVersionContentGetRes := edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)}
				client.On("GetEdgeWorkerVersionContent", mock.Anything, edgeWorkerVersionContentGetReq).Return(&edgeWorkerVersionContentGetRes, nil).Once()
			}
		}

		expectCreateEdgeWorkerWithVersion = func(t *testing.T, client *edgeworkers.Mock, name, localBundlePath, timeForCreation string, groupID, resourceTierID, edgeWorkerID int) (*edgeworkers.EdgeWorkerID, *edgeworkers.EdgeWorkerVersion) {
			edgeWorkerReq := edgeworkers.CreateEdgeWorkerIDRequest{
				Name:           name,
				GroupID:        groupID,
				ResourceTierID: resourceTierID,
			}
			createdEdgeWorker := edgeworkers.EdgeWorkerID{
				EdgeWorkerID:   edgeWorkerID,
				ResourceTierID: resourceTierID,
				GroupID:        int64(groupID),
				Name:           name,
			}
			validateBundleRes := edgeworkers.ValidateBundleResponse{
				Warnings: []edgeworkers.ValidationIssue{
					{
						Type:    "warning_type",
						Message: "warning_message",
					},
				},
			}
			bytesArray, err := convertLocalBundleFileIntoBytes(localBundlePath)
			require.NoError(t, err)
			validateBundleReq := edgeworkers.ValidateBundleRequest{
				Bundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
			}
			bytesArray, err = convertLocalBundleFileIntoBytes(localBundlePath)
			require.NoError(t, err)
			edgeWorkerVersionReq := edgeworkers.CreateEdgeWorkerVersionRequest{
				EdgeWorkerID:  edgeWorkerID,
				ContentBundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
			}
			createdEdgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
				CreatedTime:  timeForCreation,
			}
			client.On("CreateEdgeWorkerID", mock.Anything, edgeWorkerReq).Return(&createdEdgeWorker, nil).Once()
			client.On("ValidateBundle", mock.Anything, validateBundleReq).Return(&validateBundleRes, nil).Once()
			client.On("CreateEdgeWorkerVersion", mock.Anything, edgeWorkerVersionReq).Return(&createdEdgeWorkerVersion, nil).Once()

			return &createdEdgeWorker, &createdEdgeWorkerVersion
		}

		expectUpdateEdgeWorker = func(_ *testing.T, client *edgeworkers.Mock, name, localBundlePath, timeForUpdate string, groupID, resourceTierID, edgeWorkerID int) (*edgeworkers.EdgeWorkerID, *edgeworkers.EdgeWorkerVersion) {
			updatedEdgeWorker := edgeworkers.EdgeWorkerID{
				Name:           name,
				ResourceTierID: resourceTierID,
				GroupID:        int64(groupID),
				EdgeWorkerID:   edgeWorkerID,
			}
			updateEdgeWorkerID := edgeworkers.UpdateEdgeWorkerIDRequest{
				Body: edgeworkers.EdgeWorkerIDRequestBody{
					Name:           name,
					GroupID:        groupID,
					ResourceTierID: resourceTierID,
				},
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
				CreatedTime:  timeForUpdate,
			}
			client.On("UpdateEdgeWorkerID", mock.Anything, updateEdgeWorkerID).Return(&updatedEdgeWorker, nil).Once()
			return &updatedEdgeWorker, &edgeWorkerVersion
		}

		expectUpdateEdgeWorkerVersion = func(t *testing.T, client *edgeworkers.Mock, name, localBundlePath, timeForUpdate string, groupID, resourceTierID, edgeWorkerID int) (*edgeworkers.EdgeWorkerID, *edgeworkers.EdgeWorkerVersion) {
			bytesArray, err := convertLocalBundleFileIntoBytes(localBundlePath)
			require.NoError(t, err)
			validateBundleReq := edgeworkers.ValidateBundleRequest{
				Bundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
			}
			validateBundleRes := edgeworkers.ValidateBundleResponse{
				Warnings: []edgeworkers.ValidationIssue{
					{
						Type:    "warning_type",
						Message: "warning_message",
					},
				},
			}
			updatedEdgeWorker := edgeworkers.EdgeWorkerID{
				Name:           name,
				ResourceTierID: resourceTierID,
				GroupID:        int64(groupID),
				EdgeWorkerID:   edgeWorkerID,
			}
			updateEdgeWorkerID := edgeworkers.UpdateEdgeWorkerIDRequest{
				Body: edgeworkers.EdgeWorkerIDRequestBody{
					Name:           name,
					GroupID:        groupID,
					ResourceTierID: resourceTierID,
				},
				EdgeWorkerID: edgeWorkerID,
			}
			bytesArray, err = convertLocalBundleFileIntoBytes(localBundlePath)
			require.NoError(t, err)
			edgeWorkerVersionReq := edgeworkers.CreateEdgeWorkerVersionRequest{
				EdgeWorkerID:  edgeWorkerID,
				ContentBundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
			}
			edgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "2.0",
				CreatedTime:  timeForUpdate,
			}
			client.On("ValidateBundle", mock.Anything, validateBundleReq).Return(&validateBundleRes, nil).Once()
			client.On("CreateEdgeWorkerVersion", mock.Anything, edgeWorkerVersionReq).Return(&edgeWorkerVersion, nil).Once()
			client.On("UpdateEdgeWorkerID", mock.Anything, updateEdgeWorkerID).Return(&updatedEdgeWorker, nil).Once()
			return &updatedEdgeWorker, &edgeWorkerVersion
		}

		expectDeleteEdgeWorkerWithOneVersion = func(_ *testing.T, client *edgeworkers.Mock, resourceTierID, edgeWorkerID int, timeForCreation string) {
			edgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
				CreatedTime:  timeForCreation,
			}
			edgeWorkerActivationsReq := edgeworkers.ListActivationsRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerActivationsResp := edgeworkers.ListActivationsResponse{
				Activations: []edgeworkers.Activation{},
			}
			edgeWorkerVersionsReq := edgeworkers.ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerVersionResp := edgeworkers.ListEdgeWorkerVersionsResponse{
				EdgeWorkerVersions: []edgeworkers.EdgeWorkerVersion{
					edgeWorkerVersion,
				},
			}
			edgeWorkerVersionsDeleteReq := edgeworkers.DeleteEdgeWorkerVersionRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
			}
			edgeWorkerDeleteReq := edgeworkers.DeleteEdgeWorkerIDRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			client.On("ListActivations", mock.Anything, edgeWorkerActivationsReq).Return(&edgeWorkerActivationsResp, nil).Times(2)
			client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerVersionsReq).Return(&edgeWorkerVersionResp, nil).Once()
			client.On("DeleteEdgeWorkerVersion", mock.Anything, edgeWorkerVersionsDeleteReq).Return(nil).Once()
			client.On("DeleteEdgeWorkerID", mock.Anything, edgeWorkerDeleteReq).Return(nil).Once()
		}

		expectDeleteEdgeWorkerWithTwoVersions = func(_ *testing.T, client *edgeworkers.Mock, resourceTierID, edgeWorkerID int, timeForCreation, timeForUpdate string) {
			firstEdgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
				CreatedTime:  timeForCreation,
			}
			secondEdgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      "2.0",
				CreatedTime:  timeForUpdate,
			}
			edgeWorkerActivationsReq := edgeworkers.ListActivationsRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerActivationsResp := edgeworkers.ListActivationsResponse{
				Activations: []edgeworkers.Activation{},
			}
			edgeWorkerVersionsReq := edgeworkers.ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerVersionResp := edgeworkers.ListEdgeWorkerVersionsResponse{
				EdgeWorkerVersions: []edgeworkers.EdgeWorkerVersion{
					firstEdgeWorkerVersion,
					secondEdgeWorkerVersion,
				},
			}
			edgeWorkerFirstVersionsDeleteReq := edgeworkers.DeleteEdgeWorkerVersionRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      "1.0",
			}
			edgeWorkerSecondVersionsDeleteReq := edgeworkers.DeleteEdgeWorkerVersionRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      "2.0",
			}
			edgeWorkerDeleteReq := edgeworkers.DeleteEdgeWorkerIDRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			client.On("ListActivations", mock.Anything, edgeWorkerActivationsReq).Return(&edgeWorkerActivationsResp, nil).Times(2)
			client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerVersionsReq).Return(&edgeWorkerVersionResp, nil).Once()
			client.On("DeleteEdgeWorkerVersion", mock.Anything, edgeWorkerFirstVersionsDeleteReq).Return(nil).Once()
			client.On("DeleteEdgeWorkerVersion", mock.Anything, edgeWorkerSecondVersionsDeleteReq).Return(nil).Once()
			client.On("DeleteEdgeWorkerID", mock.Anything, edgeWorkerDeleteReq).Return(nil).Once()
		}

		expectImportEdgeWorkerWithOneVersion = func(t *testing.T, client *edgeworkers.Mock, localBundlePath, version string, timeForCreation string, edgeWorkerID int) {
			edgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: edgeWorkerID,
				Version:      version,
				CreatedTime:  timeForCreation,
			}
			edgeWorkerListVersionsReq := edgeworkers.ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: edgeWorkerID,
			}
			edgeWorkerVersionResp := edgeworkers.ListEdgeWorkerVersionsResponse{
				EdgeWorkerVersions: []edgeworkers.EdgeWorkerVersion{
					edgeWorkerVersion,
				},
			}
			edgeWorkerVersionContentGetReq := edgeworkers.GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      version,
			}

			bytesArray, err := convertLocalBundleFileIntoBytes(localBundlePath)
			require.NoError(t, err)

			validateBundleReq := edgeworkers.ValidateBundleRequest{
				Bundle: edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)},
			}
			validateBundleRes := &edgeworkers.ValidateBundleResponse{
				Warnings: []edgeworkers.ValidationIssue{
					{
						Type:    "warning_type",
						Message: "warning_message",
					},
				},
			}

			client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerListVersionsReq).Return(&edgeWorkerVersionResp, nil).Once()
			edgeWorkerVersionContentGetRes := edgeworkers.Bundle{Reader: bytes.NewBuffer(bytesArray)}
			client.On("GetEdgeWorkerVersionContent", mock.Anything, edgeWorkerVersionContentGetReq).Return(&edgeWorkerVersionContentGetRes, nil).Once()
			client.On("ValidateBundle", mock.Anything, validateBundleReq).Return(validateBundleRes, nil).Once()

		}

		checkAttributes = func(attrs edgeWorkerAttributes) resource.TestCheckFunc {
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "edgeworker_id", "123"),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "version", attrs.version),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "warnings.#", "1"),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "warnings.0", "{\"type\":\"warning_type\",\"message\":\"warning_message\"}"),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "name", attrs.name),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "group_id", attrs.groupID),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "resource_tier_id", strconv.Itoa(attrs.resourceTierID)),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "local_bundle", attrs.localBundle),
				resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "local_bundle_hash", attrs.localBundleHash),
			}
			if attrs.timeouts != "" {
				checks = append(checks, resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "timeouts.#", "1"))
				checks = append(checks, resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "timeouts.0.default", attrs.timeouts))
			} else {
				checks = append(checks, resource.TestCheckResourceAttr("akamai_edgeworker.edgeworker", "timeouts.#", "0"))
			}
			return resource.ComposeAggregateTestCheckFunc(checks...)
		}
	)

	t.Run("create a new edgeworker lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 2)

		expectDeleteEdgeWorkerWithOneVersion(t, client, edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, timeForCreation)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create a new edgeworker lifecycle with timeout", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 2)

		expectDeleteEdgeWorkerWithOneVersion(t, client, edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, timeForCreation)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create_with_timeout.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
							timeouts:        "2h",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	mockBundleServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bundle, err := ioutil.ReadFile("./testdata/TestResEdgeWorkersEdgeWorker/bundles/defaultBundle.tgz")
		require.NoError(t, err)
		_, err = w.Write(bundle)
		require.NoError(t, err)
	}))
	defaultBundleDownloadPath = "./testdata/TestResEdgeWorkersEdgeWorker/bundles/target/helloworld.tgz"
	defaultBundleURL = mockBundleServer.URL
	require.NoError(t, os.Setenv("EW_DEFAULT_BUNDLE_URL", defaultBundleURL))

	t.Run("create a new edgeworker with no local bundle", func(t *testing.T) {
		defer func() {
			err := os.RemoveAll("./testdata/TestResEdgeWorkersEdgeWorker/bundles/target")
			require.NoError(t, err)
		}()
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", defaultBundleURL, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, defaultBundleURL, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 2)

		expectDeleteEdgeWorkerWithOneVersion(t, client, edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, timeForCreation)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_no_bundle.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     defaultBundleURL,
							localBundleHash: defaultBundleHash,
							version:         "1.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update edgeworker local_bundle lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)
		timeForUpdate := time.Now().Add(time.Hour * 24).Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 3)

		updatedEdgeWorker, updatedEdgeWorkerVersion := expectUpdateEdgeWorkerVersion(t, client, "example", bundlePathForUpdate, timeForUpdate, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID)
		expectReadEdgeWorkerWithTwoVersions(t, client, updatedEdgeWorker.Name, bundlePathForUpdate, updatedEdgeWorkerVersion.Version, timeForCreation, timeForUpdate, int(updatedEdgeWorker.GroupID), updatedEdgeWorker.ResourceTierID, updatedEdgeWorker.EdgeWorkerID, 2)

		expectDeleteEdgeWorkerWithTwoVersions(t, client, edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, timeForCreation, timeForUpdate)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_update_local_bundle.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForUpdate,
							localBundleHash: bundleHashForUpdate,
							version:         "2.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update edgeworker local_bundle content lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		tempBundlePath := "testdata/TestResEdgeWorkersEdgeWorker/bundles/_temp_bundle.tgz"

		timeForCreation := time.Now().Format(time.RFC3339)
		timeForUpdate := time.Now().Add(time.Hour * 24).Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 3)

		updatedEdgeWorker, updatedEdgeWorkerVersion := expectUpdateEdgeWorkerVersion(t, client, "example", bundlePathForUpdate, timeForUpdate, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID)
		expectReadEdgeWorkerWithTwoVersions(t, client, updatedEdgeWorker.Name, bundlePathForUpdate, updatedEdgeWorkerVersion.Version, timeForCreation, timeForUpdate, int(updatedEdgeWorker.GroupID), updatedEdgeWorker.ResourceTierID, updatedEdgeWorker.EdgeWorkerID, 2)

		expectDeleteEdgeWorkerWithTwoVersions(t, client, edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, timeForCreation, timeForUpdate)

		prepareTempBundleLink := func(t *testing.T, existingPath, tempPath string) func() {
			return func() {
				err := os.Remove(tempPath)
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("unable to remove temp bundle file (%s): %s", tempPath, err)
				}
				err = os.Link(existingPath, tempPath)
				if err != nil {
					t.Fatalf("unable to link temp bundle file: %s", err)
				}
			}
		}

		defer func() {
			// cleanup
			err := os.Remove(tempBundlePath)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("unable to remove temp bundle file (%s): %s", tempBundlePath, err)
			}
		}()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						PreConfig: prepareTempBundleLink(t, bundlePathForCreate, tempBundlePath),
						Config:    testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_temp_bundle.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     tempBundlePath,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
					{
						PreConfig: prepareTempBundleLink(t, bundlePathForUpdate, tempBundlePath),
						Config:    testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_temp_bundle.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     tempBundlePath,
							localBundleHash: bundleHashForUpdate,
							version:         "2.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update edgeworker group_id lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)
		timeForUpdate := time.Now().Add(time.Hour * 24).Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 3)

		updatedEdgeWorker, updatedEdgeWorkerVersion := expectUpdateEdgeWorker(t, client, "example", bundlePathForCreate, timeForUpdate, 12346, edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID)
		expectReadEdgeWorkerWithOneVersion(t, client, updatedEdgeWorker.Name, bundlePathForCreate, updatedEdgeWorkerVersion.Version, timeForUpdate, int(updatedEdgeWorker.GroupID), updatedEdgeWorker.ResourceTierID, updatedEdgeWorker.EdgeWorkerID, 2)

		expectDeleteEdgeWorkerWithOneVersion(t, client, updatedEdgeWorker.ResourceTierID, updatedEdgeWorker.EdgeWorkerID, timeForCreation)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_update_group_id.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12346",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("edgeworker no update on group_id prefix change", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 4)

		expectDeleteEdgeWorkerWithOneVersion(t, client, edgeWorker.ResourceTierID, edgeWorker.EdgeWorkerID, timeForCreation)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_update_group_id_prefix.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update edgeworker name lifecycle", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		createdTime := time.Now().Format(time.RFC3339)
		updatedTime := time.Now().Add(time.Hour * 24).Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, createdTime, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, "example", bundlePathForCreate, edgeWorkerVersion.Version, createdTime, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 3)

		updatedEdgeWorker, updatedEdgeWorkerVersion := expectUpdateEdgeWorker(t, client, "example update", bundleHashForCreate, updatedTime, 12345, edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID)
		expectReadEdgeWorkerWithOneVersion(t, client, updatedEdgeWorker.Name, bundlePathForCreate, updatedEdgeWorkerVersion.Version, updatedTime, int(updatedEdgeWorker.GroupID), updatedEdgeWorker.ResourceTierID, updatedEdgeWorker.EdgeWorkerID, 2)

		expectDeleteEdgeWorkerWithOneVersion(t, client, updatedEdgeWorker.ResourceTierID, updatedEdgeWorker.EdgeWorkerID, createdTime)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_update_name.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example update",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("delete activation upon edgeworker deletion - activations on staging and production", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 2)

		// mock two activations for edgeworker version - on staging and production
		actStaging := expectActivation(123, edgeWorker.EdgeWorkerID, "123", "", "", "", stagingNetwork, activationStatusComplete, edgeWorkerVersion.Version)
		actProduction := expectActivation(123, edgeWorker.EdgeWorkerID, "123", "", "", "", productionNetwork, activationStatusComplete, edgeWorkerVersion.Version)
		edgeWorkerActivationsReq := expectListActivationsRequest(edgeWorker.EdgeWorkerID, "")
		edgeWorkerActivationsStagingResp := expectListActivationsResponse([]edgeworkers.Activation{actStaging})
		edgeWorkerActivationsProductionResp := expectListActivationsResponse([]edgeworkers.Activation{actProduction})

		// mock two empty list deactivations calls
		edgeWorkerListDeactivationsReq := expectListDeactivationsRequest(edgeWorker.EdgeWorkerID, edgeWorkerVersion.Version)
		edgeWorkerListDeactivationsResp := expectListDeactivationsResponse([]edgeworkers.Deactivation{})

		// mock deactivate version on staging and production calls
		deactIDStaging := 321
		deactIDProduction := 322
		edgeWorkerDeactivateVersionStagingReq := expectDeactivateVersionRequest(edgeWorker.EdgeWorkerID, stagingNetwork, "", edgeWorkerVersion.Version)
		edgeWorkerDeactivateVersionProductionReq := expectDeactivateVersionRequest(edgeWorker.EdgeWorkerID, productionNetwork, "", edgeWorkerVersion.Version)
		edgeWorkerDeactivateVersionStagingResp := expectDeactivation(edgeWorker.EdgeWorkerID, deactIDStaging, edgeWorkerVersion.Version, "", activationStatusPending, stagingNetwork, "", "", "", "")
		edgeWorkerDeactivateVersionProductionResp := expectDeactivation(edgeWorker.EdgeWorkerID, deactIDProduction, edgeWorkerVersion.Version, "", activationStatusPending, productionNetwork, "", "", "", "")

		// mock get deactivation calls for staging and production
		edgeWorkerGetDeactivationStagingReq := expectGetDeactivationRequest(edgeWorker.EdgeWorkerID, deactIDStaging)
		edgeWorkerGetDeactivationProductionReq := expectGetDeactivationRequest(edgeWorker.EdgeWorkerID, deactIDProduction)
		edgeWorkerGetDeactivationStagingResp := expectDeactivation(edgeWorker.EdgeWorkerID, deactIDStaging, edgeWorkerVersion.Version, "", activationStatusComplete, stagingNetwork, "", "", "", "")
		edgeWorkerGetDeactivationProductionResp := expectDeactivation(edgeWorker.EdgeWorkerID, deactIDProduction, edgeWorkerVersion.Version, "", activationStatusComplete, productionNetwork, "", "", "", "")

		// mock list edgeworker versions call
		edgeWorkerVersionsReq := expectListEdgeWorkerVersionsRequest(edgeWorker.EdgeWorkerID)
		edgeWorkerVersionResp := expectListEdgeWorkerVersionsResponse([]edgeworkers.EdgeWorkerVersion{*edgeWorkerVersion})

		// mock edgeworker version and edgeworkerID deletions
		edgeWorkerVersionsDeleteReq := expectDeleteEdgeWorkerVersionRequest(edgeWorker.EdgeWorkerID, edgeWorkerVersion.Version)
		edgeWorkerDeleteReq := expectDeleteEdgeWorkerRequest(edgeWorker.EdgeWorkerID)

		client.On("ListActivations", mock.Anything, edgeWorkerActivationsReq).Return(&edgeWorkerActivationsStagingResp, nil).Once()
		client.On("ListActivations", mock.Anything, edgeWorkerActivationsReq).Return(&edgeWorkerActivationsProductionResp, nil).Once()
		client.On("ListDeactivations", mock.Anything, edgeWorkerListDeactivationsReq).Return(&edgeWorkerListDeactivationsResp, nil).Times(2)
		client.On("DeactivateVersion", mock.Anything, edgeWorkerDeactivateVersionProductionReq).Return(&edgeWorkerDeactivateVersionProductionResp, nil).Once()
		client.On("DeactivateVersion", mock.Anything, edgeWorkerDeactivateVersionStagingReq).Return(&edgeWorkerDeactivateVersionStagingResp, nil).Once()
		client.On("GetDeactivation", mock.Anything, edgeWorkerGetDeactivationProductionReq).Return(&edgeWorkerGetDeactivationProductionResp, nil).Once()
		client.On("GetDeactivation", mock.Anything, edgeWorkerGetDeactivationStagingReq).Return(&edgeWorkerGetDeactivationStagingResp, nil).Once()
		client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerVersionsReq).Return(&edgeWorkerVersionResp, nil).Once()
		client.On("DeleteEdgeWorkerVersion", mock.Anything, edgeWorkerVersionsDeleteReq).Return(nil).Once()
		client.On("DeleteEdgeWorkerID", mock.Anything, edgeWorkerDeleteReq).Return(nil).Once()
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("delete activation upon edgeworker deletion - only on one network - staging", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)

		timeForCreation := time.Now().Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, timeForCreation, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, edgeWorker.Name, bundlePathForCreate, edgeWorkerVersion.Version, timeForCreation, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 2)

		// mock one activation for edgeworker version - on staging
		actStaging := expectActivation(123, edgeWorker.EdgeWorkerID, "123", "", "", "", stagingNetwork, activationStatusComplete, edgeWorkerVersion.Version)
		edgeWorkerActivationsReq := expectListActivationsRequest(edgeWorker.EdgeWorkerID, "")
		edgeWorkerActivationsStagingResp := expectListActivationsResponse([]edgeworkers.Activation{actStaging})

		// mock one empty list deactivations call
		edgeWorkerListDeactivationsReq := expectListDeactivationsRequest(edgeWorker.EdgeWorkerID, edgeWorkerVersion.Version)
		edgeWorkerListDeactivationsResp := expectListDeactivationsResponse([]edgeworkers.Deactivation{})

		// mock deactivate version on staging call
		deactIDStaging := 321
		edgeWorkerDeactivateVersionStagingReq := expectDeactivateVersionRequest(edgeWorker.EdgeWorkerID, stagingNetwork, "", edgeWorkerVersion.Version)
		edgeWorkerDeactivateVersionStagingResp := expectDeactivation(edgeWorker.EdgeWorkerID, deactIDStaging, edgeWorkerVersion.Version, "", activationStatusPending, stagingNetwork, "", "", "", "")

		// mock get deactivation call for staging
		edgeWorkerGetDeactivationStagingReq := expectGetDeactivationRequest(edgeWorker.EdgeWorkerID, deactIDStaging)
		edgeWorkerGetDeactivationStagingResp := expectDeactivation(edgeWorker.EdgeWorkerID, deactIDStaging, edgeWorkerVersion.Version, "", activationStatusComplete, stagingNetwork, "", "", "", "")

		// mock list edgeworker versions call
		edgeWorkerVersionsReq := expectListEdgeWorkerVersionsRequest(edgeWorker.EdgeWorkerID)
		edgeWorkerVersionResp := expectListEdgeWorkerVersionsResponse([]edgeworkers.EdgeWorkerVersion{*edgeWorkerVersion})

		// mock edgeworker version and edgeworkerID deletions
		edgeWorkerVersionsDeleteReq := expectDeleteEdgeWorkerVersionRequest(edgeWorker.EdgeWorkerID, edgeWorkerVersion.Version)
		edgeWorkerDeleteReq := expectDeleteEdgeWorkerRequest(edgeWorker.EdgeWorkerID)

		client.On("ListActivations", mock.Anything, edgeWorkerActivationsReq).Return(&edgeWorkerActivationsStagingResp, nil).Times(1)
		client.On("ListActivations", mock.Anything, edgeWorkerActivationsReq).Return(&edgeworkers.ListActivationsResponse{}, nil).Times(1)
		client.On("ListDeactivations", mock.Anything, edgeWorkerListDeactivationsReq).Return(&edgeWorkerListDeactivationsResp, nil).Times(1)
		client.On("DeactivateVersion", mock.Anything, edgeWorkerDeactivateVersionStagingReq).Return(&edgeWorkerDeactivateVersionStagingResp, nil).Once()
		client.On("GetDeactivation", mock.Anything, edgeWorkerGetDeactivationStagingReq).Return(&edgeWorkerGetDeactivationStagingResp, nil).Once()
		client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerVersionsReq).Return(&edgeWorkerVersionResp, nil).Once()
		client.On("DeleteEdgeWorkerVersion", mock.Anything, edgeWorkerVersionsDeleteReq).Return(nil).Once()
		client.On("DeleteEdgeWorkerID", mock.Anything, edgeWorkerDeleteReq).Return(nil).Once()
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
						Check: checkAttributes(edgeWorkerAttributes{
							name:            "example",
							groupID:         "12345",
							resourceTierID:  54321,
							localBundle:     bundlePathForCreate,
							localBundleHash: bundleHashForCreate,
							version:         "1.0",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("import", func(t *testing.T) {
		testDir := "testdata/TestResEdgeWorkersEdgeWorker/edgeworker_lifecycle"
		client := new(edgeworkers.Mock)
		createdTime := time.Now().Format(time.RFC3339)

		edgeWorker, edgeWorkerVersion := expectCreateEdgeWorkerWithVersion(t, client, "example", bundlePathForCreate, createdTime, 12345, 54321, 123)
		expectReadEdgeWorkerWithOneVersion(t, client, "example", bundlePathForCreate, edgeWorkerVersion.Version, createdTime, int(edgeWorker.GroupID), edgeWorker.ResourceTierID, edgeWorkerVersion.EdgeWorkerID, 3)

		expectImportEdgeWorkerWithOneVersion(t, client, bundlePathForCreate, edgeWorkerVersion.Version, createdTime, edgeWorkerVersion.EdgeWorkerID)

		expectDeleteEdgeWorkerWithOneVersion(t, client, edgeWorker.ResourceTierID, edgeWorker.EdgeWorkerID, createdTime)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/edgeworker_create.tf", testDir)),
					},
					{
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"local_bundle"},
						ImportStateId:           "123",
						ResourceName:            "akamai_edgeworker.edgeworker",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}
func TestGetLatestEdgeWorkerIDBundleVersion(t *testing.T) {
	firstVersionTimeCreation := time.Now().Format(time.RFC3339)
	secondVersionTimeCreation := time.Now().Add(time.Hour * 24).Format(time.RFC3339)

	edgeWorkerVersions := &edgeworkers.ListEdgeWorkerVersionsResponse{
		EdgeWorkerVersions: []edgeworkers.EdgeWorkerVersion{
			{
				EdgeWorkerID: 12345,
				Version:      "0.1",
				CreatedTime:  firstVersionTimeCreation,
			},
			{
				EdgeWorkerID: 123,
				Version:      "0.2",
				CreatedTime:  secondVersionTimeCreation,
			},
		},
	}
	version, err := getLatestEdgeWorkerIDBundleVersion(edgeWorkerVersions)
	require.NoError(t, err)
	assert.Equal(t, version, "0.2")
}

func TestConvertWarningsToListOfStrings(t *testing.T) {
	tests := map[string]struct {
		res              *edgeworkers.ValidateBundleResponse
		expectedWarnings []string
	}{
		"bundle with errors": {
			res: &edgeworkers.ValidateBundleResponse{
				Errors: []edgeworkers.ValidationIssue{
					{
						Type:    "error_type",
						Message: "error_message",
					},
				},
			},
			expectedWarnings: nil,
		},
		"bundle with warnings": {
			res: &edgeworkers.ValidateBundleResponse{
				Warnings: []edgeworkers.ValidationIssue{
					{
						Type:    "warning_type",
						Message: "warning_message",
					},
				},
			},
			expectedWarnings: []string{"{\"type\":\"warning_type\",\"message\":\"warning_message\"}"},
		},
		"bundle with a few warnings": {
			res: &edgeworkers.ValidateBundleResponse{
				Warnings: []edgeworkers.ValidationIssue{
					{
						Type:    "warning_type",
						Message: "warning_message",
					},
					{
						Type:    "another_warning_type",
						Message: "another_warning_message",
					},
				},
			},
			expectedWarnings: []string{"{\"type\":\"warning_type\",\"message\":\"warning_message\"}", "{\"type\":\"another_warning_type\",\"message\":\"another_warning_message\"}"},
		},
		"bundle with errors and with warnings": {
			res: &edgeworkers.ValidateBundleResponse{
				Errors: []edgeworkers.ValidationIssue{
					{
						Type:    "error_type",
						Message: "error_message",
					},
				},
				Warnings: []edgeworkers.ValidationIssue{
					{
						Type:    "warning_type",
						Message: "warning_message",
					},
				},
			},
			expectedWarnings: []string{"{\"type\":\"warning_type\",\"message\":\"warning_message\"}"},
		},
		"bundle without errors and with warnings": {
			res:              &edgeworkers.ValidateBundleResponse{},
			expectedWarnings: nil,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			warnings, err := convertWarningsToListOfStrings(test.res)
			require.NoError(t, err)
			assert.Equal(t, warnings, test.expectedWarnings)
		})
	}
}

func TestConvertLocalBundleFileIntoBytes(t *testing.T) {
	tests := map[string]struct {
		filePath         string
		expectError      bool
		expectedErrorMsg string
	}{
		"correct path": {
			filePath:    "testdata/TestResEdgeWorkersEdgeWorker/bundles/bundleForCreate.tgz",
			expectError: false,
		},
		"incorrect path": {
			filePath:         "testdata/TestResEdgeWorkersEdgeWorker/bundles/not_existing_file.tgz",
			expectError:      true,
			expectedErrorMsg: "open testdata/TestResEdgeWorkersEdgeWorker/bundles/not_existing_file.tgz: no such file or directory",
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			if test.expectError {
				_, err := convertLocalBundleFileIntoBytes(test.filePath)
				assert.Contains(t, err.Error(), test.expectedErrorMsg)
			} else {
				bytesArray, err := convertLocalBundleFileIntoBytes(test.filePath)
				require.NoError(t, err)
				assert.True(t, len(bytesArray) > 0)
			}
		})
	}
}

var (
	expectActivation = func(actID, edgeWorkerID int, accountID, createdBy, createdTime, lastModifiedTime, network, status, version string) edgeworkers.Activation {
		return edgeworkers.Activation{
			AccountID:        accountID,
			ActivationID:     actID,
			CreatedBy:        createdBy,
			CreatedTime:      createdTime,
			EdgeWorkerID:     edgeWorkerID,
			LastModifiedTime: lastModifiedTime,
			Network:          network,
			Status:           status,
			Version:          version,
		}
	}
	expectListActivationsRequest = func(edgeWorkerID int, version string) edgeworkers.ListActivationsRequest {
		return edgeworkers.ListActivationsRequest{
			EdgeWorkerID: edgeWorkerID,
			Version:      version,
		}
	}
	expectListActivationsResponse = func(activations []edgeworkers.Activation) edgeworkers.ListActivationsResponse {
		return edgeworkers.ListActivationsResponse{
			Activations: activations,
		}
	}
	expectListDeactivationsRequest = func(edgeWorkerID int, version string) edgeworkers.ListDeactivationsRequest {
		return edgeworkers.ListDeactivationsRequest{
			EdgeWorkerID: edgeWorkerID,
			Version:      version,
		}
	}
	expectListDeactivationsResponse = func(deactivations []edgeworkers.Deactivation) edgeworkers.ListDeactivationsResponse {
		return edgeworkers.ListDeactivationsResponse{
			Deactivations: deactivations,
		}
	}
	expectDeactivateVersionRequest = func(edgeWorkerID int, network, note, version string) edgeworkers.DeactivateVersionRequest {
		return edgeworkers.DeactivateVersionRequest{
			EdgeWorkerID: edgeWorkerID,
			DeactivateVersion: edgeworkers.DeactivateVersion{
				Network: edgeworkers.ActivationNetwork(network),
				Note:    note,
				Version: version,
			},
		}
	}
	expectDeactivation = func(edgeWorkerID, deactID int, version, accountID, status, network, note, createdBy, createdTime, lastModifiedTime string) edgeworkers.Deactivation {
		return edgeworkers.Deactivation{
			EdgeWorkerID:     edgeWorkerID,
			Version:          version,
			DeactivationID:   deactID,
			AccountID:        accountID,
			Status:           status,
			Network:          edgeworkers.ActivationNetwork(network),
			Note:             note,
			CreatedBy:        createdBy,
			CreatedTime:      createdTime,
			LastModifiedTime: lastModifiedTime,
		}
	}
	expectGetDeactivationRequest = func(edgeWorkerID, deactID int) edgeworkers.GetDeactivationRequest {
		return edgeworkers.GetDeactivationRequest{
			EdgeWorkerID:   edgeWorkerID,
			DeactivationID: deactID,
		}
	}
	expectListEdgeWorkerVersionsRequest = func(edgeWorkerID int) edgeworkers.ListEdgeWorkerVersionsRequest {
		return edgeworkers.ListEdgeWorkerVersionsRequest{
			EdgeWorkerID: edgeWorkerID,
		}
	}
	expectListEdgeWorkerVersionsResponse = func(versions []edgeworkers.EdgeWorkerVersion) edgeworkers.ListEdgeWorkerVersionsResponse {
		return edgeworkers.ListEdgeWorkerVersionsResponse{
			EdgeWorkerVersions: versions,
		}
	}
	expectDeleteEdgeWorkerVersionRequest = func(edgeWorkerID int, version string) edgeworkers.DeleteEdgeWorkerVersionRequest {
		return edgeworkers.DeleteEdgeWorkerVersionRequest{
			EdgeWorkerID: edgeWorkerID,
			Version:      version,
		}
	}
	expectDeleteEdgeWorkerRequest = func(edgeWorkerID int) edgeworkers.DeleteEdgeWorkerIDRequest {
		return edgeworkers.DeleteEdgeWorkerIDRequest{
			EdgeWorkerID: edgeWorkerID,
		}
	}
)
