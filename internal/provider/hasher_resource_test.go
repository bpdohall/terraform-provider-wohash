package provider

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestHasherResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccHasherInputWO("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"wohash_hasher.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact(hashInput("one")),
					),
					statecheck.ExpectKnownValue(
						"wohash_hasher.test",
						tfjsonpath.New("output"),
						knownvalue.StringExact(hashInput("one")),
					),
				},
			},
			// Update and Read testing
			{
				Config: testAccHasherInputWO("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"wohash_hasher.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact(hashInput("two")),
					),
					statecheck.ExpectKnownValue(
						"wohash_hasher.test",
						tfjsonpath.New("output"),
						knownvalue.StringExact(hashInput("two")),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccHasherInputWO("two"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccHasherInputWO("one"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccHasherWithEphemeral(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccHasherEphemeralInput(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"wohash_hasher.testephem",
						tfjsonpath.New("id"),
						knownvalue.StringExact(hashInput("testvalue")),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigVariables: map[string]config.Variable{
					"ephem": config.StringVariable("testvalue"),
				},
			},
			{
				Config: testAccHasherEphemeralInput(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"wohash_hasher.testephem",
						tfjsonpath.New("id"),
						knownvalue.StringExact(hashInput("testvalue")),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigVariables: map[string]config.Variable{
					"ephem": config.StringVariable("testvalue"),
				},
			},
			{
				Config: testAccHasherEphemeralInput(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"wohash_hasher.testephem",
						tfjsonpath.New("id"),
						knownvalue.StringExact(hashInput("newvalue")),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				ConfigVariables: map[string]config.Variable{
					"ephem": config.StringVariable("newvalue"),
				},
			},
		},
	})
}

func testAccHasherInputWO(value string) string {
	return fmt.Sprintf(`
resource "wohash_hasher" "test" {
  input_wo = %[1]q
}
`, value)
}

func testAccHasherEphemeralInput() string {
	return `
variable "ephem" {
  type = string
  ephemeral = true
}

resource "wohash_hasher" "testephem" {
  input_wo = var.ephem
}

output "hash_of_input" {
  value = wohash_hasher.testephem.output
}
`
}
