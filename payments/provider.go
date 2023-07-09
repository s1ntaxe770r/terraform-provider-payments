package payments

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/s1ntaxe770r/terraform-provider-payments/pkg/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Description: "api token",
				Required:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "email associated with your kuda buisness 	account",
				Required:    true,
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "account number associated with your kuda buisness account",
				Required:    true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"payments_banks": dataSourceBanks(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"payments_bank_transfer": resourceBankTransfer(),
		},
		ConfigureContextFunc: proivdeConfigure,
	}
}

func proivdeConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	email := d.Get("email").(string)
	accountNumber := d.Get("account_number").(string)
	apiToken := d.Get("api_token").(string)

	var dg diag.Diagnostics
	c := client.NewClient(apiToken, email, accountNumber)
	return c, dg
}
