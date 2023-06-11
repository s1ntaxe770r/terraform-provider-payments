package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/s1ntaxe770r/terraform-provider-payments/payments"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run provider with debug support")

	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return payments.Provider()
		},
	}

	plugin.Serve(opts)
}
