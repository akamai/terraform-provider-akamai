package akamai

import (
	"errors"

	"log"

	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePropertyActivation() *schema.Resource {
	return &schema.Resource{
		Create: resourcePropertyActivationCreate,
		Read:   resourcePropertyActivationRead,
		Update: resourcePropertyActivationUpdate,
		Delete: resourcePropertyActivationDelete,
		Exists: resourcePropertyActivationExists,
		Importer: &schema.ResourceImporter{
			State: resourcePropertyActivationImport,
		},
		Schema: akamaiPropertyActivationSchema,
	}
}

func resourcePropertyActivationCreate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	group, e := getGroup(d)
	if e != nil {
		return e
	}

	contract, e := getContract(d)
	if e != nil {
		return e
	}

	var property *papi.Property
	if property = findProperty(d); property == nil {
		if group == nil {
			return errors.New("group_id must be specified to activate a new property")
		}

		if contract == nil {
			return errors.New("contract_id must be specified to activate a new property")
		}

		/*
		       var e error
		   		property, e = findPropertyActivation(d)
		   		if e != nil {
		   			return e
		   		}*/
	}

	err := ensureEditableVersion(property)
	if err != nil {
		return err
	}
	d.Set("account", property.AccountID)
	d.Set("version", property.LatestVersion)

	// The API now has data, so save the partial state
	d.SetId(property.PropertyID)
	d.SetPartial("name")
	d.SetPartial("contract")
	d.SetPartial("group")
	d.SetPartial("network")

	if d.Get("activate").(bool) {
		activation, err := activateProperty(property, d)
		if err != nil {
			return err
		}
		d.SetPartial("contact")

		go activation.PollStatus(property)

	polling:
		//for activation.Status != papi.StatusActive {
		for activation.Status != "Created" {
			select {
			case statusChanged := <-activation.StatusChange:
				log.Printf("[DEBUG] Property Status: %s\n", activation.Status)
				if statusChanged == false {
					break polling
				}
				continue polling
			case <-time.After(time.Minute * 90):
				log.Println("[DEBUG] Activation Timeout (90 minutes)")
				break polling
			}
		}
	}

	d.Partial(false)
	log.Println("[DEBUG] Done")
	return nil
}

func resourcePropertyActivationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] CANCEL ACTIVATION")

	contractID, ok := d.GetOk("contract")
	if !ok {
		return errors.New("missing contract ID")
	}
	groupID, ok := d.GetOk("group")
	if !ok {
		return errors.New("missing group ID")
	}
	network, ok := d.GetOk("network")
	if !ok {
		return errors.New("missing network")
	}

	propertyID := d.Id()

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	property.Contract = &papi.Contract{ContractID: contractID.(string)}
	property.Group = &papi.Group{GroupID: groupID.(string)}

	e := property.GetProperty()
	if e != nil {
		return e
	}

	log.Printf("[DEBUG] CANCEL ACTIVATION PROPERTY %v", property)

	activations, e := property.GetActivations()
	if e != nil {
		return e
	}
	log.Printf("[DEBUG] CANCEL ACTIVATION activations %v", activations)
	activation, e := activations.GetLatestActivation(papi.NetworkValue(strings.ToUpper(network.(string))), papi.StatusActive)
	log.Printf("[DEBUG] CANCEL ACTIVATION activations get latest %v", activations)
	// an error here means there has not been any activation yet, so we can skip deactivating the property
	// if there was no error, then activations were found, this can be an Activation or a Deactivation, so we check the ActivationType
	// in case it has already been deactivated
	if e == nil && activation.ActivationType == papi.ActivationTypeActivate {
		log.Printf("[DEBUG] CANCEL ACTIVATION deactivations ")
		deactivation := papi.NewActivation(papi.NewActivations())
		deactivation.PropertyVersion = property.LatestVersion
		deactivation.ActivationType = papi.ActivationTypeDeactivate
		deactivation.Network = activation.Network
		deactivation.NotifyEmails = activation.NotifyEmails
		log.Printf("[DEBUG] CANCEL ACTIVATION deactivations %v", deactivation)
		e = deactivation.Save(property, true)
		if e != nil {
			return e
		}
		log.Printf("[DEBUG] DEACTIVATION SAVED - ID %s STATUS %s\n", deactivation.ActivationID, deactivation.Status)

		go deactivation.PollStatus(property)

	polling:
		for deactivation.Status != papi.StatusActive {
			select {
			case statusChanged := <-deactivation.StatusChange:
				log.Printf("[DEBUG] Property Status: %s\n", deactivation.Status)
				if statusChanged == false {
					break polling
				}
				continue polling
			case <-time.After(time.Minute * 90):
				log.Println("[DEBUG] Deactivation Timeout (90 minutes)")
				break polling
			}
		}
	}


	e = property.Delete()
	if e != nil {
		return e
	}
	
	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourcePropertyActivationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
	//d.Set("clone_from", property.CloneFrom.PropertyID)
	d.Set("name", property.PropertyName)
	d.Set("version", property.LatestVersion)
	d.SetId(property.PropertyID)

	return []*schema.ResourceData{d}, nil
}

func resourcePropertyActivationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Id()
	e := property.GetProperty()
	if e != nil {
		return false, e
	}

	return true, nil
}

func resourcePropertyActivationRead(d *schema.ResourceData, meta interface{}) error {
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Id()
	err := property.GetProperty()
	if err != nil {
		return err
	}

	// Cannot set clone_from. Not provided on GET requests.
	// d.Set("clone_from", nil)

	// Cannot set product. Not provided on GET requests.
	// d.Set("product", property.ProductID)

	d.Set("account", property.AccountID)
	d.Set("contract", property.ContractID)
	d.Set("group", property.GroupID)
	d.Set("name", property.PropertyName)
	//d.Set("note", property.Note)

	d.Set("version", property.LatestVersion)
	if property.StagingVersion > 0 {
		d.Set("staging_version", property.StagingVersion)
	}
	if property.ProductionVersion > 0 {
		d.Set("production_version", property.ProductionVersion)
	}

	return nil
}

var akamaiPropertyActivationSchema = map[string]*schema.Schema{
	"account": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"contract": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"group": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"network": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "staging",
	},

	"activate": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"version": &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	},
	"staging_version": &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	},
	"production_version": &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	},
	/*"rule_format": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"ipv6": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	},*/
	"hostname": &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"contact": &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
}

func resourcePropertyActivationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATING")
	d.Partial(true)

	property, e := getProperty(d)
	if e != nil {
		return e
	}

	err := ensureEditableVersion(property)
	if err != nil {
		return err
	}
	d.Set("version", property.LatestVersion)

	d.SetPartial("default")
	d.SetPartial("origin")

	// an existing activation on this property will be automatically deactivated upon
	// creation of this new activation
	if d.Get("activate").(bool) {
		activation, err := activateProperty(property, d)
		if err != nil {
			return err
		}
		d.SetPartial("contact")

		go activation.PollStatus(property)

	polling:
		for activation.Status != papi.StatusActive {
			select {
			case statusChanged := <-activation.StatusChange:
				log.Printf("[DEBUG] Property Status: %s\n", activation.Status)
				if statusChanged == false {
					break polling
				}
				continue polling
			case <-time.After(time.Minute * 90):
				log.Println("[DEBUG] Activation Timeout (90 minutes)")
				break polling
			}
		}
	}

	d.Partial(false)

	log.Println("[DEBUG] Done")
	return nil
}

// Helpers

func findPropertyActivation(d *schema.ResourceData) (*papi.Property, error) {
	results, err := papi.Search(papi.SearchByPropertyName, d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if results == nil || len(results.Versions.Items) == 0 {
		for _, hostname := range d.Get("hostname").(*schema.Set).List() {
			results, err = papi.Search(papi.SearchByHostname, hostname.(string))
			if err == nil && (results == nil || len(results.Versions.Items) != 0) {
				break
			}
		}

		if err != nil || (results == nil || len(results.Versions.Items) == 0) {
			return nil, err
		}
	}

	property := &papi.Property{
		PropertyID: results.Versions.Items[0].PropertyID,
		Group: &papi.Group{
			GroupID: results.Versions.Items[0].GroupID,
		},
		Contract: &papi.Contract{
			ContractID: results.Versions.Items[0].ContractID,
		},
	}

	err = property.GetProperty()
	if err != nil {
		return nil, err
	}

	return property, err
}
