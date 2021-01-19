package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataFile(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataFileConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.foo", "bar"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.bar", "1"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.baz.0", "qux"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.baz.1", "quux"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.qux.0", "1"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.qux.1", "2"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.quux.quuz", "corge"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.quux.grault", "garply"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.quuz.corge", "grault"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.quuz.garply.0", "waldo"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.quuz.garply.1", "fred"),
					resource.TestCheckResourceAttr("data.tfvars_file.test", "variables.corge.0.grault.garply", "waldo"),
				),
			},
		},
	})
}

const testDataFileConfig = `
data "tfvars_file" "test" {
  filename = "testdata/test.tfvars"
}
`
