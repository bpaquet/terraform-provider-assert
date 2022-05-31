package assert

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
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

func TestAccResource_sns_wrong_arn(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceConfig_sns_wrong_arn,
				ExpectError: regexp.MustCompile("Unable to parse ARN"),
			},
		},
	})
}

const testAccResourceConfig_sns_wrong_arn = `
provider "assert" {
	sns_topic_arn = "foo"
}

data "assert_equal" "my_first_test" {
	current = [1, 2, 3]
	expected = [1, 2]
}
`

type mockPublishApi struct {
	subject *string
	message *string
}

func (r *mockPublishApi) PublishMessage(ctx context.Context, subject, message *string) error {
	r.subject = subject
	r.message = message
	return nil
}

func TestAccResource_sns_send_message(t *testing.T) {
	api := &mockPublishApi{
		subject: nil,
		message: nil,
	}
	var mock PublishApi = api
	MockedPublishApi = &mock
	dsn := "data.assert_equal.my_first_test"
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_sns_send_message,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
				),
			},
		},
	})
	assert.Equal(t, *api.subject, "Error: This is a message")
	assert.Equal(t, *api.message, "Error: This is a message\nCurrent: [1 2 3]\nExpected: [1 2]")
}

const testAccResourceConfig_sns_send_message = `
provider "assert" {
	sns_topic_arn = "arn:aws:sns:eu-west-1:11111:terraform-github-notifications"
	fail_on_assert = false
}

data "assert_equal" "my_first_test" {
	current = [1, 2, 3]
	expected = [1, 2]
	message = "This is a message"
}
`
