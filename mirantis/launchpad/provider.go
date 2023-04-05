package launchpad

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"mirantis-launchpad_launchpad": ResourceConfig(),
			"mirantis-launchpad_yaml":      ResourceYamlConfig(),
		},
	}
}
