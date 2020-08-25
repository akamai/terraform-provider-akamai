package gtm

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	akamai.Provider(hclog.Default(), Subprovider())

	testAccProvider = inst.Provider
	testAccProviders = map[string]*schema.Provider{
		"akamai": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := inst.Provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {

}

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

func Test_getGTMV1_3Service(t *testing.T) {

	var tests = []testsStruct{
		testsStruct{
			name: "no valid config",
			args: args{
				schema: schema.TestResourceDataRaw(t, inst.Provider.Schema, map[string]interface{}{}),
			},
			edgerc:  ``,
			wantErr: errors.New("Unable to create instance using environment or .edgerc file"),
		},
		testsStruct{
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
		testsStruct{
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
		testsStruct{
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
		testsStruct{
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
		/*
			testsStruct{
				name: "gtm block complete",
				args: args{
					schema: &data{
						data: map[string]interface{}{
							"gtm": &data{
								data: map[string]interface{}{
									"host":          "block",
									"access_token":  "block",
									"client_token":  "block",
									"client_secret": "block",
									"max_body":      1,
								},
							},
						},
					},
				},
				want: &edgegrid.Config{
					Host:         "block",
					AccessToken:  "block",
					ClientToken:  "block",
					ClientSecret: "block",
					MaxBody:      1,
				},
			},
		*/
	}

	// Invoke tests
	testGetConfigServiceExec(t, tests, getConfigGTMV1Service)

}

type testsStruct struct {
	name    string
	args    args
	want    *edgegrid.Config
	wantErr error
	edgerc  string
	env     map[string]string
}

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
