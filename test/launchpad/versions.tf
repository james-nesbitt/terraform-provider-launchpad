
terraform {
  required_version = ">= 1.0.0"
  required_providers {
    launchpad = {
      source = "Mirantis/launchpad"
    }
    mirantis-msr-connect = {
      source = "Mirantis/msr"
    }
    mirantis-mke-connect = {
      source = "Mirantis/mke"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.16.0"
    }
  }
}
