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
				Description: "email associated with your kuda buisness 	account",
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "account number associated with your kuda buisness account",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"payments_banks": dataSourceBanks(),
		},
	}
}
