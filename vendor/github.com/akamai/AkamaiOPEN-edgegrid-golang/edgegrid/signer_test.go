package edgegrid

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/stretchr/testify/assert"
)

var (
	testFile  = "../testdata/testdata.json"
	timestamp = "20140321T19:34:21+0000"
	nonce     = "nonce-xx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	config    = Config{
		Host:         "https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
		Debug:        true,
		HeaderToSign: []string{
			"X-Test1",
			"X-Test2",
			"X-Test3",
		},
	}
)

type JSONTests struct {
	Tests []Test `json:"tests"`
}
type Test struct {
	Name    string `json:"testName"`
	Request struct {
		Method  string              `json:"method"`
		Path    string              `json:"path"`
		Headers []map[string]string `json:"headers"`
		Data    string              `json:"data"`
	} `json:"request"`
	ExpectedAuthorization string `json:"expectedAuthorization"`
}

func TestMakeEdgeTimeStamp(t *testing.T) {
	actual := makeEdgeTimeStamp()
	expected := regexp.MustCompile(`^\d{4}[0-1][0-9][0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]\+0000$`)
	if assert.Regexp(t, expected, actual, "Fail: Regex do not match") {
		t.Log("Pass: Regex matches")
	}
}

func TestCreateNonce(t *testing.T) {
	actual := createNonce()
	for i := 0; i < 100; i++ {
		expected := createNonce()
		assert.NotEqual(t, actual, expected, "Fail: Nonce matches")
	}
}

func TestCreateAuthHeader(t *testing.T) {
	var edgegrid JSONTests
	byt, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Test file not found, err %s", err)
	}
	url, err := url.Parse(config.Host)
	if err != nil {
		t.Fatalf("URL is not parsable, err %s", err)
	}
	err = jsonhooks.Unmarshal(byt, &edgegrid)
	if err != nil {
		t.Fatalf("JSON is not parsable, err %s", err)
	}
	for _, edge := range edgegrid.Tests {
		url.Path = edge.Request.Path
		req, _ := http.NewRequest(
			edge.Request.Method,
			url.String(),
			bytes.NewBuffer([]byte(edge.Request.Data)),
		)
		for _, header := range edge.Request.Headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		actual := createAuthHeader(config, req, timestamp, nonce)
		if assert.Equal(t, edge.ExpectedAuthorization, actual, fmt.Sprintf("Fail: %s", edge.Name)) {
			t.Logf("Pass: %s\n", edge.Name)
			t.Logf("Expected: %s - Actual %s", edge.ExpectedAuthorization, actual)
		}

	}
}

func TestInitConfigBroken(t *testing.T) {
	testSample := "../testdata/sample_edgerc"
	testConfigBroken, err := InitEdgeRc(testSample, "broken")
	assert.Equal(t, err, nil)
	assert.Equal(t, testConfigBroken.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigBroken.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigBroken.MaxBody, 128*1024)
	assert.Equal(t, testConfigBroken.HeaderToSign, []string(nil))
}

func TestInitConfigUnparsable(t *testing.T) {
	testSample := "../testdata/edgerc_that_doesnt_parse"
	_, err := InitEdgeRc(testSample, "")
	assert.Error(t, err)
}

func TestInitConfigNotFound(t *testing.T) {
	testSample := "edgerc_not_found"
	_, err := InitEdgeRc(testSample, "")
	assert.Error(t, err)
}

func TestInitConfigDashes(t *testing.T) {
	testSample := "../testdata/sample_edgerc"
	_, err := InitEdgeRc(testSample, "dashes")
	assert.Error(t, err)
}

