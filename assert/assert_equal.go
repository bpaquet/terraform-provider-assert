package assert

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
)

func assertEqualDataSource() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`assert_equal`" + ` resource compares the two lists provided as arugments, and fail during apply if they are not equal.`,

		Read: resourceRead,

		Schema: map[string]*schema.Schema{
			"current": {
				Description: "A list describing the current state. Example: a list of instance ids frem a data source.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"expected": {
				Description: "A list describing the expected state. Example: a list of instance ids extracted from resources declared in the same workspace.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return nil
}
