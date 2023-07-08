package main

import (
	"os"

	"github.com/s1ntaxe770r/terraform-provider-payments/pkg/client"
	"github.com/sirupsen/logrus"
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
	apikey := os.Getenv("API_KEY")
	c := client.NewClient(apikey, "hello@jubril.xyz")
	token := c.GetAuthToken()

	name, err := c.GetSenderName("2182383852", "000004", token)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info(name)

}
