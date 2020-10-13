package gtm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

var testAccProviders map[string]*schema.Provider
var testProvider *schema.Provider

func init() {
	testProvider = akamai.Provider(Subprovider())()
	testAccProviders = map[string]*schema.Provider{
		"akamai": testProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	TODO(t, "Check not implemented")
}

func getTestProvider() *schema.Provider {
	return testProvider
}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client gtm.GTM, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	defer func() {
		inst.client = orig
		clientLock.Unlock()
	}()

	f()
}

// TODO marks a test as being in a "pending" state and logs a message telling the user why. Such tests are expected to
// fail for the time being and may exist for the sake of unfinished/future features or to document known buggy cases
// that won't be fixed right away. The failure of a pending test is not considered an error and the test will therefore
// be skipped unless the TEST_TODO environment variable is set to a non-empty value.
func TODO(t *testing.T, message string) {
	t.Helper()
	t.Log(fmt.Sprintf("TODO: %s", message))

	if os.Getenv("TEST_TODO") == "" {
		t.Skip("TODO: Set TEST_TODO=1 in env to run this test")
	}
}

/*
///////////////////////////////

type data struct {
	data map[string]interface{}
}

func (d *data) Get(key string) interface{} {
	if value, ok := d.data[key]; ok {
		return value
	}
	return nil
}

func (d *data) GetOk(key string) (interface{}, bool) {
	if value, ok := d.data[key]; ok {
		return value, true
	}
	return nil, false
}

func (d *data) List() []interface{} {
	return []interface{}{d.data}
}

type args struct {
	schema resourceData
}
*/

func Test_getGTMV1_3Service(t *testing.T) {
	type args struct {
		schema *schema.ResourceData
	}

	tests := []struct {
		name string
		args args
		want *edgegrid.Config
		//	wantErr  error
		checkErr func(err error) bool
		edgerc   string
		env      map[string]string
	}{
		{
			name: "no valid config",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{}),
			},
			edgerc: ``,
			checkErr: func(err error) bool {
				// We do this because DeepEqual with errors or interfaces is BAD
				// Ideally the edgegrid call will return an os error for NotFound, etc.
				return err.Error() == "Unable to create instance using environment or .edgerc file"
			},
		},
		{
			name: "undefined .edgerc, undefined section",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{}),
			},
			edgerc: `[default]
		host = default
		access_token = default
		client_token = default
		client_secret = default
		max_body = 1`,
			want: &edgegrid.Config{
				Host:         "default",
				AccessToken:  "default",
				ClientToken:  "default",
				ClientSecret: "default",
				MaxBody:      1,
			},
		},
		{
			name: "undefined .edgerc, gtm default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"gtm_section": "default",
				}),
			},
			edgerc: `[default]
		host = default
		access_token = default
		client_token = default
		client_secret = default
		max_body = 1

		[not_default]
		host = not_default
		access_token = not_default
		client_token = not_default
		client_secret = not_default
		max_body = 2`,
			want: &edgegrid.Config{
				Host:         "default",
				AccessToken:  "default",
				ClientToken:  "default",
				ClientSecret: "default",
				MaxBody:      1,
			},
		},
		{
			name: "undefined .edgerc, gtm not_default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"gtm_section": "not_default",
				}),
			},
			edgerc: `[default]
		host = default
		access_token = default
		client_token = default
		client_secret = default
		max_body = 1

		[not_default]
		host = not_default
		access_token = not_default
		client_token = not_default
		client_secret = not_default
		max_body = 2`,
			want: &edgegrid.Config{
				Host:         "not_default",
				AccessToken:  "not_default",
				ClientToken:  "not_default",
				ClientSecret: "not_default",
				MaxBody:      2,
			},
		},
		{
			name: "no edgerc gtm section with env",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"gtm_section": "gtm",
				}),
			},
			env: map[string]string{
				"AKAMAI_GTM_HOST":          "env",
				"AKAMAI_GTM_ACCESS_TOKEN":  "env",
				"AKAMAI_GTM_CLIENT_TOKEN":  "env",
				"AKAMAI_GTM_CLIENT_SECRET": "env",
				"AKAMAI_GTM_MAX_BODY":      "1",
			},
			want: &edgegrid.Config{
				Host:         "env",
				AccessToken:  "env",
				ClientToken:  "env",
				ClientSecret: "env",
				MaxBody:      1,
			},
			checkErr: func(err error) bool {
				// We do this because DeepEqual with errors or interfaces is BAD
				// Ideally the edgegrid call will return an os error for NotFound, etc.
				return err.Error() == "Unable to create instance using environment or .edgerc file"
			},
		},
		{
			name: "gtm block complete",
			args: args{
				schema: func() *schema.ResourceData {
					resource := schema.Resource{
						Schema: map[string]*schema.Schema{
							"property": {
								Optional: true,
								Type:     schema.TypeSet,
								Elem:     config.Options("gtm"),
							},
						},
					}
					rd := resource.TestResourceData()
					rd.Set("gtm", schema.NewSet(func(i interface{}) int {
						return 0
					}, []interface{}{
						map[string]interface{}{
							"host":          "block",
							"access_token":  "block",
							"client_token":  "block",
							"client_secret": "block",
							"max_body":      1,
						},
					}))
					return rd
				}(),
			},
			want: &edgegrid.Config{
				Host:         "block",
				AccessToken:  "block",
				ClientToken:  "block",
				ClientSecret: "block",
				MaxBody:      1,
			},
		},
	}
	homedir.DisableCache = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempdir, err := ioutil.TempDir(os.TempDir(), "terraform-provider-akamai-test")
			if err != nil {
				t.Fatalf("unable to write .edgerc to temporary dir: %v [test: %s]", err, tt.name)
			}
			defer os.RemoveAll(tempdir)

			defer restoreEnv(os.Environ())
			setEnv(tempdir, tt.env)

			if tt.edgerc != "" {
				edgerc := path.Join(tempdir, ".edgerc")
				err = ioutil.WriteFile(
					edgerc,
					[]byte(tt.edgerc),
					0775,
				)
				if err != nil {
					t.Fatalf("unable to write .edgerc to temporary file: %s [test: %s]", edgerc, tt.name)
				}
			}

			got, err := getConfigGTMV1Service(tt.args.schema)
			if err != nil {
				if tt.checkErr != nil && tt.checkErr(err) {
					return
				}

				t.Errorf("getGTMV1Service() unexpected error = %w", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGTMV1Service() = %v, want %v", got, tt.want)
			}
		})
	}
}

