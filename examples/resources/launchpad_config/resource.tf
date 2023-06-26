# install Mirantis products using parametrized launchpad
resource "launchpad_config" "example" {
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
      install_flags  = ["--default-orchestrator"]
    }
    msr {
      version = "2.9.4"
    }

    host {
      role = "manager"
      ssh {
        address  = "manager1.example.org"
        key_path = "./key.pem"
        user     = "ubuntu"
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
