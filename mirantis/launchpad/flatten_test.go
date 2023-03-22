package launchpad_test

import (
	"reflect"
	"testing"

	common "github.com/Mirantis/mcc/pkg/product/common/api"
	mcc_mke "github.com/Mirantis/mcc/pkg/product/mke"
	mcc_api "github.com/Mirantis/mcc/pkg/product/mke/api"
	"github.com/Mirantis/terraform-provider-launchpad/mirantis/launchpad"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	k0s_rig "github.com/k0sproject/rig"
)

var DUMMYSTR = "test"
var DUMMYROLE = "manager"
var DUMMYPORT = 22

func dummyConn(isSSH bool) k0s_rig.Connection {
	if isSSH {
		return k0s_rig.Connection{
			SSH: &k0s_rig.SSH{
				Address: DUMMYSTR,
				KeyPath: DUMMYSTR,
				User:    DUMMYSTR,
				Port:    DUMMYPORT,
			},
		}
	} else {
		return k0s_rig.Connection{
			WinRM: &k0s_rig.WinRM{
				Address:  DUMMYSTR,
				Password: DUMMYSTR,
				User:     DUMMYSTR,
				Port:     DUMMYPORT,
				UseHTTPS: true,
				Insecure: true,
			},
		}
	}
}

func dummyMKEObj(conn k0s_rig.Connection, hasMSR bool) mcc_mke.MKE {
	mke := mcc_mke.MKE{}
	eBeforeHooks := []string{
		DUMMYSTR,
	}
	eAfterHooks := []string{
		DUMMYSTR,
	}
	eHooks := common.Hooks{
		"apply": {
			"before": eBeforeHooks,
			"after":  eAfterHooks,
		},
	}

	eHosts := mcc_api.Hosts{}
	eHosts = append(eHosts, &mcc_api.Host{
		Role:        DUMMYROLE,
		Connection:  conn,
		Hooks:       eHooks,
		MSRMetadata: &mcc_api.MSRMetadata{},
	})

	mkeIpgradeFlags := common.Flags{DUMMYSTR}
	mkeUpgradeFlags := common.Flags{DUMMYSTR}
	mkeConfig := mcc_api.MKEConfig{
		AdminUsername: DUMMYSTR,
		AdminPassword: DUMMYSTR,
		ImageRepo:     DUMMYSTR,
		Version:       DUMMYSTR,
		InstallFlags:  mkeIpgradeFlags,
		UpgradeFlags:  mkeUpgradeFlags,

		Metadata: &mcc_api.MKEMetadata{},
	}

	mcrConfig := common.MCRConfig{
		Version:           DUMMYSTR,
		InstallURLLinux:   DUMMYSTR,
		InstallURLWindows: DUMMYSTR,
		RepoURL:           DUMMYSTR,
		Channel:           DUMMYSTR,
	}

	var msrConfig *mcc_api.MSRConfig
	if hasMSR {
		tempMSRConfig := mcc_api.MSRConfig{}
		tempMSRConfig.ImageRepo = DUMMYSTR
		tempMSRConfig.InstallFlags = common.Flags{DUMMYSTR}
		tempMSRConfig.Version = DUMMYSTR
		tempMSRConfig.ReplicaIDs = DUMMYSTR

		msrConfig = &tempMSRConfig
	}

	mke.ClusterConfig = mcc_api.ClusterConfig{
		APIVersion: "launchpad.mirantis.com/mke/v1.4",
		Kind:       "mke",
		Metadata: &mcc_api.ClusterMeta{
			Name: DUMMYSTR,
		},
		Spec: &mcc_api.ClusterSpec{
			Hosts: eHosts,
			Cluster: mcc_api.Cluster{
				Prune: true,
			},
			MKE: mkeConfig,
			MCR: mcrConfig,
			MSR: msrConfig,
		},
	}

	return mke
}

