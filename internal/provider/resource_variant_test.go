package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVariantResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVariantResourceConfigDefault("Staging"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_variant.test", "id", "Test-w4a5n4@Staging"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "name", "Staging"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "public", "false"),
					resource.TestCheckNoResourceAttr("apollographql_variant.test", "url"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apollographql_variant.test",
				ImportState:       true,
				ImportStateId:     "Test-w4a5n4:Staging",
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccVariantResourceConfigDefault("Staging"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_variant.test", "id", "Test-w4a5n4@Staging"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "name", "Staging"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "public", "false"),
					resource.TestCheckNoResourceAttr("apollographql_variant.test", "url"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// Update and Read testing
			{
				Config: testAccVariantResourceConfigNonDefault("Staging"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_variant.test", "id", "Test-w4a5n4@Staging"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "name", "Staging"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "public", "true"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "url", "https://example.com"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apollographql_variant.test",
				ImportState:       true,
				ImportStateId:     "Test-w4a5n4:Staging",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccVariantResourceNonDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVariantResourceConfigNonDefault("QA"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_variant.test", "id", "Test-w4a5n4@QA"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "name", "QA"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "public", "true"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "url", "https://example.com"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apollographql_variant.test",
				ImportState:       true,
				ImportStateId:     "Test-w4a5n4:QA",
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccVariantResourceConfigNonDefault("QA"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_variant.test", "id", "Test-w4a5n4@QA"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "name", "QA"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "public", "true"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "url", "https://example.com"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// Update and Read testing
			{
				Config: testAccVariantResourceConfigDefault("QA"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollographql_variant.test", "id", "Test-w4a5n4@QA"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "name", "QA"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "public", "false"),
					resource.TestCheckNoResourceAttr("apollographql_variant.test", "url"),
					resource.TestCheckResourceAttr("apollographql_variant.test", "graph_id", "Test-w4a5n4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apollographql_variant.test",
				ImportState:       true,
				ImportStateId:     "Test-w4a5n4:QA",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccVariantResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "apollographql_variant" "test" {
  name = "%s"
  graph_id = "Test-w4a5n4"
}
`, name)
}

func testAccVariantResourceConfigNonDefault(name string) string {
	return fmt.Sprintf(`
resource "apollographql_variant" "test" {
  name = "%s"
  public = true
  url = "https://example.com"
  graph_id = "Test-w4a5n4"
}
`, name)
}
