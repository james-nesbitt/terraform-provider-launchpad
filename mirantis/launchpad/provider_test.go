package launchpad_test

import (
	"testing"

	launchpad "github.com/Mirantis/terraform-provider-launchpad/mirantis/launchpad"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := launchpad.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = launchpad.Provider()
}
