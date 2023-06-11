package payments

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Description: "api token",
			},
			"email": {
				Type:        schema.TypeString,
				Description: "email associated with your kuda account",
			},
		},
	}
}
