package assert

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func assertEqualDataSource() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`assert_equal`" + ` resource compares the two lists provided as arguments, and fails if they are not equal.`,

		ReadContext: resourceRead,

		Schema: map[string]*schema.Schema{
			"current": {
				Description: "A list describing the current state. Example: a list of instance ids from a data source.",
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

			"message": {
				Description: "The error message if the expected != current",
				Type:        schema.TypeString,
				Optional:    true,
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

type templateParams struct {
	Message  string
	Expected string
	Current  string
}

func executeTemplate(tmpl *template.Template, params *templateParams) (*string, error) {
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, params); err != nil {
		return nil, err
	}
	result := doc.String()
	return &result, nil
}

func resourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%d", rand.Int()))

	t := &handler{}
	assert.ElementsMatch(t, d.Get("expected"), d.Get("current"), d.Get("message"))
	providerConfig := m.(*providerConfig)
	if providerConfig.publishApi != nil && t.result != nil {
		params := &templateParams{
			Message:  d.Get("message").(string),
			Expected: fmt.Sprintf("%v", d.Get("expected")),
			Current:  fmt.Sprintf("%v", d.Get("current")),
		}
		body, err := executeTemplate(providerConfig.snsBodyTemplate, params)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, fmt.Sprintf("Sending SNS notification, body %s", *body))
		subject, err := executeTemplate(providerConfig.snsSubjectTemplate, params)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, fmt.Sprintf("Sending SNS notification, subject %s", *subject))
		err = providerConfig.publishApi.PublishMessage(ctx, subject, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if providerConfig.failOnAssert {
		return diag.FromErr(t.result)
	} else {
		return nil
	}
}
