package launchpad

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"launchpad_config":      ResourceConfig(),
			"launchpad_yaml_config": ResourceYamlConfig(),
		},
	}
}
