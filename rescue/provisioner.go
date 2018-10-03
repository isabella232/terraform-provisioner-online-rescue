package rescue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// TokenEnvVar is the environment variable name used to get the Online.net API token
const TokenEnvVar = "ONLINE_TOKEN"

// Provisioner describes this provisioner configuration.
func Provisioner() terraform.ResourceProvisioner {
	return &schema.Provisioner{
		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "This will either enable or disable rescue mode",
			},
			"server": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The server to set the rescue mode on",
			},
			"image": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This image to boot rescue mode into (required if enabled is true)",
			},
			// This is a hack used to send the credentials back to where this is called inside Terraform
			// provisioners currently do not support output (https://github.com/hashicorp/terraform/issues/610)
			// this is why we will write the rescue details to disk in individual files in a given directory
			// so they can easily be read by ${file()}
			"credentials_dir": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This path will be used to store the rescue credentials in as separate files (required if enabled is true)",
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc(TokenEnvVar, nil),
				Required:    true,
				Sensitive:   true,
				Description: "Online.net private API token, by default the ONLINE_TOKEN environment variable is used.",
			},
		},
		ApplyFunc: applyRescueProvisioner,
	}
}
