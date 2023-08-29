package edgeworkers

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	expectReadEdgeWorkersEdgeWorker = func(t *testing.T, client *edgeworkers.Mock, data testDataForEdgeWorker, timesToRun int) {
		edgeWorkerGetReq := edgeworkers.GetEdgeWorkerIDRequest{
			EdgeWorkerID: data.EdgeWorkerID,
		}

		edgeWorkerGetRes := edgeworkers.EdgeWorkerID{
			EdgeWorkerID:   data.EdgeWorkerID,
			Name:           data.Name,
			GroupID:        data.GroupID,
			ResourceTierID: data.ResourceTierID,
		}

		var edgeWorkerVersions []edgeworkers.EdgeWorkerVersion
		for i := 0; i < len(data.Version); i++ {
			edgeWorkerVersion := edgeworkers.EdgeWorkerVersion{
				EdgeWorkerID: data.EdgeWorkerID,
				Version:      data.Version[i].Version,
				CreatedTime:  data.Version[i].CreatedTime,
			}
			edgeWorkerVersions = append(edgeWorkerVersions, edgeWorkerVersion)
		}

		edgeWorkerListVersionsReq := edgeworkers.ListEdgeWorkerVersionsRequest{EdgeWorkerID: data.EdgeWorkerID}
		edgeWorkerListVersionsRes := edgeworkers.ListEdgeWorkerVersionsResponse{
			EdgeWorkerVersions: edgeWorkerVersions,
		}

		latestVersion, _ := getLatestEdgeWorkerIDBundleVersion(&edgeWorkerListVersionsRes)
		edgeWorkerGetVersionContentReq := edgeworkers.GetEdgeWorkerVersionContentRequest{
			EdgeWorkerID: data.EdgeWorkerID,
			Version:      latestVersion,
		}

		bytesArray, err := convertLocalBundleFileIntoBytes(data.ExpectedBundleFile)
		require.NoError(t, err)

		edgeWorkerValidateBundleRes := edgeworkers.ValidateBundleResponse{
			Errors:   nil,
			Warnings: data.Warnings,
		}

		client.On("GetEdgeWorkerID", mock.Anything, edgeWorkerGetReq).Return(&edgeWorkerGetRes, nil).Times(timesToRun)
		client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerListVersionsReq).Return(&edgeWorkerListVersionsRes, nil).Times(timesToRun)
		for i := 0; i < timesToRun; i++ {
			edgeWorkerGerVersionContentRes := edgeworkers.Bundle{Reader: bytes.NewReader(bytesArray)}
			client.On("GetEdgeWorkerVersionContent", mock.Anything, edgeWorkerGetVersionContentReq).Return(&edgeWorkerGerVersionContentRes, nil).Once()
		}
		client.On("ValidateBundle", mock.Anything, mock.Anything).Return(&edgeWorkerValidateBundleRes, nil).Times(timesToRun)
	}

	expectGetEdgeWorkerError = func(client *edgeworkers.Mock, errorMessage string) {
		edgeWorkerGetReq := edgeworkers.GetEdgeWorkerIDRequest{
			EdgeWorkerID: 1,
		}
		client.On("GetEdgeWorkerID", mock.Anything, edgeWorkerGetReq).Return(nil, fmt.Errorf(errorMessage)).Once()
	}

	expectListEdgeWorkerVersionsError = func(client *edgeworkers.Mock, errorMessage string) {
		edgeWorkerGetReq := edgeworkers.GetEdgeWorkerIDRequest{
			EdgeWorkerID: 1,
		}
		client.On("GetEdgeWorkerID", mock.Anything, edgeWorkerGetReq).Return(nil, nil).Once()
		client.On("ListEdgeWorkerVersions", mock.Anything, edgeworkers.ListEdgeWorkerVersionsRequest{EdgeWorkerID: 1}).Return(nil, fmt.Errorf(errorMessage)).Once()
	}

	expectReadEdgeWorkerNoVersions = func(client *edgeworkers.Mock, data testDataForEdgeWorker, timesToRun int) {
		edgeWorkerGetReq := edgeworkers.GetEdgeWorkerIDRequest{
			EdgeWorkerID: data.EdgeWorkerID,
		}

		edgeWorkerGetRes := edgeworkers.EdgeWorkerID{
			EdgeWorkerID:   data.EdgeWorkerID,
			Name:           data.Name,
			GroupID:        data.GroupID,
			ResourceTierID: data.ResourceTierID,
		}

		var edgeWorkerVersions []edgeworkers.EdgeWorkerVersion
		edgeWorkerListVersionsReq := edgeworkers.ListEdgeWorkerVersionsRequest{EdgeWorkerID: data.EdgeWorkerID}
		edgeWorkerListVersionsRes := edgeworkers.ListEdgeWorkerVersionsResponse{
			EdgeWorkerVersions: edgeWorkerVersions,
		}

		client.On("GetEdgeWorkerID", mock.Anything, edgeWorkerGetReq).Return(&edgeWorkerGetRes, nil).Times(timesToRun)
		client.On("ListEdgeWorkerVersions", mock.Anything, edgeWorkerListVersionsReq).Return(&edgeWorkerListVersionsRes, nil).Times(timesToRun)
	}

	oneVersionData = testDataForEdgeWorker{
		EdgeWorkerID:       1,
		GroupID:            11,
		ResourceTierID:     1000,
		Name:               "Test Name",
		LocalBundlePath:    "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_one_version.tgz",
		LocalBundleHash:    "ba1ca447bdfebf06dee5be85eb17745b9f5dd6c718a3020409a5848f341d510f",
		ExpectedBundleFile: "testdata/TestDataEdgeWorkersEdgeWorker/defaultBundle.tgz",
		Version: []edgeworkers.EdgeWorkerVersion{
			{
				EdgeWorkerID: 1,
				Version:      "1.0",
				CreatedTime:  "2006-01-02T15:04:05Z",
			},
		},
		Warnings: nil,
	}

	twoVersionsData = testDataForEdgeWorker{
		EdgeWorkerID:       2,
		GroupID:            22,
		ResourceTierID:     300,
		Name:               "Test Name2",
		LocalBundlePath:    "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_two_versions.tgz",
		LocalBundleHash:    "ba1ca447bdfebf06dee5be85eb17745b9f5dd6c718a3020409a5848f341d510f",
		ExpectedBundleFile: "testdata/TestDataEdgeWorkersEdgeWorker/defaultBundle.tgz",
		Version: []edgeworkers.EdgeWorkerVersion{
			{
				EdgeWorkerID: 2,
				Version:      "2.0",
				CreatedTime:  time.Now().Format(time.RFC3339),
			},
			{
				EdgeWorkerID: 2,
				Version:      "3.0",
				CreatedTime:  time.Now().Format(time.RFC3339),
			},
		},
		Warnings: nil,
	}

	oneWarningData = testDataForEdgeWorker{
		EdgeWorkerID:       3,
		GroupID:            33,
		ResourceTierID:     200,
		Name:               "Test Name",
		LocalBundlePath:    "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_one_warning.tgz",
		LocalBundleHash:    "ba1ca447bdfebf06dee5be85eb17745b9f5dd6c718a3020409a5848f341d510f",
		ExpectedBundleFile: "testdata/TestDataEdgeWorkersEdgeWorker/defaultBundle.tgz",
		Version: []edgeworkers.EdgeWorkerVersion{
			{
				EdgeWorkerID: 3,
				Version:      "3.0",
				CreatedTime:  time.Now().Format(time.RFC3339),
			},
		},
		Warnings: []edgeworkers.ValidationIssue{
			{
				Type:    "warning",
				Message: "warning one",
			},
		},
	}

	threeWarningsData = testDataForEdgeWorker{
		EdgeWorkerID:       4,
		GroupID:            44,
		ResourceTierID:     100,
		Name:               "Test Name",
		LocalBundlePath:    "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_three_warnings.tgz",
		LocalBundleHash:    "ba1ca447bdfebf06dee5be85eb17745b9f5dd6c718a3020409a5848f341d510f",
		ExpectedBundleFile: "testdata/TestDataEdgeWorkersEdgeWorker/defaultBundle.tgz",
		Version: []edgeworkers.EdgeWorkerVersion{
			{
				EdgeWorkerID: 4,
				Version:      "3.0",
				CreatedTime:  time.Now().Format(time.RFC3339),
			},
		},
		Warnings: []edgeworkers.ValidationIssue{
			{
				Type:    "warning",
				Message: "warning one",
			},
			{
				Type:    "warning",
				Message: "warning two",
			},
			{
				Type:    "warning",
				Message: "warning three",
			},
		},
	}

	noVersionsData = testDataForEdgeWorker{
		EdgeWorkerID:       5,
		GroupID:            55,
		ResourceTierID:     500,
		Name:               "Test Name",
		LocalBundlePath:    "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/no_versions.tgz",
		LocalBundleHash:    "",
		ExpectedBundleFile: "",
		Version:            []edgeworkers.EdgeWorkerVersion{},
		Warnings:           nil,
	}

	defaultBundlePathData = testDataForEdgeWorker{
		EdgeWorkerID:       1,
		GroupID:            11,
		ResourceTierID:     1,
		Name:               "Test Name",
		LocalBundlePath:    "default_name.tgz",
		LocalBundleHash:    "ba1ca447bdfebf06dee5be85eb17745b9f5dd6c718a3020409a5848f341d510f",
		ExpectedBundleFile: "testdata/TestDataEdgeWorkersEdgeWorker/defaultBundle.tgz",
		Version: []edgeworkers.EdgeWorkerVersion{
			{
				EdgeWorkerID: 1,
				Version:      "1.0",
				CreatedTime:  "2006-01-02T15:04:05Z",
			},
		},
		Warnings: nil,
	}
)