type args struct {
	schema *schema.ResourceData
}

/*
	var tests = []testsStruct{
		{
			name: "no valid config",
			args: args{
				schema: schema.TestResourceDataRaw(t, inst.Provider.Schema, map[string]interface{}{}),
			},
			edgerc:  ``,
			wantErr: fmt.Errorf("Unable to create instance using environment or .edgerc file"),
		},
		{
			name: "undefined .edgerc, undefined section",
			args: args{
				schema: schema.TestResourceDataRaw(t, inst.Provider.Schema, map[string]interface{}{}),
			},
			edgerc: `[default]
host = default
access_token = default
client_token = default
client_secret = default
max_body = 1`,
			want: &edgegrid.Config{
				Host:         "default",
				AccessToken:  "default",
				ClientToken:  "default",
				ClientSecret: "default",
				MaxBody:      1,
			},
		},
		// Section specific follows
		{
			name: "undefined .edgerc, default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, inst.Provider.Schema, map[string]interface{}{
					"gtm_section": "default",
				}),
			},
			edgerc: `[default]
host = default
access_token = default
client_token = default
client_secret = default
max_body = 1

[not_default]
host = not_default
access_token = not_default
client_token = not_default
client_secret = not_default
max_body = 2`,
			want: &edgegrid.Config{
				Host:         "default",
				AccessToken:  "default",
				ClientToken:  "default",
				ClientSecret: "default",
				MaxBody:      1,
			},
		},
		{
			name: "undefined .edgerc, not_default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, inst.Provider.Schema, map[string]interface{}{
					"gtm_section": "not_default",
				}),
			},
			edgerc: `[default]
host = default
access_token = default
client_token = default
client_secret = default
max_body = 1

[not_default]
host = not_default
access_token = not_default
client_token = not_default
client_secret = not_default
max_body = 2`,
			want: &edgegrid.Config{
				Host:         "not_default",
				AccessToken:  "not_default",
				ClientToken:  "not_default",
				ClientSecret: "not_default",
				MaxBody:      2,
			},
		},
		{
			name: "no edgerc gtm section with env",
			args: args{
				schema: schema.TestResourceDataRaw(t, inst.Provider.Schema, map[string]interface{}{
					"gtm_section": "gtm",
				}),
			},
			env: map[string]string{
				"AKAMAI_GTM_HOST":          "env",
				"AKAMAI_GTM_ACCESS_TOKEN":  "env",
				"AKAMAI_GTM_CLIENT_TOKEN":  "env",
				"AKAMAI_GTM_CLIENT_SECRET": "env",
				"AKAMAI_GTM_MAX_BODY":      "1",
			},
			want: &edgegrid.Config{
				Host:         "env",
				AccessToken:  "env",
				ClientToken:  "env",
				ClientSecret: "env",
				MaxBody:      1,
			},
		},
	}

	// Invoke tests
	testGetConfigServiceExec(t, tests, getConfigGTMV1Service)

}
*/