func TestInitConfigDefault(t *testing.T) {
	var configDefault = []string{
		"",
		"default",
	}
	for _, section := range configDefault {
		testConfigDefault, err := InitEdgeRc("../testdata/sample_edgerc", section)
		assert.Equal(t, err, nil)
		assert.Equal(t, testConfigDefault.ClientToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
		assert.Equal(t, testConfigDefault.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
		assert.Equal(t, testConfigDefault.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
		assert.Equal(t, testConfigDefault.MaxBody, 131072)
		assert.Equal(t, testConfigDefault.HeaderToSign, []string(nil))
	}
}

func TestInitConfigSection(t *testing.T) {
	testConfigDefault, err := InitEdgeRc("../testdata/sample_edgerc", "test")
	assert.Equal(t, err, nil)
	assert.Equal(t, testConfigDefault.Host, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, testConfigDefault.ClientToken, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.ClientSecret, "testxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigDefault.AccessToken, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.MaxBody, 131072)
	assert.Equal(t, testConfigDefault.HeaderToSign, []string(nil))
}

func TestInitEdgeRcBroken(t *testing.T) {
	testSample := "../testdata/sample_edgerc"
	testConfigBroken, err := InitEdgeRc(testSample, "broken")
	assert.NoError(t, err)
	assert.Equal(t, testConfigBroken.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigBroken.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigBroken.MaxBody, 128*1024)
	assert.Equal(t, testConfigBroken.HeaderToSign, []string(nil))
}

func TestInitEdgeRcDefault(t *testing.T) {
	var configDefault = []string{
		"",
		"default",
	}
	for _, section := range configDefault {
		testConfigDefault, err := InitEdgeRc("../testdata/sample_edgerc", section)
		assert.NoError(t, err)
		assert.Equal(t, testConfigDefault.Host, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
		assert.Equal(t, testConfigDefault.ClientToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
		assert.Equal(t, testConfigDefault.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
		assert.Equal(t, testConfigDefault.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
		assert.Equal(t, testConfigDefault.MaxBody, 131072)
		assert.Equal(t, testConfigDefault.HeaderToSign, []string(nil))
	}
}

func TestInitEdgeRcSection(t *testing.T) {
	testConfigDefault, err := InitEdgeRc("../testdata/sample_edgerc", "test")
	assert.NoError(t, err)
	assert.Equal(t, testConfigDefault.Host, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, testConfigDefault.ClientToken, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.ClientSecret, "testxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigDefault.AccessToken, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.MaxBody, 131072)
	assert.Equal(t, testConfigDefault.HeaderToSign, []string(nil))
}

func TestInitEnv(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("AKAMAI_HOST", "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.NoError(t, err)

	err = os.Setenv("AKAMAI_CLIENT_TOKEN", "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_CLIENT_SECRET", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_ACCESS_TOKEN", "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)

	c, err := InitEnv("")
	assert.NoError(t, err)
	assert.Equal(t, c.Host, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, c.ClientToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, c.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.MaxBody, 131072)
	assert.Equal(t, c.HeaderToSign, []string(nil))
}

func TestInitEnvIncomplete(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("AKAMAI_HOST", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.NoError(t, err)

	_, err = InitEnv("")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Fatal missing required environment variables: [AKAMAI_CLIENT_TOKEN AKAMAI_CLIENT_SECRET AKAMAI_ACCESS_TOKEN]")
}

func TestInitEnvMaxBody(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("AKAMAI_HOST", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_CLIENT_TOKEN", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_CLIENT_SECRET", "envxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_ACCESS_TOKEN", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_MAX_BODY", "42")
	assert.NoError(t, err)

	c, err := InitEnv("")
	assert.NoError(t, err)
	assert.Equal(t, c.Host, "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, c.ClientToken, "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.ClientSecret, "envxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, c.AccessToken, "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.MaxBody, 42)
	assert.Equal(t, c.HeaderToSign, []string(nil))
}

func TestInitWithEnv(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("AKAMAI_HOST", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_CLIENT_TOKEN", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_CLIENT_SECRET", "envxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_ACCESS_TOKEN", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)

	c, err := InitEnv("")
	assert.NoError(t, err)
	assert.Equal(t, c.Host, "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, c.ClientToken, "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.ClientSecret, "envxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, c.AccessToken, "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.MaxBody, 131072)
	assert.Equal(t, c.HeaderToSign, []string(nil))
}

func TestInitWithoutEnv(t *testing.T) {
	os.Clearenv()

	c, err := InitEnv("")
	assert.Error(t, err)
	assert.NotEqual(t, c.Host, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.NotEqual(t, c.ClientToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NotEqual(t, c.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.NotEqual(t, c.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NotEqual(t, c.MaxBody, 131072)
	assert.Equal(t, c.HeaderToSign, []string(nil))
}

func TestInitWithSectionEnv(t *testing.T) {
	os.Clearenv()

	err := os.Setenv("AKAMAI_TEST_HOST", "testenv-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_TEST_CLIENT_TOKEN", "testenv-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_TEST_CLIENT_SECRET", "testenvxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.NoError(t, err)
	err = os.Setenv("AKAMAI_TEST_ACCESS_TOKEN", "testenv-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.NoError(t, err)

	c, err := InitEnv("test")
	assert.NoError(t, err)
	assert.Equal(t, c.Host, "testenv-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, c.ClientToken, "testenv-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.ClientSecret, "testenvxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, c.AccessToken, "testenv-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.MaxBody, 131072)
	assert.Equal(t, c.HeaderToSign, []string(nil))
}
