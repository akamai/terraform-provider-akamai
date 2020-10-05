package property

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

	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
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

}

// Only allow one test at a time to patch the client via useClient()
var clientLock sync.Mutex

// useClient swaps out the client on the global instance for the duration of the given func
func useClient(client papi.PAPI, f func()) {
	clientLock.Lock()
	orig := inst.client
	inst.client = client

	defer func() {
		inst.client = orig
		clientLock.Unlock()
	}()

	f()
}

func getTestProvider() *schema.Provider {
	testProvider = inst.Provider
	testProvider.Schema["edgerc"] = &schema.Schema{
		Optional:    true,
		Type:        schema.TypeString,
		DefaultFunc: schema.EnvDefaultFunc("EDGERC", nil),
	}
	testProvider.Schema["config_section"] = &schema.Schema{
		Description: "The section of the edgerc file to use for configuration",
		Optional:    true,
		Type:        schema.TypeString,
		Default:     "default",
	}
	return testProvider
}

func Test_getPAPIV1Service(t *testing.T) {
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
			name: "undefined .edgerc, property default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"property_section": "default",
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
			name: "undefined .edgerc, papi default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"papi_section": "default",
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
			name: "undefined .edgerc, property not_default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"property_section": "not_default",
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
			name: "undefined .edgerc, papi not_default section",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"papi_section": "not_default",
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
			name: "no edgerc property section with env",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"property_section": "property",
				}),
			},
			env: map[string]string{
				"AKAMAI_PROPERTY_HOST":          "env",
				"AKAMAI_PROPERTY_ACCESS_TOKEN":  "env",
				"AKAMAI_PROPERTY_CLIENT_TOKEN":  "env",
				"AKAMAI_PROPERTY_CLIENT_SECRET": "env",
				"AKAMAI_PROPERTY_MAX_BODY":      "1",
			},
			want: &edgegrid.Config{
				Host:         "env",
				AccessToken:  "env",
				ClientToken:  "env",
				ClientSecret: "env",
				MaxBody:      1,
			},
		},
		{
			name: "no edgerc papi section with env",
			args: args{
				schema: schema.TestResourceDataRaw(t, getTestProvider().Schema, map[string]interface{}{
					"papi_section": "papi",
				}),
			},
			env: map[string]string{
				"AKAMAI_PAPI_HOST":          "env",
				"AKAMAI_PAPI_ACCESS_TOKEN":  "env",
				"AKAMAI_PAPI_CLIENT_TOKEN":  "env",
				"AKAMAI_PAPI_CLIENT_SECRET": "env",
				"AKAMAI_PAPI_MAX_BODY":      "1",
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
			name: "property block complete",
			args: args{
				schema: func() *schema.ResourceData {
					resource := schema.Resource{
						Schema: map[string]*schema.Schema{
							"property": {
								Optional: true,
								Type:     schema.TypeSet,
								Elem:     config.Options("property"),
							},
						},
					}
					rd := resource.TestResourceData()
					rd.Set("property", schema.NewSet(func(i interface{}) int {
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

			got, err := getPAPIV1Service(tt.args.schema)
			if err != nil {
				if tt.checkErr != nil && tt.checkErr(err) {
					return
				}

				t.Errorf("getPAPIV1Service() unexpected error = %w", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPAPIV1Service() = %v, want %v", got, tt.want)
			}
		})
	}
}

type args struct {
	schema *schema.ResourceData
}

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