func TestFlattenInputConfigModelValidSSH(t *testing.T) {
	c := struct {
		input    *schema.ResourceData
		expected mcc_mke.MKE
	}{
		input:    launchpad.ResourceConfig().TestResourceData(),
		expected: dummyMKEObj(dummyConn(true), false),
	}

	metadata := make(map[string]interface{})
	metadata["name"] = DUMMYSTR

	if err := c.input.Set("metadata", []interface{}{metadata}); err != nil {
		t.Fatalf("Error setting schema: '%#v'. Error: %s", metadata, err)
	}

	spec := make(map[string]interface{})
	clusterList := make([]map[string]interface{}, 1)

	cluster := map[string]interface{}{
		"prune": true,
	}
	clusterList[0] = cluster
	spec["cluster"] = clusterList

	hostsList := make([]map[string]interface{}, 1)
	hooksList := make([]map[string]interface{}, 1)
	hooks := map[string]interface{}{
		"before": []string{
			DUMMYSTR,
		},
		"after": []string{
			DUMMYSTR,
		},
	}
	hooksList[0] = hooks

	sshList := make([]map[string]interface{}, 1)
	ssh := map[string]interface{}{
		"address":  DUMMYSTR,
		"key_path": DUMMYSTR,
		"user":     DUMMYSTR,
		"port":     DUMMYPORT,
	}
	sshList[0] = ssh

	host := map[string]interface{}{
		"role":  DUMMYROLE,
		"hooks": hooksList,
		"ssh":   sshList,
	}

	hostsList[0] = host
	spec["host"] = hostsList

	mcrList := make([]map[string]interface{}, 1)
	mcr := map[string]interface{}{
		"channel":             DUMMYSTR,
		"install_url_linux":   DUMMYSTR,
		"install_url_windows": DUMMYSTR,
		"repo_url":            DUMMYSTR,
		"version":             DUMMYSTR,
	}

	mcrList[0] = mcr
	spec["mcr"] = mcrList

	mkeList := make([]map[string]interface{}, 1)
	mkeIFlag := []string{
		DUMMYSTR,
	}
	mkeUFlag := []string{
		DUMMYSTR,
	}
	mke := map[string]interface{}{
		"admin_password": DUMMYSTR,
		"admin_username": DUMMYSTR,
		"image_repo":     DUMMYSTR,
		"version":        DUMMYSTR,
		"install_flags":  mkeIFlag,
		"upgrade_flags":  mkeUFlag,
	}

	mkeList[0] = mke
	spec["mke"] = mkeList

	if err := c.input.Set("spec", []interface{}{spec}); err != nil {
		t.Fatalf("Error setting schema: '%#v'. Error %s", metadata, err)
	}

	mkeClient, err := launchpad.FlattenInputConfigModel(c.input)

	if err != nil {
		t.Errorf("unexpected error: (%v)", err)
	}

	if !reflect.DeepEqual(mkeClient, c.expected) {
		t.Fatalf("Error matching output and expected: %#v vs %#v", mkeClient, c.expected)
	}
}

func TestFlattenInputConfigModelValidWINRM(t *testing.T) {
	c := struct {
		input    *schema.ResourceData
		expected mcc_mke.MKE
	}{
		input:    launchpad.ResourceConfig().TestResourceData(),
		expected: dummyMKEObj(dummyConn(false), false),
	}

	metadata := make(map[string]interface{})
	metadata["name"] = DUMMYSTR

	if err := c.input.Set("metadata", []interface{}{metadata}); err != nil {
		t.Fatalf("Error setting schema: '%#v'. Error: %s", metadata, err)
	}

	spec := make(map[string]interface{})
	clusterList := make([]map[string]interface{}, 1)

	cluster := map[string]interface{}{
		"prune": true,
	}
	clusterList[0] = cluster
	spec["cluster"] = clusterList

	hostsList := make([]map[string]interface{}, 1)
	hooksList := make([]map[string]interface{}, 1)
	hooks := map[string]interface{}{
		"before": []string{
			DUMMYSTR,
		},
		"after": []string{
			DUMMYSTR,
		},
	}
	hooksList[0] = hooks

	winrmList := make([]map[string]interface{}, 1)
	winrm := map[string]interface{}{
		"address":   DUMMYSTR,
		"user":      DUMMYSTR,
		"password":  DUMMYSTR,
		"port":      DUMMYPORT,
		"use_https": true,
		"insecure":  true,
	}
	winrmList[0] = winrm

	host := map[string]interface{}{
		"role":  DUMMYROLE,
		"hooks": hooksList,
		"winrm": winrmList,
	}

	hostsList[0] = host
	spec["host"] = hostsList

	mcrList := make([]map[string]interface{}, 1)
	mcr := map[string]interface{}{
		"channel":             DUMMYSTR,
		"install_url_linux":   DUMMYSTR,
		"install_url_windows": DUMMYSTR,
		"repo_url":            DUMMYSTR,
		"version":             DUMMYSTR,
	}

	mcrList[0] = mcr
	spec["mcr"] = mcrList

	mkeList := make([]map[string]interface{}, 1)
	mkeIFlag := []string{
		DUMMYSTR,
	}
	mkeUFlag := []string{
		DUMMYSTR,
	}
	mke := map[string]interface{}{
		"admin_password": DUMMYSTR,
		"admin_username": DUMMYSTR,
		"image_repo":     DUMMYSTR,
		"version":        DUMMYSTR,
		"install_flags":  mkeIFlag,
		"upgrade_flags":  mkeUFlag,
	}

	mkeList[0] = mke
	spec["mke"] = mkeList

	if err := c.input.Set("spec", []interface{}{spec}); err != nil {
		t.Fatalf("Error setting schema. Error %s", err)
	}

	mkeClient, err := launchpad.FlattenInputConfigModel(c.input)

	if err != nil {
		t.Errorf("unexpected error: (%v)", err)
	}

	if !reflect.DeepEqual(mkeClient, c.expected) {
		t.Fatalf("Error matching output and expected: %#v vs %#v", mkeClient, c.expected)
	}
}

