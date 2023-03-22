# MCC Terraform Provider

This is a doc which will explain and demonstrate what MCC terraform provider is and how to use it.



- [MCC Terraform Provider](#mcc-terraform-provider)
  - [Introduction](#introduction)
  - [Schema](#schema)
  - [Usage](#usage)
  - [Examples](#examples)
    - [Multiple staticly defined *hosts* attributes:](#multiple-staticly-defined-hosts-attributes)
    - [Dynamic blocks example:](#dynamic-blocks-example)
    - [Full MCC terraform provider example](#full-mcc-terraform-provider-example)
  - [References](#references)


## Introduction

The MCC terraform provider directly imports the Launchpad code, therefore the terraform provider is always working with the latest Launchpad api.

## Schema

This section will explain the terraform schema for the MCC terraform provider.

- `spec`: A terraform block containing *cluster* and *host* attributes.

  - `cluster`: Containing general cluster flags - like *prune*
    - `prune`: Removes hosts(nodes) which are not part of the spec

  - `host`: Each is a node which will be in the MKE cluster. Each of those hosts should have *role* and *connection* block(either *ssh* or *winrm*)
    -  `ssh` or `winrm`: Method for Launchpad to connect to the remote hosts. Depending on the host type(Linux vs Windows). The below attributes are required for the remote conection.
      - `ssh`: Connection block for Linux hosts. Contains host *ip*, pem file *key_path*, ssh *user*
      - `winrm`: Connection block for Windows hosts. Contains host ip   *address*, *user*, *password*, *port*.
      - `hooks`: These are commands(terminal commands) which can be executed *before* or *after* the provisioning of a host
    - `role`: The role of the host, so that Launchpad knows where to install the appropriate software (MKE, MSR). These roles can be one of the following: *manager*, *worker*, *msr*.

  - `mcr`: The terraform block for the MCR product containing all the required attributes for the installation.
    - `channel`: The type of engine channel to use.
    - `install_url_linux`: The engine install script location for linux
    - `install_url_windows`: The engine install script location for windows
    - `image_repo`: The engine repo where to pull the install script from
    - `version`: The engine version to install

  - `mke`: The terraform block for the MKE product containing all the required attributes for the installation.
    - `admin_username`: MKE's admin username
    - `admin_password`:  MKE's admin password
    - `version`: MKE version to install
    - `image_repo`: Where to pull the MKE installation images
    - `install_flags`: The MKE flags that you can set for the MKE installation, i.e. san, orchestrator, etc.
    - `upgrade_flags`: Upgrade flags which are used on performing MKE upgrade

  - `msr`: The optional terraform block for the MSR product containing all the required attributes for the installation.
    - `image_repo`: The repository where to pull the MSR installation images
    - `version`: Which MSR version to be installed
    - `replica_ids`: Used to identify the *type* of assigning that MSR does on its hosts. MSR finds the highest replica id and assigns sequential ones starting from that to all the hosts without replica ids.
    - `install_flags`: A list of the installation flags used when performing MSR installation.

## Usage
This section will show you how to import the MCC terraform provider

Since the terraform provider is not released on the Hashicorp registry, we need to point to its source locally. All local terraform providers/plugins are installed under *~/.terraform.d/plugins/*(however we don't need to pre-pend this path). The version is also important since you can have multiple local versions installed. Simply, point to the companie's directory/providers/"name of provider". The example below demonstrates this.
```
terraform {
  required_providers {
    mirantis-launchpad = {
      version = "= 0.9.0"
      source  = "mirantis.com/providers/mirantis-launchpad"
    }
  }
}
```
Afterwards, you just need to import the terraform provider with the following code:
```
provider "mirantis-launchpad" {}
```
You now have the access to the provider **resources/data sources**.

### Launchpad installer resources

You can use launchpad in install MCR/MKE/MSR-classic by using the `mirantis_install_launchpad` resource.
This is a standalone resource that does not require that launchpad is installed.

```
// launchpad install simple (static hosts)
resource "mirantis-launchpad_launchpad" "my-cluster" {

  metadata {
    name = "my-cluster-name"
  }
  spec {
    cluster {
      prune = true
    }

    host {
      role = "manager"
      ssh {
        address = "52.35.136.67"
        user = "ubuntu"
        key_path = "./systest-cluster.pem"
      }
    }
    host {
      role = "worker"
      hooks {
        after = ["ls -la"]
      }
      ssh {
        address = "34.222.18.107"
        user = "ubuntu"
        key_path = "./systest-cluster.pem"
      }
    }

    mcr {
      channel = "stable"
      install_url_linux = "https://get.mirantis.com/"
      install_url_windows = "https://get.mirantis.com/install.ps1"
      image_repo = "https://repos.mirantis.com"
      version = "20.10.9"
    } // mcr

    mke {
      admin_password = "password"
      admin_username = "admin"
      image_repo = "docker.io/mirantis"
      version = var.mke_version
      install_flags = ["--san=${module.elb_mke.lb_dns_name}", "--default-node-orchestrator=kubernetes", "--nodeport-range=32768-35535"]
      upgrade_flags = ["--force-recent-backup", "--force-minimums"]
    } // mke

    msr {
      image_repo = "docker.io/mirantis"
      version = "2.8.6"
      replica_ids = "sequential"
      install_flags = local.msr_install_flags
    } // msr

  } // spec
}
```

Note that static host lists can be inflexible, but terraform offers a feature
for dynamically writing such blocks from data structures such as lists:

```
   dynamic  "host" {
      for_each = local.hosts

      content {
        role = hosts.value.role
        ssh  {
          address = hosts.value.instance.public_ip
          user = hosts.value.ssh.user
          key_path = hosts.value.ssh.keyPath
        }
      }
    }
```

## References

A link to the private MCC github repo: https://github.com/Mirantis/mcc
A link to the public Mirantis Launchpad github repo: https://github.com/Mirantis/launchpad
