package main

import (
	"context"
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/philips-software/terraform-provider-hsdp/internal/provider"

	"log"
)

var commit = "deadbeef"
var release = "v0.0.0"
var date = "0000-00-00"
var buildSource = "unknown"
var buildVersion = release + "-" + commit + "." + date + "." + buildSource

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: func() *schema.Provider {
		return provider.Provider(buildVersion)
	}}
	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/philips-software/hsdp", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
	plugin.Serve(opts)
}
