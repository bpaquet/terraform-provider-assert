package assert

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
)

func assertEqualResource() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`assert_equal`" + ` resource compares the two lists provided as arugments, and fail during apply if they are not equal.`,

		Create: resourceCreate,
		Read:   resourceRead,
		Delete: resourceDelete,

		Schema: map[string]*schema.Schema{
			"current": {
				Description: "A list describing the current state. Example: a list of instance ids frem a data source.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"expected": {
				Description: "A list describing the expected state. Example: a list of instance ids extracted from resources declared in the same workspace.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return nil
}

func resourceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
