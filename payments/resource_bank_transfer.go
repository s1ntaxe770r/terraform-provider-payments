package payments

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/s1ntaxe770r/terraform-provider-payments/pkg/client"
)

func resourceBankTransfer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBankTransferCreate,
		Read:   schema.Noop,
		Delete: schema.Noop,
		Schema: map[string]*schema.Schema{
			"amount": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bank_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_number": &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
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
			"response_code": &schema.Schema{
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
		return err
	}
	if err := d.Set("message", resp.Message); err != nil {
		return errors.New("error setting message")
	}
	if err := d.Set("status", resp.Status); err != nil {
		return errors.New("error setting status")
	}
	if err := d.Set("request_reference", resp.RequestReference); err != nil {
		return errors.New("error setting request_reference")
	}
	if err := d.Set("response_code", resp.ResponseCode); err != nil {
		return errors.New("error setting response_code")
	}

	d.SetId(resp.RequestReference)
	return nil
}
