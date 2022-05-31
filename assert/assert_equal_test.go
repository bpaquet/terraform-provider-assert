package assert

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResource_empty(t *testing.T) {
	dsn := "data.assert_equal.my_first_test"
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_empty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
				),
			},
		},
	})
}

const testAccResourceConfig_empty = `
data "assert_equal" "my_first_test" {
	current = []
	expected = []
}
`

func TestAccResource_equal(t *testing.T) {
	dsn := "data.assert_equal.my_first_test"
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_equal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
				),
			},
		},
	})
}

const testAccResourceConfig_equal = `
data "assert_equal" "my_first_test" {
	current = [1, 2, 3]
	expected = [1, 2, 3]
}
`

func TestAccResource_not_equal(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceConfig_not_equal,
				ExpectError: regexp.MustCompile("elements differ"),
			},
		},
	})
}

const testAccResourceConfig_not_equal = `
data "assert_equal" "my_first_test" {
	current = [1, 2, 3]
	expected = [1, 2]
}
`

func TestAccResource_not_equal_no_fail(t *testing.T) {
	dsn := "data.assert_equal.my_first_test"
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_not_equal_no_fail,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
				),
			},
		},
	})
}

const testAccResourceConfig_not_equal_no_fail = `
provider "assert" {
	fail_on_assert = false
}
data "assert_equal" "my_first_test" {
	current = [1, 2, 3]
	expected = [1, 2]
}
`

func TestAccResource_not_equal_with_error_message(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceConfig_not_equal_with_error_message,
				ExpectError: regexp.MustCompile("This is an incredible error mesaage"),
			},
		},
	})
}

const testAccResourceConfig_not_equal_with_error_message = `
data "assert_equal" "my_first_test" {
	current = [1, 2, 3]
	expected = [1, 2]
	message = "This is an incredible error mesaage"
}
`
