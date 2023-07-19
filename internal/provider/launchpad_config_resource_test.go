package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLaunchpadConfigResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccLaunchpadConfigResourceConfig_minimal(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("launchpad_config.test", "skip_destroy", "false"),
					resource.TestCheckResourceAttr("launchpad_config.test", "spec.host.0.hooks.0.apply.0.before.0", "ls -la"),
				),
			},
		},
	})
}

func testAccLaunchpadConfigResourceConfig_minimal() string {
	return `
resource "launchpad_config" "test" {
    metadata {
        name = "test"
    }
    spec {
        mcr {
            version = "22.10"
        }
        mke {
            version        = "3.6.4"
            admin_password = "mypassword"
            install_flags  = ["--flag1", "--flag2" ]
        }
        msr {
            version = "2.9.1"
        }

        host {
            role = "manager"
            ssh {
                address  = "manager1.example.org"
                key_path = "./key.pem"
                user     = "ubuntu"
            }

            hooks {
                apply {
                    before = [ "ls -la", "pwd" ]
                }
            }

            mcr_config {
                debug = true
                bip = "172.20.0.1/16"

                default_address_pools = [
                    {
                        base="172.21.0.0",
                        size=16
                        test="test" // this should produce an error but it doesn't
                    },
                    {
                        base="172.22.0.0",
                        size=16
                    }
                ]
            }
        }

        host {
            role = "worker"
            ssh {
                address  = "worker1.example.org"
                key_path = "./key.pem"
                user     = "ubuntu"
            }
        }

        host {
            role = "worker"
            winrm {
                address  = "windowsworker1.example.org"
                user     = "ubuntu"
                password = "my-win-password"
            }
        }

        host {
            role = "msr"
            ssh {
                address  = "msr1.example.org"
                key_path = "./key.pem"
                user     = "ubuntu"
            }
        }
    }
}
`
}
