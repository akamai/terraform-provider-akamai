package akamai

import (
	"encoding/json"
	"errors"
	"log"
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
	property.PropertyID = d.Get("property").(string)
	err := property.GetProperty()
	if err != nil {
		return errors.New("unable to find property")
	}

	// The API now has data, so save the partial state
	d.SetPartial("network")
	d.Set("property", property.PropertyID)

	if d.Get("activate").(bool) {
		activation, err := activateProperty(property, d)
		if err != nil {
			return err
		}

		d.SetId(activation.ActivationID)
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
	property.PropertyID = d.Get("property").(string)
	e := property.GetProperty()
	if e != nil {
		return e
	}

	log.Printf("[DEBUG] DEACTIVE PROPERTY %v", property)

	activations, e := property.GetActivations()
	if e != nil {
		return e
	}

	network := papi.NetworkValue(d.Get("network").(string))
	version := d.Get("version").(int)
	for _, activation := range activations.Activations.Items {
		if activation.Network == network && activation.PropertyVersion == version && activation.Status != papi.StatusInactive && activation.Status != papi.StatusDeactivated {
			// The version is not inactive, so we need to deactivate it
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
	}

	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourcePropertyActivationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Get("property").(string)
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
	property.PropertyID = d.Get("property").(string)
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

	property, e := getProperty(d)
	if e != nil {
		return e
	}

	if d.Get("activate").(bool) {
		activations, err := property.GetActivations()
		if err != nil {
			// No activations found
			return nil
		}

		old, new := d.GetChange("network")
		if old.(string) != new.(string) {
			// deactivate on the old network, we don't need to wait for this
			deactivateProperty(property, d, papi.NetworkValue(old.(string)))
		}

		var activation *papi.Activation
		network := papi.NetworkValue(d.Get("network").(string))
		version := d.Get("version").(int)
		for _, a := range activations.Activations.Items {
			if a.Network == network && a.PropertyVersion == version && a.Status != papi.StatusFailed && a.Status != papi.StatusDeactivated && a.Status != papi.StatusAborted && a.ActivationType != papi.ActivationTypeDeactivate {
				activation = a
				break
			}
		}

		// Already an activation in process, we'll just wait on that one
		if activation == nil {
			activation, err = activateProperty(property, d)
			if err != nil {
				return err
			}
		}

		d.SetId(activation.ActivationID)
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
		d.Set("status", string(activation.Status))
	} else {
		return resourcePropertyRead(d, meta)
	}

	log.Println("[DEBUG] Done")
	return nil
}

func activateProperty(property *papi.Property, d *schema.ResourceData) (*papi.Activation, error) {
	activation := getActivation(d, papi.ActivationTypeActivate, papi.NetworkValue(d.Get("network").(string)))
	err := activation.Save(property, true)
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

	activation := getActivation(d, papi.ActivationTypeDeactivate, network)
	err = activation.Save(property, true)
	if err != nil {
		body, _ := json.Marshal(activation)
		log.Printf("[DEBUG] API Request Body: %s\n", string(body))
		return nil, err
	}
	log.Println("[DEBUG] Deactivation submitted successfully")

	return activation, nil
}

func getActivation(d *schema.ResourceData, activationType papi.ActivationValue, network papi.NetworkValue) *papi.Activation {
	log.Println("[DEBUG] Creating new activation")
	activation := papi.NewActivation(papi.NewActivations())
	activation.PropertyVersion = d.Get("version").(int)
	activation.Network = network
	for _, email := range d.Get("contact").(*schema.Set).List() {
		activation.NotifyEmails = append(activation.NotifyEmails, email.(string))
	}
	activation.Note = "Using Terraform"

	activation.ActivationType = activationType

	log.Println("[DEBUG] Activating")
	return activation
}
