package provider

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	mcc_common_api "github.com/Mirantis/mcc/pkg/product/common/api"
	mcc_mke_api "github.com/Mirantis/mcc/pkg/product/mke/api"
	k0s_dig "github.com/k0sproject/dig"
	k0s_rig "github.com/k0sproject/rig"
)

const (
	HostRoleMSR = "msr"
)

func launchpadSchema14() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Mirantis installation using launchpad, parametrized",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"skip_destroy": schema.BoolAttribute{
				MarkdownDescription: "Do not bother uninstalling on destroy",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},

		Blocks: map[string]schema.Block{

			"metadata": schema.SingleNestedBlock{
				MarkdownDescription: "Metadata for the launchpad cluster",

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Cluster name",
						Required:            true,
					},
				},
			},

			"spec": schema.SingleNestedBlock{
				MarkdownDescription: "Launchpad install specifications",

				Blocks: map[string]schema.Block{

					"cluster": schema.ListNestedBlock{
						MarkdownDescription: "MSR installation configuration",

						Validators: []validator.List{
							listvalidator.SizeAtMost(1),
						},
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"prune": schema.BoolAttribute{
									MarkdownDescription: "Prune cluster resources orphaned on apply",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
								},
							},
						},
					},

					"mcr": schema.SingleNestedBlock{
						MarkdownDescription: "MCR installation configuration",

						Attributes: map[string]schema.Attribute{
							"version": schema.StringAttribute{
								MarkdownDescription: "MCR version to install",
								Required:            true,
							},
							"channel": schema.StringAttribute{
								MarkdownDescription: "Repitory installation channel",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("stable"),
							},
							"repo_url": schema.StringAttribute{
								MarkdownDescription: "Repository installation URL for installation script",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("https://repos.mirantis.com"),
							},
							"install_url_linux": schema.StringAttribute{
								MarkdownDescription: "MCR installation script for linux installations",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("https://get.mirantis.com/"),
							},
							"install_url_windows": schema.StringAttribute{
								MarkdownDescription: "MCR installation script for windows installations",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("https://get.mirantis.com/install.ps1"),
							},
						},
					},

					"mke": schema.SingleNestedBlock{
						MarkdownDescription: "MKE installation configuration",

						Attributes: map[string]schema.Attribute{
							"version": schema.StringAttribute{
								MarkdownDescription: "MKE version to install",
								Required:            true,
							},
							"image_repo": schema.StringAttribute{
								MarkdownDescription: "Image repo for MKE images",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("docker.io/mirantis"),
							},
							"admin_username": schema.StringAttribute{
								MarkdownDescription: "MKE admin user name",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("admin"),
							},
							"admin_password": schema.StringAttribute{
								MarkdownDescription: "MKE admin user password",
								Required:            true,
								Sensitive:           true,
							},
							"license_file_path": schema.StringAttribute{
								MarkdownDescription: "MKE license file path",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString(""),
							},

							"install_flags": schema.ListAttribute{
								MarkdownDescription: "Optional MKE bootstrapper install flags",
								ElementType:         types.StringType,
								Optional:            true,
								Computed:            true,
								Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
							},
							"upgrade_flags": schema.ListAttribute{
								MarkdownDescription: "Optional MKE bootstrapper update flags",
								ElementType:         types.StringType,
								Optional:            true,
								Computed:            true,
								Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
							},
						},
					},

					"msr": schema.ListNestedBlock{
						MarkdownDescription: "MSR installation configuration",

						Validators: []validator.List{
							listvalidator.SizeAtMost(1),
						},
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"version": schema.StringAttribute{
									MarkdownDescription: "MCR version to install",
									Required:            true,
								},
								"image_repo": schema.StringAttribute{
									MarkdownDescription: "Image repo for MSR images",
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString("docker.io/mirantis"),
								},
								"replica_ids": schema.StringAttribute{
									MarkdownDescription: "MSR replica IDs as a string",
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString("admin"),
								},

								"install_flags": schema.ListAttribute{
									MarkdownDescription: "Optional MSR bootstrapper install flags",
									ElementType:         types.StringType,
									Optional:            true,
									Computed:            true,
									Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
								},
								"upgrade_flags": schema.ListAttribute{
									MarkdownDescription: "Optional MSR bootstrapper update flags",
									ElementType:         types.StringType,
									Optional:            true,
									Computed:            true,
									Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
								},
							},
						},
					},

					"host": schema.ListNestedBlock{
						MarkdownDescription: "Individual host configuration, for each machine in the cluster",

						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},

						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"role": schema.StringAttribute{
									MarkdownDescription: "Host machine role in the cluster",
									Required:            true,
								},
							},
							Blocks: map[string]schema.Block{

								"hooks": schema.ListNestedBlock{
									MarkdownDescription: "Hook configuration for the host",

									Validators: []validator.List{
										listvalidator.SizeAtMost(1),
									},

									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{},
										Blocks: map[string]schema.Block{

											"apply": schema.ListNestedBlock{
												MarkdownDescription: "Launchpad.Apply string hooks for the host",

												Validators: []validator.List{
													listvalidator.SizeAtMost(1),
												},

												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"before": schema.ListAttribute{
															MarkdownDescription: "String hooks to run on hosts before the Apply operation is run.",
															ElementType:         types.StringType,
															Optional:            true,
															Computed:            true,
															Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
														},
														"after": schema.ListAttribute{
															MarkdownDescription: "String hooks to run on hosts after the Apply operation is run.",
															ElementType:         types.StringType,
															Optional:            true,
															Computed:            true,
															Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
														},
													},
												},
											},
										},
									},
								},

								"mcr_config": schema.ListNestedBlock{
									MarkdownDescription: "MCR configuration for the host",

									Validators: []validator.List{
										listvalidator.SizeAtMost(1),
									},

									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"debug": schema.BoolAttribute{
												MarkdownDescription: "Log level",
												Optional:            true,
												Computed:            true,
												Default:             booldefault.StaticBool(false),
											},
											"bip": schema.StringAttribute{
												MarkdownDescription: "Base IP",
												Optional:            true,
												Computed:            true,
												Default:             stringdefault.StaticString(""),
											},

											"default_address_pools": schema.ListNestedAttribute{
												MarkdownDescription: "Reassign docker subnets",

												Optional: true,

												NestedObject: schema.NestedAttributeObject{

													Validators: []validator.Object{
														// @todo validate that only base and size attributes were added
														// because in our unit test I was able to add fields that
														// did not exist.
														// @see: "github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
													},

													Attributes: map[string]schema.Attribute{
														"base": schema.StringAttribute{
															MarkdownDescription: "The CIDR range allocated for bridge networks in each IP address pool.",
															Required:            true,
														},
														"size": schema.Int64Attribute{
															MarkdownDescription: "The CIDR netmask that determines the subnet size to allocate from the base pool.",
															Default:             int64default.StaticInt64(16),
															Optional:            true,
															Computed:            true,
														},
													},
												},
											},
										},
									},
								},

								"ssh": schema.ListNestedBlock{
									MarkdownDescription: "SSH configuration for the host",

									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"address": schema.StringAttribute{
												MarkdownDescription: "SSH endpoint",
												Required:            true,
											},
											"key_path": schema.StringAttribute{
												MarkdownDescription: "SSH endpoint",
												Required:            true,
											},
											"user": schema.StringAttribute{
												MarkdownDescription: "SSH endpoint",
												Required:            true,
											},
											"port": schema.Int64Attribute{
												MarkdownDescription: "SSH Port",
												Optional:            true,
												Computed:            true,
												Default:             int64default.StaticInt64(22),
											},
										},
									},
								},
								"winrm": schema.ListNestedBlock{
									MarkdownDescription: "WinRM configuration for the host",

									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"address": schema.StringAttribute{
												MarkdownDescription: "WinRM endpoint",
												Required:            true,
											},
											"user": schema.StringAttribute{
												MarkdownDescription: "WinRM user",
												Required:            true,
											},
											"password": schema.StringAttribute{
												MarkdownDescription: "WinRM password",
												Required:            true,
											},
											"port": schema.Int64Attribute{
												MarkdownDescription: "WinRM Port",
												Optional:            true,
												Computed:            true,
												Default:             int64default.StaticInt64(5985),
											},
											"use_https": schema.BoolAttribute{
												MarkdownDescription: "If false, then no HTTP is used for winrm transport",
												Optional:            true,
												Computed:            true,
												Default:             booldefault.StaticBool(true),
											},
											"insecure": schema.BoolAttribute{
												MarkdownDescription: "If false, then no SSL certificate validation is used",
												Optional:            true,
												Computed:            true,
												Default:             booldefault.StaticBool(true),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type launchpadSchema14Model struct {
	Id          types.String `tfsdk:"id"`
	SkipDestroy types.Bool   `tfsdk:"skip_destroy"`

	Metadata launchpadSchema14ModelMetadata `tfsdk:"metadata"`
	Spec     launchpadSchema14ModelSpec     `tfsdk:"spec"`
}

// ClusterEqual compare with another state, to see it they are different enough to warrant running launchpad.
func (ls launchpadSchema14Model) ClusterEqual(c launchpadSchema14Model) bool {
	return reflect.DeepEqual(ls.Spec, c.Spec)
}

// ClusterConfig convert this state object into a proper ClusterConfig.
func (ls launchpadSchema14Model) ClusterConfig(diags *diag.Diagnostics) mcc_mke_api.ClusterConfig {
	cc := mcc_mke_api.ClusterConfig{
		APIVersion: "launchpad.mirantis.com/mke/v1.4",
		Kind:       "mke",

		Metadata: func() *mcc_mke_api.ClusterMeta {
			return &mcc_mke_api.ClusterMeta{
				Name: ls.Metadata.Name.String(),
			}
		}(),

		Spec: &mcc_mke_api.ClusterSpec{
			Cluster: mcc_mke_api.Cluster{
				Prune: false,
			},

			Hosts: mcc_mke_api.Hosts{},

			MCR: mcc_common_api.MCRConfig{
				Version:           ls.Spec.MCR.Version.ValueString(),
				InstallURLLinux:   ls.Spec.MCR.InstallURLLinux.ValueString(),
				InstallURLWindows: ls.Spec.MCR.InstallURLWindows.ValueString(),
				RepoURL:           ls.Spec.MCR.RepoURL.ValueString(),
				Channel:           ls.Spec.MCR.Channel.ValueString(),
			},

			MKE: mcc_mke_api.MKEConfig{
				AdminUsername:   ls.Spec.MKE.AdminUsername.ValueString(),
				AdminPassword:   ls.Spec.MKE.AdminPassword.ValueString(),
				ImageRepo:       ls.Spec.MKE.ImageRepo.ValueString(),
				Version:         ls.Spec.MKE.Version.ValueString(),
				InstallFlags:    mcc_common_api.Flags{},
				UpgradeFlags:    mcc_common_api.Flags{},
				Metadata:        &mcc_mke_api.MKEMetadata{},
				LicenseFilePath: ls.Spec.MKE.LicenseFilePath.ValueString(),
			},

			MSR: nil,
		},
	}

	for _, msr := range ls.Spec.MSR {
		cc.Spec.MSR = &mcc_mke_api.MSRConfig{
			ImageRepo:    msr.ImageRepo.ValueString(),
			Version:      msr.Version.ValueString(),
			ReplicaIDs:   msr.ReplicaIDs.ValueString(),
			InstallFlags: mcc_common_api.Flags{},
			UpgradeFlags: mcc_common_api.Flags{},
		}

		if !msr.InstallFlags.IsNull() {
			var fvs []string
			if diag := msr.InstallFlags.ElementsAs(context.Background(), &fvs, true); diag == nil {
				cc.Spec.MSR.InstallFlags = mcc_common_api.Flags(fvs)
			}
		}
	}

	hasMSRHosts := false
	for _, host := range ls.Spec.Hosts {
		if host.Role.ValueString() == HostRoleMSR {
			hasMSRHosts = true
		}
	}

	if hasMSRHosts && cc.Spec.MSR == nil {
		diags.AddError("MSR hosts were provided, but there is no MSR configuration.", "You have MSR hosts in your configuration, but you have not provided an MSR configuration block. You must add the msr block to apply.")
	} else if cc.Spec.MSR != nil && !hasMSRHosts {
		diags.AddError("MSR config passed without hosts", "You have MSR setup in your configuration, but have not added any msr hosts. This sounds benign, but it causes a pointer error in the launchpad VerifyFacts.ValidateMSRVersionJump() method.")
	}

	if !ls.Spec.MKE.InstallFlags.IsNull() {
		var fvs []string
		if diag := ls.Spec.MKE.InstallFlags.ElementsAs(context.Background(), &fvs, true); diag == nil {
			cc.Spec.MKE.InstallFlags = mcc_common_api.Flags(fvs)
		}
	}
	if !ls.Spec.MKE.UpgradeFlags.IsNull() {
		var fvs []string
		if diag := ls.Spec.MKE.UpgradeFlags.ElementsAs(context.Background(), &fvs, true); diag == nil {
			cc.Spec.MKE.UpgradeFlags = mcc_common_api.Flags(fvs)
		}
	}

	for _, host := range ls.Spec.Hosts {
		mccHost := mcc_mke_api.Host{
			Role:  host.Role.ValueString(),
			Hooks: mcc_common_api.Hooks{},
		}

		if len(host.SSH) > 0 {
			hssh := host.SSH[0]

			mccHost.Connection = k0s_rig.Connection{
				SSH: &k0s_rig.SSH{
					Address: hssh.Address.ValueString(),
					KeyPath: hssh.KeyPath.ValueStringPointer(),
					User:    hssh.User.ValueString(),
					Port:    int(hssh.Port.ValueInt64()),
				},
			}
		} else if len(host.WinRM) > 0 {
			hwinrm := host.WinRM[0]

			mccHost.Connection = k0s_rig.Connection{
				WinRM: &k0s_rig.WinRM{
					Address:  hwinrm.Address.ValueString(),
					Password: hwinrm.Password.ValueString(),
					User:     hwinrm.User.ValueString(),
					Port:     int(hwinrm.Port.ValueInt64()),
					UseHTTPS: hwinrm.UseHTTPS.ValueBool(),
					Insecure: hwinrm.Insecure.ValueBool(),
				},
			}
		}

		if len(host.MCRConfig) > 0 {
			mcrConfig := host.MCRConfig[0]
			daemonConfig := k0s_dig.Mapping{}

			/**
			 * Schema Daemon config, and the dig.Mapping struct
			 *
			 * Launchpad expects a dig.Mapping for config, which is really just a
			 * map[string]interface{}. Launchpad doesn't validate the struct at all
			 * but does Marshall the result into json.
			 * All you need to do here is ensure that you treat it like a string
			 * key map, and ensure that any value that you add can be marshalled
			 * using `json.Marshall`.
			 *
			 * You will want to add some interpretation for any daemon confit that
			 * you expect to be able to pass.
			 *
			 */

			daemonConfig["debug"] = mcrConfig.Debug.ValueBool()
			daemonConfig["bip"] = mcrConfig.Bip.ValueString()

			if len(mcrConfig.DefaultAddressPools) > 0 {
				daps := []interface{}{}
				for _, dap := range mcrConfig.DefaultAddressPools {
					dapm := map[string]interface{}{
						"base": dap.Base.ValueString(),
						"size": dap.Size.ValueInt64(),
					}
					daps = append(daps, dapm)
				}
				daemonConfig["default-address-pools"] = daps
			}

			mccHost.DaemonConfig = daemonConfig
		}

		if len(host.Hooks) > 0 {
			sh := host.Hooks[0]

			if len(sh.Apply) > 0 {
				ha := sh.Apply[0]

				hha := map[string][]string{
					"before": {},
					"after":  {},
				}
				var shab []string
				if diag := ha.Before.ElementsAs(context.Background(), &shab, true); diag == nil {
					hha["before"] = shab
				}
				var shaa []string
				if diag := ha.After.ElementsAs(context.Background(), &shaa, true); diag == nil {
					hha["after"] = shab
				}

				mccHost.Hooks["apply"] = hha
			}

		}

		cc.Spec.Hosts = append(cc.Spec.Hosts, &mccHost)
	}

	return cc
}

type launchpadSchema14ModelMetadata struct {
	Name types.String `tfsdk:"name" json:"name"`
}

type launchpadSchema14ModelSpec struct {
	Cluster []launchpadSchema14ModelCluster  `tfsdk:"cluster"`
	Hosts   []launchpadSchema14ModelSpecHost `tfsdk:"host"`
	MCR     launchpadSchema14ModelSpecMCR    `tfsdk:"mcr"`
	MKE     launchpadSchema14ModelSpecMKE    `tfsdk:"mke"`
	MSR     []launchpadSchema14ModelSpecMSR  `tfsdk:"msr"`
}

type launchpadSchema14ModelCluster struct {
	Prune types.Bool `tfsdk:"prune"`
}

type launchpadSchema14ModelSpecMCR struct {
	Version           types.String `tfsdk:"version"`
	Channel           types.String `tfsdk:"channel"`
	InstallURLLinux   types.String `tfsdk:"install_url_linux"`
	InstallURLWindows types.String `tfsdk:"install_url_windows"`
	RepoURL           types.String `tfsdk:"repo_url"`
}

type launchpadSchema14ModelSpecMKE struct {
	AdminPassword   types.String `tfsdk:"admin_password"`
	AdminUsername   types.String `tfsdk:"admin_username"`
	ImageRepo       types.String `tfsdk:"image_repo"`
	Version         types.String `tfsdk:"version"`
	InstallFlags    types.List   `tfsdk:"install_flags"`
	UpgradeFlags    types.List   `tfsdk:"upgrade_flags"`
	LicenseFilePath types.String `tfsdk:"license_file_path"`
}

type launchpadSchema14ModelSpecMSR struct {
	ImageRepo    types.String `tfsdk:"image_repo"`
	Version      types.String `tfsdk:"version"`
	ReplicaIDs   types.String `tfsdk:"replica_ids"`
	InstallFlags types.List   `tfsdk:"install_flags"`
	UpgradeFlags types.List   `tfsdk:"upgrade_flags"`
}

type launchpadSchema14ModelSpecHost struct {
	Role      types.String                              `tfsdk:"role"`
	Hooks     []launchpadSchema14ModelSpecHostHooks     `tfsdk:"hooks"`
	SSH       []launchpadSchema14ModelSpecHostSSH       `tfsdk:"ssh"`
	WinRM     []launchpadSchema14ModelSpecHostWinrm     `tfsdk:"winrm"`
	MCRConfig []launchpadSchema14ModelSpecHostMCRconfig `tfsdk:"mcr_config"`
}
type launchpadSchema14ModelSpecHostHooks struct {
	Apply []launchpadSchema14ModelSpecHostHookAction `tfsdk:"apply"`
}
type launchpadSchema14ModelSpecHostMCRconfig struct {
	Debug               types.Bool                                                   `json:"debug" tfsdk:"debug"`
	Bip                 types.String                                                 `json:"bip" tfsdk:"bip"`
	DefaultAddressPools []launchpadSchema14ModelSpecHostMCRconfigDefaultAddressPools `json:"default-address-pools" tfsdk:"default_address_pools"`
}
type launchpadSchema14ModelSpecHostMCRconfigDefaultAddressPools struct {
	Base types.String `json:"base" tfsdk:"base"`
	Size types.Int64  `json:"size" tfsdk:"size"`
}
type launchpadSchema14ModelSpecHostHookAction struct {
	Before types.List `tfsdk:"before"`
	After  types.List `tfsdk:"after"`
}
type launchpadSchema14ModelSpecHostSSH struct {
	Address types.String `tfsdk:"address"`
	KeyPath types.String `tfsdk:"key_path"`
	User    types.String `tfsdk:"user"`
	Port    types.Int64  `tfsdk:"port"`
}
type launchpadSchema14ModelSpecHostWinrm struct {
	Address  types.String `tfsdk:"address"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	UseHTTPS types.Bool   `tfsdk:"use_https"`
	Insecure types.Bool   `tfsdk:"insecure"`
}
