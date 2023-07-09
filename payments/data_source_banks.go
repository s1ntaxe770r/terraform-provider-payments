package payments

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/s1ntaxe770r/terraform-provider-payments/pkg/client"
)

func dataSourceBanks() *schema.Resource {
	return &schema.Resource{
		Description: "obtain a list of banks and their codes",
		Read:        dataSourceBanksRead,
		Schema: map[string]*schema.Schema{
			"banks": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBanksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.ApiClient)
	authToken := client.GetAuthToken()
	bankListResponse, err := client.GetBankList(authToken)

	if err != nil {
		return err
	}

	banks := flattenBanks(bankListResponse.Data.Banks)

	d.SetId("banks") // Set a unique ID for the resource data
	d.Set("banks", banks)

	return nil
}

func flattenBanks(banks []client.Bank) (values []map[string]interface{}) {
	if banks != nil {
		for _, bank := range banks {
			value := map[string]interface{}{
				"name": bank.BankName,
				"code": bank.BankCode,
			}
			values = append(values, value)
		}

	}
	return values
}
