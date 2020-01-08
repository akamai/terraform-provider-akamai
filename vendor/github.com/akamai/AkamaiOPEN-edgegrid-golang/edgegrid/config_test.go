package edgegrid

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitEdgeRc(t *testing.T) {
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

func TestInitEdgeRc_ConfigBroken(t *testing.T) {
	testSample := "../testdata/sample_edgerc"
	testConfigBroken, err := InitEdgeRc(testSample, "broken")
	assert.Equal(t, err, nil)
	assert.Equal(t, testConfigBroken.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigBroken.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigBroken.MaxBody, 128*1024)
	assert.Equal(t, testConfigBroken.HeaderToSign, []string(nil))
}

func TestInitEdgeRc_ConfigUnparsable(t *testing.T) {
	testSample := "../testdata/edgerc_that_doesnt_parse"
	_, err := InitEdgeRc(testSample, "")
	assert.Error(t, err)
}

func TestInitEdgeRc_ConfigNotFound(t *testing.T) {
	testSample := "edgerc_not_found"
	_, err := InitEdgeRc(testSample, "")
	assert.Error(t, err)
}

func TestInitEdgeRc_ConfigDashes(t *testing.T) {
	testSample := "../testdata/sample_edgerc"
	_, err := InitEdgeRc(testSample, "dashes")
	assert.Error(t, err)
}

func TestInitEdgeRc_ConfigSection(t *testing.T) {
	testConfigDefault, err := InitEdgeRc("../testdata/sample_edgerc", "test")
	assert.Equal(t, err, nil)
	assert.Equal(t, testConfigDefault.Host, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, testConfigDefault.ClientToken, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.ClientSecret, "testxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigDefault.AccessToken, "test-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigDefault.MaxBody, 131072)
	assert.Equal(t, testConfigDefault.HeaderToSign, []string(nil))
}

func TestInitEdgeRc_Broken(t *testing.T) {
	testSample := "../testdata/sample_edgerc"
	testConfigBroken, err := InitEdgeRc(testSample, "broken")
	assert.NoError(t, err)
	assert.Equal(t, testConfigBroken.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, testConfigBroken.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, testConfigBroken.MaxBody, 128*1024)
	assert.Equal(t, testConfigBroken.HeaderToSign, []string(nil))
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

func TestInitEnv_Incomplete(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("AKAMAI_HOST", "env-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.NoError(t, err)

	_, err = InitEnv("")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Fatal missing required environment variables: [AKAMAI_CLIENT_TOKEN AKAMAI_CLIENT_SECRET AKAMAI_ACCESS_TOKEN]")
}

func TestInitEnv_MaxBody(t *testing.T) {
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

func TestInit_WithEnv(t *testing.T) {
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

func TestInit_WithoutEnv(t *testing.T) {
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

func TestInit_WithSectionEnv(t *testing.T) {
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

func TestInitEdgeRc_NoDefault(t *testing.T) {
	c, err := InitEdgeRc("../testdata/nodefault_edgerc", "nodefault")
	assert.NoError(t, err)
	assert.Equal(t, c.Host, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/")
	assert.Equal(t, c.ClientToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.ClientSecret, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=")
	assert.Equal(t, c.AccessToken, "xxxx-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx")
	assert.Equal(t, c.MaxBody, 131072)
	assert.Equal(t, c.HeaderToSign, []string(nil))
}
