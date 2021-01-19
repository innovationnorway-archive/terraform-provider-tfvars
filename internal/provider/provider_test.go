package provider

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var protoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"tfvars": func() (tfprotov5.ProviderServer, error) {
		return New("dev")(), nil
	},
}
