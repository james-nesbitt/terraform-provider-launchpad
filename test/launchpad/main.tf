# This aws secrets manager secret will be used in monitoring-infra-mke-apps terraform configuration as well
// provision cluster machines
module "provision" {
  source  = "terraform-mirantis-modules/launchpad-aws/mirantis"
  version = "0.0.1"

  aws_region   = "us-west-2"
  master_count = 1
  worker_count = 3
  cluster_name = var.cluster_name
}

// Mirantis installing terraform provider
provider "launchpad" {}

// launchpad install from provisioned cluster
resource "launchpad_config" "cluster" {
  skip_destroy = true

  metadata {
    name = var.cluster_name
  }
  spec {
    cluster {
      prune = true
    }

    dynamic "host" {
      for_each = module.provision.hosts

      content {
        role = host.value.role

        dynamic "ssh" {
          for_each = can(host.value.ssh) ? [1] : []

          content {
            address  = host.value.ssh.address
            user     = host.value.ssh.user
            key_path = host.value.ssh.keyPath
            port     = 22
          }
        }

        dynamic "winrm" {
          for_each = can(host.value.winRM) ? [1] : []

          content {
            address   = host.value.winRM.address
            port      = 5985
            user      = host.value.winRM.user
            password  = host.value.winRM.password
            use_https = host.value.winRM.useHTTPS
            insecure  = host.value.winRM.insecure
          }
        }

      }
    }

    mcr {
      channel             = "stable"
      install_url_linux   = "https://get.mirantis.com/"
      install_url_windows = "https://get.mirantis.com/install.ps1"
      repo_url            = "https://repos.mirantis.com"
      version             = var.mcr_version
    } // mcr

    mke {
      admin_password = var.admin_password
      admin_username = var.admin_username
      image_repo     = "docker.io/mirantis"
      version        = var.mke_version
      install_flags  = ["--san=${module.provision.mke_lb}", "--default-node-orchestrator=kubernetes", "--nodeport-range=32768-35535"]
      upgrade_flags  = ["--force-recent-backup", "--force-minimums"]
    } // mke
  }   // spec
}
