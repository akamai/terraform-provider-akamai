package akamai

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourcePropertyRules() *schema.Resource {
	return &schema.Resource{
		Create: resourcePropertyRulesCreate,
		Read:   resourcePropertyRulesRead,
		Update: resourcePropertyRulesUpdate,
		Delete: resourcePropertyRulesDelete,
		Exists: resourcePropertyRulesExists,
		Importer: &schema.ResourceImporter{
			State: resourcePropertyImport,
		},
		Schema: akamaiPropertyRulesSchema,
	}
}

var akpsRulesOption = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var akpsRulesCriteria = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"option": akpsRulesOption,
		},
	},
}

var akpsRulesBehavior = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"option": akpsRulesOption,
		},
	},
}

var akamaiPropertyRulesSchema = map[string]*schema.Schema{
	// rules tree can go max 5 levels deep
	"variables": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"rules": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"criteria_match": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "all",
				},
				"behavior": akpsBehavior,
				"rule": &schema.Schema{
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"comment": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"criteria_match": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "all",
							},
							"criteria": akpsCriteria,
							"behavior": akpsBehavior,
							"rule": &schema.Schema{
								Type:     schema.TypeSet,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"comment": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"criteria_match": {
											Type:     schema.TypeString,
											Optional: true,
											Default:  "all",
										},
										"criteria": akpsCriteria,
										"behavior": akpsBehavior,
										"rule": &schema.Schema{
											Type:     schema.TypeSet,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"name": {
														Type:     schema.TypeString,
														Required: true,
													},
													"comment": {
														Type:     schema.TypeString,
														Optional: true,
													},
													"criteria_match": {
														Type:     schema.TypeString,
														Optional: true,
														Default:  "all",
													},
													"criteria": akpsCriteria,
													"behavior": akpsBehavior,
													"rule": &schema.Schema{
														Type:     schema.TypeSet,
														Optional: true,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"name": {
																	Type:     schema.TypeString,
																	Required: true,
																},
																"comment": {
																	Type:     schema.TypeString,
																	Optional: true,
																},
																"criteria_match": {
																	Type:     schema.TypeString,
																	Optional: true,
																	Default:  "all",
																},
																"criteria": akpsCriteria,
																"behavior": akpsBehavior,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"variable": &schema.Schema{
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"description": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"hidden": {
								Type:     schema.TypeBool,
								Required: true,
							},
							"sensitive": {
								Type:     schema.TypeBool,
								Required: true,
							},
							"value": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	},
	"json": {
		Type: schema.TypeString,
		//Type: schema.TypeSet,
		Computed:    true,
		Description: "JSON Rule representation",
	},
}

func resourcePropertyRulesCreate(d *schema.ResourceData, meta interface{}) error {

	rules := papi.NewRules()

	// get rules from the TF config
	unmarshalRules(d, rules)

	jsonBody, err := jsonhooks.Marshal(rules)
	if err != nil {
		return err
	}

	sha := getSHAString(string(jsonBody))
	d.Set("json", string(jsonBody))

	d.SetId(sha)

	log.Println("[DEBUG] Done")
	return nil
}

func resourcePropertyRulesDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DELETING")
	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourcePropertyRulesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceID := d.Id()
	propertyID := resourceID

	if !strings.HasPrefix(resourceID, "prp_") {
		for _, searchKey := range []papi.SearchKey{papi.SearchByPropertyName, papi.SearchByHostname, papi.SearchByEdgeHostname} {
			results, err := papi.Search(searchKey, resourceID)
			if err != nil {
				continue
			}

			if results != nil && len(results.Versions.Items) > 0 {
				propertyID = results.Versions.Items[0].PropertyID
				break
			}
		}
	}

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	e := property.GetProperty()
	if e != nil {
		return nil, e
	}

	d.Set("account", property.AccountID)
	d.Set("contract", property.ContractID)
	d.Set("group", property.GroupID)

	d.Set("name", property.PropertyName)
	d.Set("version", property.LatestVersion)
	d.SetId(property.PropertyID)

	return []*schema.ResourceData{d}, nil
}

func resourcePropertyRulesExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func resourcePropertyRulesRead(d *schema.ResourceData, meta interface{}) error {

	log.Println("[DEBUG] READ" + d.Id())

	return nil
}

func resourcePropertyRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATING")
	d.Partial(true)
	rules := papi.NewRules()

	// get rules from the TF config
	unmarshalRules(d, rules)

	jsonBody, err := jsonhooks.Marshal(rules)
	if err != nil {
		return err
	}

	sha := getSHAString(string(jsonBody))
	d.Set("json", string(jsonBody))
	//d.Set("json", jsonBody)
	d.SetId(sha)
	d.Partial(false)

	log.Println("[DEBUG] Done")
	return nil
}