func TestFlattenInputConfigModelValidMSR(t *testing.T) {
	c := struct {
		input    *schema.ResourceData
		expected mcc_mke.MKE
	}{
		input:    launchpad.ResourceConfig().TestResourceData(),
		expected: dummyMKEObj(dummyConn(false), true),
	}

	metadata := make(map[string]interface{})
	metadata["name"] = DUMMYSTR

	if err := c.input.Set("metadata", []interface{}{metadata}); err != nil {
		t.Fatalf("Error setting schema: '%#v'. Error: %s", metadata, err)
	}

	spec := make(map[string]interface{})
	clusterList := make([]map[string]interface{}, 1)

	cluster := map[string]interface{}{
		"prune": true,
	}
	clusterList[0] = cluster
	spec["cluster"] = clusterList

	hostsList := make([]map[string]interface{}, 1)
	hooksList := make([]map[string]interface{}, 1)
	hooks := map[string]interface{}{
		"before": []string{
			DUMMYSTR,
		},
		"after": []string{
			DUMMYSTR,
		},
	}
	hooksList[0] = hooks

	winrmList := make([]map[string]interface{}, 1)
	winrm := map[string]interface{}{
		"address":   DUMMYSTR,
		"user":      DUMMYSTR,
		"password":  DUMMYSTR,
		"port":      DUMMYPORT,
		"use_https": true,
		"insecure":  true,
	}
	winrmList[0] = winrm

	host := map[string]interface{}{
		"role":  DUMMYROLE,
		"hooks": hooksList,
		"winrm": winrmList,
	}

	hostsList[0] = host
	spec["host"] = hostsList

	mcrList := make([]map[string]interface{}, 1)
	mcr := map[string]interface{}{
		"channel":             DUMMYSTR,
		"install_url_linux":   DUMMYSTR,
		"install_url_windows": DUMMYSTR,
		"repo_url":            DUMMYSTR,
		"version":             DUMMYSTR,
	}

	mcrList[0] = mcr
	spec["mcr"] = mcrList

	mkeList := make([]map[string]interface{}, 1)
	mkeIFlag := []string{
		DUMMYSTR,
	}
	mkeUFlag := []string{
		DUMMYSTR,
	}
	mke := map[string]interface{}{
		"admin_password": DUMMYSTR,
		"admin_username": DUMMYSTR,
		"image_repo":     DUMMYSTR,
		"version":        DUMMYSTR,
		"install_flags":  mkeIFlag,
		"upgrade_flags":  mkeUFlag,
	}

	mkeList[0] = mke
	spec["mke"] = mkeList

	msrList := make([]map[string]interface{}, 1)
	msrIFlag := []string{
		DUMMYSTR,
	}

	msr := map[string]interface{}{
		"version":       DUMMYSTR,
		"replica_ids":   DUMMYSTR,
		"image_repo":    DUMMYSTR,
		"install_flags": msrIFlag,
	}

	msrList[0] = msr
	spec["msr"] = msrList

	if err := c.input.Set("spec", []interface{}{spec}); err != nil {
		t.Fatalf("Error setting schema. Error %s", err)
	}

	mkeClient, err := launchpad.FlattenInputConfigModel(c.input)

	if err != nil {
		t.Errorf("unexpected error: (%v)", err)
	}

	if !reflect.DeepEqual(mkeClient, c.expected) {
		t.Fatalf("Error matching output and expected: %#v vs %#v", mkeClient, c.expected)
	}
}
