package payments

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/s1ntaxe770r/terraform-provider-payments/pkg/client"
)

func resourceBankTransfer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBankTransferCreate,
		Schema: map[string]*schema.Schema{
			"amount": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bank_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"message": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"request_reference": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBankTransferCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.ApiClient)
	authToken := c.GetAuthToken()
	amount := d.Get("amount").(string)
	bankCode := d.Get("bank_code").(string)
	accountNumber := d.Get("account_number").(string)

	resp, err := c.SingleFundTransfer(accountNumber, amount, bankCode, authToken)

	if err != nil {
		return errors.New("unable to make transfer" + err.Error())

	}
	if err := d.Set("message", resp.Message); err != nil {
		return errors.New("unable to set message" + err.Error())
	}

	if err := d.Set("status", resp.Status); err != nil {
		return errors.New("unable to set status" + err.Error())
	}

	if err := d.Set("request_reference", resp.RequestReference); err != nil {
		return errors.New("unable to set request reference" + err.Error())
	}
	return nil
}
