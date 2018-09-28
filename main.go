package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/src-d/terraform-provisioner-online-rescue/rescue"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProvisionerFunc: rescue.Provisioner,
	})
}
