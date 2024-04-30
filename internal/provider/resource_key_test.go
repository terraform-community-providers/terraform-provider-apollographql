package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeyResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccKeyResourceConfigDefault("deploy"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apollographql_key.test", "id"),
					resource.TestCheckResourceAttr("apollographql_key.test", "name", "deploy"),
					resource.TestCheckResourceAttr("apollographql_key.test", "role", "GRAPH_ADMIN"),
					resource.TestCheckResourceAttrSet("apollographql_key.test", "token"),
					resource.TestCheckResourceAttr("apollographql_key.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// Update with default values
			{
				Config: testAccKeyResourceConfigDefault("deploy"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apollographql_key.test", "id"),
					resource.TestCheckResourceAttr("apollographql_key.test", "name", "deploy"),
					resource.TestCheckResourceAttr("apollographql_key.test", "role", "GRAPH_ADMIN"),
					resource.TestCheckResourceAttrSet("apollographql_key.test", "token"),
					resource.TestCheckResourceAttr("apollographql_key.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// Update and Read testing
			{
				Config: testAccKeyResourceConfigDefault("github-deploy"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apollographql_key.test", "id"),
					resource.TestCheckResourceAttr("apollographql_key.test", "name", "github-deploy"),
					resource.TestCheckResourceAttr("apollographql_key.test", "role", "GRAPH_ADMIN"),
					resource.TestCheckResourceAttrSet("apollographql_key.test", "token"),
					resource.TestCheckResourceAttr("apollographql_key.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccKeyResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "apollographql_key" "test" {
  name = "%s"
  graph_id = "Test-w4a5n4"
}
`, name)
}
