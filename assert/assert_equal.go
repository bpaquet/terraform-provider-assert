package assert

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func assertEqualDataSource() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`assert_equal`" + ` resource compares the two lists provided as arugments, and fail during apply if they are not equal.`,

		Read: resourceRead,

		Schema: map[string]*schema.Schema{
			"current": {
				Description: "A list describing the current state. Example: a list of instance ids frem a data source.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"expected": {
				Description: "A list describing the expected state. Example: a list of instance ids extracted from resources declared in the same workspace.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

type handler struct {
	testing.T
	result error
}

func (handler *handler) Errorf(format string, args ...any) {
	handler.result = fmt.Errorf(format, args...)
}

func resourceRead(d *schema.ResourceData, meta any) error {
	d.SetId(fmt.Sprintf("%d", rand.Int()))

	t := &handler{}
	assert.ElementsMatch(t, d.Get("current"), d.Get("expected"))
	return t.result
}
