package launchpad

import (
	"fmt"

	common "github.com/Mirantis/mcc/pkg/product/common/api"
	mcc_mke "github.com/Mirantis/mcc/pkg/product/mke"
	mcc_api "github.com/Mirantis/mcc/pkg/product/mke/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	k0s_rig "github.com/k0sproject/rig"
)

// Flattens the config input
func FlattenInputConfigModel(d *schema.ResourceData) (mcc_mke.MKE, error) {
	// Retrieve the metadata
	metadataList := d.Get("metadata").([]interface{})[0]
	m := metadataList.(map[string]interface{})
	metadataName := m["name"].(string)

	// parse spec's cluster flags
	specList := d.Get("spec").([]interface{})[0]
	spec := specList.(map[string]interface{})
	clusterList := spec["cluster"].([]interface{})[0]
	cluster := clusterList.(map[string]interface{})
	prune := cluster["prune"].(bool)

	// parse spec's hosts info from Terraform config
	hosts := mcc_api.Hosts{}
	for _, h := range spec["host"].([]interface{}) {
		host := h.(map[string]interface{})
		role := host["role"].(string)

		extractedHooks := common.Hooks{}
		if len(host["hooks"].([]interface{})) > 0 {
			hooksList := host["hooks"].([]interface{})[0]
			hooks := hooksList.(map[string]interface{})

			_beforeHooks := []string{}
			if len(hooks["before"].([]interface{})) > 0 {
				for _, hook := range hooks["before"].([]interface{}) {
					_beforeHooks = append(_beforeHooks, hook.(string))
				}
			}
			_afterHooks := []string{}
			if len(hooks["after"].([]interface{})) > 0 {
				for _, hook := range hooks["after"].([]interface{}) {
					_afterHooks = append(_afterHooks, hook.(string))
				}
			}
			extractedHooks = common.Hooks{
				"apply": {
					"before": _beforeHooks,
					"after":  _afterHooks,
				},
			}
		}
		connection := k0s_rig.Connection{}
		if len(host["ssh"].([]interface{})) > 0 {
			sshList := host["ssh"].([]interface{})[0]
			ssh := sshList.(map[string]interface{})
			address := ssh["address"].(string)
			key_path := ssh["key_path"].(string)
			user := ssh["user"].(string)
			port := ssh["port"].(int)

			connection = k0s_rig.Connection{
				SSH: &k0s_rig.SSH{
					Address: address,
					KeyPath: key_path,
					User:    user,
					Port:    port,
				},
			}
		} else if len(host["winrm"].([]interface{})) > 0 {
			winrmList := host["winrm"].([]interface{})[0]
			winrm := winrmList.(map[string]interface{})
			address := winrm["address"].(string)
			user := winrm["user"].(string)
			password := winrm["password"].(string)
			port := winrm["port"].(int)
			useHTTPS := winrm["use_https"].(bool)
			insecure := winrm["insecure"].(bool)

			connection = k0s_rig.Connection{
				WinRM: &k0s_rig.WinRM{
					Address:  address,
					Password: password,
					User:     user,
					Port:     port,
					UseHTTPS: useHTTPS,
					Insecure: insecure,
				},
			}
		} else {
			return mcc_mke.MKE{}, fmt.Errorf("%w: %s", ErrMissingConnectionBlock, h)
		}

		extractedHost := &mcc_api.Host{
			Role:        role,
			Connection:  connection,
			Hooks:       extractedHooks,
			MSRMetadata: &mcc_api.MSRMetadata{},
		}
		hosts = append(hosts, extractedHost)
	}

	// parse mcr config
	mcrList := spec["mcr"].([]interface{})[0]
	mcr := mcrList.(map[string]interface{})
	mcrChannel := mcr["channel"].(string)
	mcrInstallURLLinux := mcr["install_url_linux"].(string)
	mcrInstallURLWindows := mcr["install_url_windows"].(string)
	mcrRepoURL := mcr["repo_url"].(string)
	mcrVersion := mcr["version"].(string)

	mcrConfig := common.MCRConfig{
		Version:           mcrVersion,
		InstallURLLinux:   mcrInstallURLLinux,
		InstallURLWindows: mcrInstallURLWindows,
		RepoURL:           mcrRepoURL,
		Channel:           mcrChannel,
	}

	// parse MKE's flags
	mkeList := spec["mke"].([]interface{})[0]
	mke := mkeList.(map[string]interface{})
	mkeAdminUser := mke["admin_username"].(string)
	mkeAdminPass := mke["admin_password"].(string)
	mkeImageRepo := mke["image_repo"].(string)
	mkeVersion := mke["version"].(string)
	// MKE's install flags
	mkeInstallFlags := common.Flags{}
	if len(mke["install_flags"].([]interface{})) > 0 {
		for _, f := range mke["install_flags"].([]interface{}) {
			mkeInstallFlags.Add(f.(string))
		}
	}
	// MKE's upgrade flags
	mkeUpgradeFlags := common.Flags{}
	if len(mke["upgrade_flags"].([]interface{})) > 0 {
		for _, f := range mke["upgrade_flags"].([]interface{}) {
			mkeUpgradeFlags.Add(f.(string))
		}
	}

	mkeConfig := mcc_api.MKEConfig{
		AdminUsername: mkeAdminUser,
		AdminPassword: mkeAdminPass,
		ImageRepo:     mkeImageRepo,
		Version:       mkeVersion,
		InstallFlags:  mkeInstallFlags,
		UpgradeFlags:  mkeUpgradeFlags,

		Metadata: &mcc_api.MKEMetadata{},
	}

	var msrConfig *mcc_api.MSRConfig
	// parse MSR's flags
	if len(spec["msr"].([]interface{})) > 0 {
		tempMSRConfig := mcc_api.MSRConfig{}
		msrList := spec["msr"].([]interface{})[0]
		msr := msrList.(map[string]interface{})
		version := msr["version"].(string)
		image_repo := msr["image_repo"].(string)
		replica_ids := msr["replica_ids"].(string)

		extractedInstallFlags := common.Flags{}
		if len(msr["install_flags"].([]interface{})) > 0 {
			for _, flag := range msr["install_flags"].([]interface{}) {
				extractedInstallFlags.Add(flag.(string))
			}
		}

		tempMSRConfig.Version = version
		tempMSRConfig.ReplicaIDs = replica_ids
		tempMSRConfig.ImageRepo = image_repo
		tempMSRConfig.InstallFlags = extractedInstallFlags
		msrConfig = &tempMSRConfig
	}

	clusterConfig := mcc_api.ClusterConfig{
		APIVersion: "launchpad.mirantis.com/mke/v1.4",
		Kind:       "mke",
		Metadata: &mcc_api.ClusterMeta{
			Name: metadataName,
		},
		Spec: &mcc_api.ClusterSpec{
			Hosts: hosts,
			Cluster: mcc_api.Cluster{
				Prune: prune,
			},
			MKE: mkeConfig,
			MCR: mcrConfig,
			MSR: msrConfig,
		},
	}

	if err := clusterConfig.Validate(); err != nil {
		return mcc_mke.MKE{}, err
	}

	return mcc_mke.MKE{ClusterConfig: clusterConfig}, nil
}
