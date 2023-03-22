
terraform {
  required_version = ">= 1.0.0"
  required_providers {
    mirantis-launchpad = {
      version = ">= 0.9.0"
      source  = "mirantis.com/providers/mirantis-launchpad"
    }
    mirantis-msr-connect = {
      version = ">= 0.9.0"
      source  = "mirantis.com/providers/mirantis-msr-connect"
    }
    mirantis-mke-connect = {
      version = ">= 0.9.0"
      source  = "mirantis.com/providers/mirantis-mke-connect"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.16.0"
    }
  }
}
