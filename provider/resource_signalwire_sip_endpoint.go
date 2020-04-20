package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSignalwireSipEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceSignalwireSipEndpointCreate,
		Read:   resourceSignalwireSipEndpointRead,
		Update: resourceSignalwireSipEndpointUpdate,
		Delete: resourceSignalwireSipEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSignalwireSipEndpointImport,
		},
		Schema: map[string]*schema.Schema{
			"space": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The space your SIP endpoint will be created in (https://your-space.signalwire.com)",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "String representing the username portion of the endpoint. Must be unique across your project and must not container white space characters or @.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A password to authenticate registrations to this endpoint.",
			},
			"caller_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Friendly Caller ID used as the CNAM when dialing a phone number or the From when dialing another SIP Endpoint.",
			},
			//"send_as": {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	Description: "The e164 formatted number you which to set as the originating number when dialing PSTN phone numbers from this SIP Endpoint.",
			//},
			"ciphers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "A list of encryption ciphers this endpoint will support.",
			},
			"codecs": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "A list of codecs this endpoint will support.",
			},
			"encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A string representing whether connections to this endpoint require encryption or if encryption is optional. Encryption will always be used if possible.",
			},
		},
	}
}

func resourceSignalwireSipEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	resp, err := client.Req("POST", d.Get("space").(string), "endpoints/sip", schemaToSipEndpoint(d))

	if err != nil {
		return err
	}

	d.SetId(resp["id"].(string))

	return resourceSignalwireSipEndpointRead(d, meta)
}

func resourceSignalwireSipEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	resp, err := client.Req("GET", d.Get("space").(string), "endpoints/sip/"+d.Id(), nil)

	if err != nil {
		return err
	}

	respToSipEndpointData(d, resp)
	return nil
}

func resourceSignalwireSipEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	resp, err := client.Req("PUT", d.Get("space").(string), "endpoints/sip/"+d.Id(), schemaToSipEndpoint(d))

	if err != nil {
		return err
	}

	respToSipEndpointData(d, resp)
	return nil
}

func resourceSignalwireSipEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	_, err := client.Req("DELETE", d.Get("space").(string), "endpoints/sip/"+d.Id(), nil)

	return err
}

func schemaToSipEndpoint(d *schema.ResourceData) map[string]interface{} {
	request := map[string]interface{}{
		"username": d.Get("username").(string),
		"password": d.Get("password").(string),
	}

	optionalsStr := []string{"caller_id", "encryption"}
	for _, optional := range optionalsStr {
		if val, ok := d.GetOkExists(optional); ok {
			request[optional] = val.(string)
		}
	}
	optionalsSlices := []string{"ciphers", "codecs"}
	for _, optionalSlice := range optionalsSlices {
		if val, ok := d.GetOkExists(optionalSlice); ok {
			var strSlice []string
			for _, strVal := range val.([]interface{}) {
				strSlice = append(strSlice, strVal.(string))
			}
			request[optionalSlice] = strSlice
		}
	}
	return request
}

func respToSipEndpointData(d *schema.ResourceData, resp map[string]interface{}) {
	d.Set("username", resp["username"].(string))
	d.Set("caller_id", resp["caller_id"].(string))
	d.Set("encryption", resp["encryption"].(string))

	slices := []string{"ciphers", "codecs"}
	for _, slice := range slices {
		var strSlice []string
		for _, strVal := range resp[slice].([]interface{}) {
			strSlice = append(strSlice, strVal.(string))
		}
		d.Set(slice, strSlice)
	}
}
