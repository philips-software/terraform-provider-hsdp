package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/philips-software/terraform-provider-hsdp/hsdp"
)

var commit = "deadbeef"
var release = "v0.0.0"
var date = "0000-00-00"
var buildSource = "unknown"
var buildVersion = release + "-" + commit + "." + date + "." + buildSource

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return hsdp.Provider(buildVersion)
		},
	})
}
