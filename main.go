package main

import (
	"fmt"

	"github.com/s1ntaxe770r/terraform-provider-payments/pkg/client"
	"golang.org/x/exp/slog"
)

func main() {
	// var debugMode bool

	// flag.BoolVar(&debugMode, "debug", false, "set to true to run provider with debug support")

	// flag.Parse()

	// opts := &plugin.ServeOpts{
	// 	ProviderFunc: func() *schema.Provider {
	// 		return payments.Provider()
	// 	},
	// }

	// plugin.Serve(opts)

	c := client.NewClient("", "hello@jubril.xyz")

	resp, err := c.GetBankList()

	if err != nil {
		fmt.Println(err)
	}

	slog.Info(resp.Data.Banks[0].BankName)
	slog.Info(resp.Data.Banks[0].BankCode)

}
