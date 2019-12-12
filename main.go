package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/philips-software/terraform-provider-hsdp/hsdp"
)

var commit = "deadbeef"
var release = "v0.0.0"
var buildVersion = release + "-" + commit

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return hsdp.Provider(buildVersion)
		},
	})
}
