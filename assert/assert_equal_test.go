package assert

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResource_basic(t *testing.T) {
	dsn := "data.assert_equal.my_first_test"
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
				),
			},
		},
	})
}

const testAccResourceConfig_basic = `
data "assert_equal" "my_first_test" {
	current = []
	expected = []
}
`