type testsStruct struct {
	name    string
	args    args
	want    *edgegrid.Config
	wantErr error
	edgerc  string
	env     map[string]string
}

type getConfigServiceSig func(tools.ResourceDataFetcher) (*edgegrid.Config, error)

func testGetConfigServiceExec(t *testing.T, tests []testsStruct, configService getConfigServiceSig) {

	homedir.DisableCache = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempdir, err := ioutil.TempDir(os.TempDir(), "terraform-provider-akamai-test")
			if err != nil {
				t.Fatalf("unable to write .edgerc to temporary dir: %v [test: %s]", err, tt.name)
			}
			defer os.RemoveAll(tempdir)

			defer restoreEnv(os.Environ())
			setEnv(tempdir, tt.env)

			if tt.edgerc != "" {
				edgerc := path.Join(tempdir, ".edgerc")
				err = ioutil.WriteFile(
					edgerc,
					[]byte(tt.edgerc),
					0775,
				)
				if err != nil {
					t.Fatalf("unable to write .edgerc to temporary file: %s [test: %s]", edgerc, tt.name)
				}
			}

			got, err := configService(tt.args.schema)

			if err != nil {
				if reflect.DeepEqual(err, tt.wantErr) {
					return
				}
				t.Errorf("getService() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getService() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
type getConfigServiceSig func(resourceData) (*edgegrid.Config, error)

func testGetConfigServiceExec(t *testing.T, tests []testsStruct, configService getConfigServiceSig) {

	homedir.DisableCache = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempdir, err := ioutil.TempDir(os.TempDir(), "terraform-provider-akamai-test")
			if err != nil {
				t.Fatalf("unable to write .edgerc to temporary dir: %v [test: %s]", err, tt.name)
			}
			defer os.RemoveAll(tempdir)

			defer restoreEnv(os.Environ())
			setEnv(tempdir, tt.env)

			if tt.edgerc != "" {
				edgerc := path.Join(tempdir, ".edgerc")
				err = ioutil.WriteFile(
					edgerc,
					[]byte(tt.edgerc),
					0775,
				)
				if err != nil {
					t.Fatalf("unable to write .edgerc to temporary file: %s [test: %s]", edgerc, tt.name)
				}
			}

			got, err := configService(tt.args.schema)
			if err != nil {
				// TODO: Should not use deep equal with errors (RR)
				if reflect.DeepEqual(err, tt.wantErr) {
					return
				}
				t.Errorf("getService() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getService() = %v, want %v", got, tt.want)
			}
		})
	}
}

*/

func setEnv(home string, env map[string]string) {
	os.Clearenv()
	os.Setenv("HOME", home)
	if len(env) > 0 {
		for key, val := range env {
			os.Setenv(key, val)
		}
	}
}

func restoreEnv(env []string) {
	os.Clearenv()
	for _, value := range env {
		envVar := strings.Split(value, "=")
		os.Setenv(envVar[0], envVar[1])
	}
}

// loadFixtureBytes returns the entire contents of the given file as a byte slice
func loadFixtureBytes(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// loadFixtureString returns the entire contents of the given file as a string
func loadFixtureString(path string) string {
	return string(loadFixtureBytes(path))
}

// compactJSON converts a JSON-encoded byte slice to a compact form (so our JSON fixtures can be readable)
func compactJSON(encoded []byte) string {
	buf := bytes.Buffer{}
	if err := json.Compact(&buf, encoded); err != nil {
		panic(fmt.Sprintf("%s: %s", err, string(encoded)))
	}

	return buf.String()
}
