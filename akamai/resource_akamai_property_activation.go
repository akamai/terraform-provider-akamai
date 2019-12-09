package akamai

import (
	"encoding/json"
	"errors"
	"fmt"
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
		Schema: akamaiPropertyActivationSchema,
	}
}

var akamaiPropertyActivationSchema = map[string]*schema.Schema{
	"property": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"version": &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
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
	"contact": &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"status": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
}

func resourcePropertyActivationCreate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	property := papi.NewProperty(papi.NewProperties())
	id := strings.Split(d.Get("property").(string), "-")
	property.PropertyID = id[0]
	err := property.GetProperty()
	if err != nil {
		return errors.New("unable to find property")
	}

	// The API now has data, so save the partial state
	d.SetPartial("network")
	d.Set("property", property.PropertyID+"-"+id[1])

	if d.Get("activate").(bool) {
		activation, err := activateProperty(property, d)
		if err != nil {
			return err
		}

		d.SetId(activation.ActivationID)
		d.Set("version", activation.PropertyVersion)
		d.Set("status", string(activation.Status))
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
	} else {
		d.SetId("none")
	}

	d.Partial(false)
	log.Println("[DEBUG] Done")
	return nil
}

func resourcePropertyActivationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DEACTIVATE PROPERTY")

	property := papi.NewProperty(papi.NewProperties())
	id := strings.Split(d.Get("property").(string), "-")
	property.PropertyID = id[0]
	e := property.GetProperty()
	if e != nil {
		return e
	}

	log.Printf("[DEBUG] DEACTIVE PROPERTY %v", property)

	network := papi.NetworkValue(d.Get("network").(string))
	propertyVersion := property.ProductionVersion
	if network == "STAGING" {
		propertyVersion = property.StagingVersion
	}
	version := d.Get("version").(int)
	log.Printf("[DEBUG] Version to deactivate is %d and current active %s version is %d\n", version, network, propertyVersion)

	if propertyVersion == version {
		// The current active version is the one we need to deactivate
		log.Printf("[DEBUG] Deactivating %s version %d \n", network, version)
		activation, err := deactivateProperty(property, d, papi.NetworkValue(d.Get("network").(string)))
		if err != nil {
			return err
		}

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

		d.Set("status", string(activation.Status))
	}

	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourcePropertyActivationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	property := papi.NewProperty(papi.NewProperties())
	id := strings.Split(d.Get("property").(string), "-")
	property.PropertyID = id[0]
	err := property.GetProperty()
	if err != nil {
		return false, err
	}

	activations, err := property.GetActivations()
	if err != nil {
		// No activations found
		return false, nil
	}

	network := papi.NetworkValue(d.Get("network").(string))
	version := d.Get("version").(int)
	for _, activation := range activations.Activations.Items {
		if activation.Network == network && activation.PropertyVersion == version {
			return true, nil
		}
	}

	return false, nil
}

func resourcePropertyActivationRead(d *schema.ResourceData, meta interface{}) error {
	property := papi.NewProperty(papi.NewProperties())
	id := strings.Split(d.Get("property").(string), "-")
	property.PropertyID = id[0]
	err := property.GetProperty()
	if err != nil {
		return err
	}

	d.SetId("")
	activations, err := property.GetActivations()
	if err != nil {
		// No activations found
		return nil
	}

	network := papi.NetworkValue(d.Get("network").(string))
	version := d.Get("version").(int)
	for _, activation := range activations.Activations.Items {
		if activation.Network == network && activation.PropertyVersion == version {
			d.SetId(activation.ActivationID)
			d.Set("status", string(activation.Status))
		}
	}

	return nil
}

func resourcePropertyActivationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATING")

	log.Println("[DEBUG] Fetching property")
	id := strings.Split(d.Get("property").(string), "-")
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = id[0]
	e := property.GetProperty()
	if e != nil {
		return e
	}

	activation, err := getActivation(d, property, papi.ActivationTypeActivate, papi.NetworkValue(d.Get("network").(string)))
	if err != nil {
		return err
	}

	a, err := findExistingActivation(property, activation)
	if err == nil {
		activation = a
	}

	if d.Get("activate").(bool) {
		old, new := d.GetChange("network")
		if old.(string) != new.(string) {
			// deactivate on the old network, we don't need to wait for this
			deactivateProperty(property, d, papi.NetworkValue(old.(string)))
		}

		// No activation in progress, create a new one
		if a == nil {
			activation, err = activateProperty(property, d)
			if err != nil {
				return err
			}
		}

		d.SetId(activation.ActivationID)
		d.Set("version", activation.PropertyVersion)
		d.Set("status", string(activation.Status))

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
		d.Set("version", activation.PropertyVersion)
		d.Set("status", string(activation.Status))
	} else {
		return resourcePropertyRead(d, meta)
	}

	log.Println("[DEBUG] Done")
	return nil
}

func activateProperty(property *papi.Property, d *schema.ResourceData) (*papi.Activation, error) {
	activation, err := getActivation(d, property, papi.ActivationTypeActivate, papi.NetworkValue(d.Get("network").(string)))
	if err != nil {
		return nil, err
	}

	if a, err := findExistingActivation(property, activation); err == nil && a != nil {
		return a, nil
	}

	err = activation.Save(property, true)
	if err != nil {
		body, _ := json.Marshal(activation)
		log.Printf("[DEBUG] API Request Body: %s\n", string(body))
		return nil, err
	}
	log.Println("[DEBUG] Activation submitted successfully")

	return activation, nil
}

func deactivateProperty(property *papi.Property, d *schema.ResourceData, network papi.NetworkValue) (*papi.Activation, error) {
	version, err := property.GetLatestVersion(network)
	if err != nil || version == nil {
		// Not active
		return nil, nil
	}

	activation, err := getActivation(d, property, papi.ActivationTypeDeactivate, network)
	if err != nil {
		return nil, err
	}

	if a, err := findExistingActivation(property, activation); err == nil && a != nil {
		return a, nil
	}

	err = activation.Save(property, true)
	if err != nil {
		body, _ := json.Marshal(activation)
		log.Printf("[DEBUG] API Request Body: %s\n", string(body))
		return nil, err
	}
	log.Println("[DEBUG] Deactivation submitted successfully")

	return activation, nil
}

func getActivation(d *schema.ResourceData, property *papi.Property, activationType papi.ActivationValue, network papi.NetworkValue) (*papi.Activation, error) {
	log.Println("[DEBUG] Creating new activation")
	activation := papi.NewActivation(papi.NewActivations())
	if version, ok := d.GetOk("version"); ok && version.(int) != 0 {
		activation.PropertyVersion = version.(int)
	} else {
		version, err := property.GetLatestVersion("")
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] Using latest version: %d\n", version.PropertyVersion)
		activation.PropertyVersion = version.PropertyVersion
	}
	activation.Network = network
	for _, email := range d.Get("contact").(*schema.Set).List() {
		activation.NotifyEmails = append(activation.NotifyEmails, email.(string))
	}
	activation.Note = "Using Terraform"

	activation.ActivationType = activationType

	log.Println("[DEBUG] Activating")
	return activation, nil
}

func findExistingActivation(property *papi.Property, activation *papi.Activation) (*papi.Activation, error) {
	activations, err := property.GetActivations()
	if err != nil {
		return nil, err
	}

	inProgressStates := map[papi.StatusValue]bool{
		papi.StatusActive:              true,
		papi.StatusNew:                 true,
		papi.StatusPending:             true,
		papi.StatusPendingDeactivation: true,
		papi.StatusZone1:               true,
		papi.StatusZone2:               true,
		papi.StatusZone3:               true,
	}
	for _, a := range activations.Activations.Items {
		if _, ok := inProgressStates[a.Status]; !ok {
			continue
		}

		// There is an activation in progress, if it's for the same version/network/type we can re-use it
		if a.PropertyVersion != activation.PropertyVersion || a.ActivationType != activation.ActivationType || a.Network != activation.Network {
			return nil, fmt.Errorf("%s already in progress: v%d on %s", activation.ActivationType, activation.PropertyVersion, activation.Network)
		}

		log.Println("[DEBUG] Existing activation found")
		return a, nil
	}

	return nil, nil
}
