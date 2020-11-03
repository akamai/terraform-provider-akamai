package akamai

import (
	"context"
	"testing"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	cacheSubprovider struct{}
)

var (
	testInst *cacheSubprovider

	testAccProviders map[string]*schema.Provider

	testAccProvider *schema.Provider
)

func init() {
	testAccProvider = Provider(newCacheProvider())()
	testAccProviders = map[string]*schema.Provider{
		"akamai": testAccProvider,
	}
}

func cacheWriteTest() string {
	return `
provider "akamai" {
	cache_enabled = true
}

resource "akamai_cache" "test" {
	key = "foo"
	value = "bar"
}
`
}

func TestCache(t *testing.T) {
	t.Run("CacheSet", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: cacheWriteTest(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("foo", "id"),
					),
				},
			},
		})
	})
}

func newCacheProvider() Subprovider {
	testInst = &cacheSubprovider{}
	return testInst
}

func testDatasource() *schema.Resource {
	return &schema.Resource{
		ReadContext: testCacheGet,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func testResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: testCacheSet,
		ReadContext:   testCacheGet,
		DeleteContext: schema.NoopContext,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func testCacheSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	meta := Meta(m)
	logger := meta.Log("cache", "testCacheSet")

	logger.Debug("testing cache set")

	key := d.Get("key").(string)
	value := d.Get("value").(string)

	if err := meta.CacheSet(testInst, key, value); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})

		return diags
	}

	d.SetId(key)

	return nil
}

func testCacheGet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	meta := Meta(m)
	logger := meta.Log("cache", "testCacheGet")

	logger.Debug("testing cache get")

	var value string
	key := d.Get("key").(string)

	if err := meta.CacheGet(testInst, key, &value); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})

		return diags
	}

	d.Set("value", value)

	d.SetId(key)

	return nil
}

func (d *cacheSubprovider) Name() string {
	return "test"
}

func (d *cacheSubprovider) Version() string {
	return "0.0"
}

func (c *cacheSubprovider) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (c *cacheSubprovider) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cache": testResource(),
	}
}

func (c *cacheSubprovider) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"akamai_cache": testDatasource(),
	}
}

func (c *cacheSubprovider) Configure(log log.Interface, d *schema.ResourceData) diag.Diagnostics {
	log.Debug("START Configure")

	return nil
}
