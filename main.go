package main

import (
	"flag"

	"github.com/philips-software/terraform-provider-hsdp/hsdp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
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
		return hsdp.Provider(buildVersion)
	}}
	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "registry.terraform.io/philips-software/hsdp"
	}
	plugin.Serve(opts)
}
