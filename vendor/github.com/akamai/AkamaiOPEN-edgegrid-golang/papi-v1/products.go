package papi

import (
	"encoding/json"
	"fmt"

        edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/patrickmn/go-cache"
)

// Products represents a collection of products
type Products struct {
	client.Resource
	AccountID  string `json:"accountId"`
	ContractID string `json:"contractId"`
	Products   struct {
		Items []*Product `json:"items"`
	} `json:"products"`
}

// NewProducts creates a new Products
func NewProducts() *Products {
	products := &Products{}
	products.Init()

	return products
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (products *Products) PostUnmarshalJSON() error {
	products.Init()

	for key, product := range products.Products.Items {
		products.Products.Items[key].parent = products
		if err := product.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	return nil
}

// GetProducts populates Products with product data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listproducts
// Endpoint: GET /papi/v1/products/{?contractId}
func (products *Products) GetProducts(contract *Contract) error {
	cacheproducts, found := Profilecache.Get("products")
	if found {
		json.Unmarshal(cacheproducts.([]byte), products)
		return nil
	} else {
		req, err := client.NewRequest(
			Config,
			"GET",
			fmt.Sprintf(
				"/papi/v1/products?contractId=%s",
				contract.ContractID,
			),
			nil,
		)
		if err != nil {
			return err
		}

		edge.PrintHttpRequest(req, true)

		res, err := client.Do(Config, req)
		if err != nil {
			return err
		}

		edge.PrintHttpResponse(res, true)

		if client.IsError(res) {
			return client.NewAPIError(res)
		}

		if err = client.BodyJSON(res, products); err != nil {
			return err
		}

		byt, _ := json.Marshal(products)
		Profilecache.Set("products", byt, cache.DefaultExpiration)
		return nil
	}

}

// FindProduct finds a specific product by ID
func (products *Products) FindProduct(id string) (*Product, error) {
	var product *Product
	var productFound bool
	for _, product = range products.Products.Items {
		if product.ProductID == id {
			productFound = true
			break
		}
	}

	if !productFound {
		return nil, fmt.Errorf("Unable to find product: \"%s\"", id)
	}

	return product, nil
}

// Product represents a product resource
type Product struct {
	client.Resource
	parent      *Products
	ProductName string `json:"productName"`
	ProductID   string `json:"productId"`
}

// NewProduct creates a new Product
func NewProduct(parent *Products) *Product {
	product := &Product{parent: parent}
	product.Init()

	return product
}
