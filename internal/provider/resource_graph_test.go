package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var id = fmt.Sprintf("todo-api-%s", uuid.New().String())

func TestAccGraphResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGraphResourceConfigDefault(id, "Todo API"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_graph.test", "id", id),
					resource.TestCheckResourceAttr("apollographql_graph.test", "title", "Todo API"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "onboarding_architecture", "MONOLITH"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "organization_id", "pksunkara"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "description", ""),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apollographql_graph.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccGraphResourceConfigDefault(id, "Todo API"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_graph.test", "id", id),
					resource.TestCheckResourceAttr("apollographql_graph.test", "title", "Todo API"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "onboarding_architecture", "MONOLITH"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "organization_id", "pksunkara"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "description", ""),
				),
			},
			// Update and Read testing
			{
				Config: testAccGraphResourceConfigDefaultUpdate(id, "Todo app API"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_graph.test", "id", id),
					resource.TestCheckResourceAttr("apollographql_graph.test", "title", "Todo app API"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "onboarding_architecture", "MONOLITH"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "organization_id", "pksunkara"),
					resource.TestCheckResourceAttr("apollographql_graph.test", "description", "API for our todo app"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apollographql_graph.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGraphResourceConfigDefault(id string, title string) string {
	return fmt.Sprintf(`
resource "apollographql_graph" "test" {
  id = "%s"
  title = "%s"

  organization_id = "pksunkara"
}
`, id, title)
}

func testAccGraphResourceConfigDefaultUpdate(id string, title string) string {
	return fmt.Sprintf(`
resource "apollographql_graph" "test" {
  id = "%s"
  title = "%s"

  organization_id = "pksunkara"
  description = "API for our todo app"
}
`, id, title)
}