type testDataForEdgeWorker struct {
	EdgeWorkerID       int
	GroupID            int64
	ResourceTierID     int
	Name               string
	LocalBundlePath    string
	LocalBundleHash    string
	ExpectedBundleFile string
	Version            []edgeworkers.EdgeWorkerVersion
	Warnings           []edgeworkers.ValidationIssue
}

func TestDataEdgeWorkersEdgeWorker(t *testing.T) {
	tests := map[string]struct {
		init       func(*testing.T, *edgeworkers.Mock, testDataForEdgeWorker)
		mockData   testDataForEdgeWorker
		configPath string
		error      *regexp.Regexp
	}{
		"happy path with one version": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectReadEdgeWorkersEdgeWorker(t, m, testData, 5)
			},
			mockData:   oneVersionData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_one_version.tf",
			error:      nil,
		},
		"happy path with 2 versions": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectReadEdgeWorkersEdgeWorker(t, m, testData, 5)
			},
			mockData:   twoVersionsData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_two_versions.tf",
			error:      nil,
		},
		"happy path with one warning": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectReadEdgeWorkersEdgeWorker(t, m, testData, 5)
			},
			mockData:   oneWarningData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_one_warning.tf",
			error:      nil,
		},
		"happy path with three warnings": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectReadEdgeWorkersEdgeWorker(t, m, testData, 5)
			},
			mockData:   threeWarningsData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_three_warnings.tf",
			error:      nil,
		},
		"happy path without local bundle path specified": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectReadEdgeWorkersEdgeWorker(t, m, testData, 5)
			},
			mockData:   defaultBundlePathData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_no_local_bundle.tf",
			error:      nil,
		},
		"no versions": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectReadEdgeWorkerNoVersions(m, noVersionsData, 5)
			},
			mockData:   noVersionsData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_no_versions.tf",
			error:      nil,
		},
		"could not get an edgeworker_id": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectGetEdgeWorkerError(m, "could not get an edgeworker")
			},
			mockData:   oneVersionData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_one_version.tf",
			error:      regexp.MustCompile("could not get an edgeworker"),
		},
		"could not list versions": {
			init: func(t *testing.T, m *edgeworkers.Mock, testData testDataForEdgeWorker) {
				expectListEdgeWorkerVersionsError(m, "could not list edgeworker versions")
			},
			mockData:   oneVersionData,
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_one_version.tf",
			error:      regexp.MustCompile("could not list edgeworker versions"),
		},
		"edgeworker_id not provided": {
			init:       func(t *testing.T, m *edgeworkers.Mock, worker testDataForEdgeWorker) {},
			mockData:   testDataForEdgeWorker{},
			configPath: "testdata/TestDataEdgeWorkersEdgeWorker/edgeworker_no_edgeworker_id.tf",
			error:      regexp.MustCompile("Missing required argument"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(t, client, test.mockData)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkAttrsForEdgeWorker(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
	if _, err := os.Stat("default_name.tgz"); err == nil {
		err = os.Remove("default_name.tgz")
		if err != nil {
			t.Fatalf("unable to remove temp bundle file (%s): %s", "default_name.tgz", err)
		}
	}

}

func checkAttrsForEdgeWorker(data testDataForEdgeWorker) resource.TestCheckFunc {
	if len(data.Version) == 0 {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "edgeworker_id", strconv.Itoa(data.EdgeWorkerID)),
			resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "local_bundle", data.LocalBundlePath),
			resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "name", data.Name),
			resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "group_id", strconv.FormatInt(data.GroupID, 10)),
			resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "resource_tier_id", strconv.Itoa(data.ResourceTierID)),
			resource.TestCheckNoResourceAttr("data.akamai_edgeworker.test", "local_bundle_hash"),
			resource.TestCheckNoResourceAttr("data.akamai_edgeworker.test", "version"),
			resource.TestCheckNoResourceAttr("data.akamai_edgeworker.test", "warnings"))
	}
	latestVersion, _ := getLatestEdgeWorkerIDBundleVersion(&edgeworkers.ListEdgeWorkerVersionsResponse{EdgeWorkerVersions: data.Version})
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "edgeworker_id", strconv.Itoa(data.EdgeWorkerID)),
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "local_bundle", data.LocalBundlePath),
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "name", data.Name),
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "group_id", strconv.FormatInt(data.GroupID, 10)),
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "resource_tier_id", strconv.Itoa(data.ResourceTierID)),
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "local_bundle_hash", data.LocalBundleHash),
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "version", latestVersion),
		resource.TestCheckResourceAttr("data.akamai_edgeworker.test", "warnings.#", strconv.Itoa(len(data.Warnings))),
		checkWarningsAttrForEdgeWorker(data.Warnings),
	)
}

func checkWarningsAttrForEdgeWorker(warnings []edgeworkers.ValidationIssue) resource.TestCheckFunc {
	var warningsCheckFuncs []resource.TestCheckFunc
	for i := 0; i < len(warnings); i++ {
		warningsCheckFuncs = append(warningsCheckFuncs, resource.TestCheckResourceAttr("data.akamai_edgeworker.test", fmt.Sprintf("warnings.%d", i), createWarningEntry(warnings[i].Type, warnings[i].Message)))
	}
	return resource.ComposeAggregateTestCheckFunc(warningsCheckFuncs...)
}

func createWarningEntry(warningType, warningMessage string) string {
	return fmt.Sprintf("{\"type\":\"%s\",\"message\":\"%s\"}", warningType, warningMessage)
}
