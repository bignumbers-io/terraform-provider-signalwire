package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SIGNALWIRE_PROJECT_ID", nil),
				Description: "Signalwire Project ID",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SIGNALWIRE_AUTH_TOKEN", "https://kitt.signalwire.com/credentials/new"),
				Description: "The API token used to connect to Signalwire",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"signalwire_sip_endpoint": resourceSignalwireSipEndpoint(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token:     d.Get("auth_token").(string),
		ProjectId: d.Get("project_id").(string),
	}

	return config.Client()
}
