package assert

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sns_topic_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SNS topic ARN to be notified in case of error",
			},
			"sns_body_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Error: {{.Message}}\nCurrent: {{.Current}}\nExpected: {{.Expected}}",
				Description: "Template to build the body of the message for SNS notifications",
			},
			"sns_subject_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Error: {{.Message}}",
				Description: "Template to build the subject of the message for SNS notifications",
			},
			"fail_on_assert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, any assertion failed will through an error in the provider",
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"assert_equal": assertEqualDataSource(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return p
}

type PublishApi interface {
	PublishMessage(ctx context.Context, subject, message *string) error
}

var MockedPublishApi *PublishApi = nil

type realPublishApi struct {
	snsTopicARN string
	client      *sns.Client
}

func (r *realPublishApi) PublishMessage(ctx context.Context, subject, message *string) error {
	input := &sns.PublishInput{
		Message:  message,
		Subject:  subject,
		TopicArn: &r.snsTopicARN,
	}
	_, err := r.client.Publish(ctx, input)
	return err
}

type providerConfig struct {
	snsBodyTemplate    *template.Template
	snsSubjectTemplate *template.Template
	failOnAssert       bool
	publishApi         PublishApi
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	snsTopicARN := d.Get("sns_topic_arn").(string)
	snsBodyTemplate, err := template.New("sns_body_template").Parse(d.Get("sns_body_template").(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	snsSubjectTemplate, err := template.New("sns_subject_template").Parse(d.Get("sns_subject_template").(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	failOnAssert := d.Get("fail_on_assert").(bool)
	providerConfig := &providerConfig{
		snsBodyTemplate:    snsBodyTemplate,
		snsSubjectTemplate: snsSubjectTemplate,
		failOnAssert:       failOnAssert,
	}
	if snsTopicARN != "" {
		snsTopicARNSplitted := strings.Split(snsTopicARN, ":")
		if len(snsTopicARNSplitted) != 6 {
			return nil, diag.FromErr(fmt.Errorf("Unable to parse ARN: %s", snsTopicARN))
		}
		region := snsTopicARNSplitted[3]

		if MockedPublishApi == nil {
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
			if err != nil {
				return nil, diag.FromErr(err)
			}
			client := sns.NewFromConfig(cfg)
			providerConfig.publishApi = &realPublishApi{
				snsTopicARN: snsTopicARN,
				client:      client,
			}
		} else {
			providerConfig.publishApi = *MockedPublishApi
		}
	}
	return providerConfig, nil
}
