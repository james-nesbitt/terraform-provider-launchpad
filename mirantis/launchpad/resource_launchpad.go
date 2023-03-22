package launchpad

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"
)

// ResourceConfig for Launchpad config schema
func ResourceConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigCreate,
		ReadContext:   resourceConfigRead,
		UpdateContext: resourceConfigUpdate,
		DeleteContext: resourceConfigDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"skip_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								if v == "" {
									errs = append(errs, fmt.Errorf("%q can't be empty string.: '%s'", key, v))
								}
								return
							},
						},
					},
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prune": {
										Type:     schema.TypeBool,
										Default:  true,
										Optional: true,
									},
								},
							},
						},
						"host": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:     schema.TypeString,
										Required: true,
									},
									"hooks": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"before": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"after": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"ssh": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"key_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"user": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
													Default:  22,
												},
											},
										},
									}, // ssh
									"winrm": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"user": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"password": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"use_https": {
													Type:     schema.TypeBool,
													Default:  true,
													Optional: true,
												},
												"insecure": {
													Type:     schema.TypeBool,
													Default:  true,
													Optional: true,
												},
											},
										},
									},
								}, // winrm
							},
						}, // hosts
						"mcr": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel": {
										Type:     schema.TypeString,
										Required: true,
									},
									"install_url_linux": {
										Type:     schema.TypeString,
										Required: true,
									},
									"install_url_windows": {
										Type:     schema.TypeString,
										Required: true,
									},
									"repo_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						}, // mcr
						"mke": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"admin_password": {
										Type:     schema.TypeString,
										Required: true,
									},
									"admin_username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"image_repo": {
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Type:     schema.TypeString,
										Required: true,
									},
									"install_flags": {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"upgrade_flags": {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						}, // mke
						"msr": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_repo": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"replica_ids": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"install_flags": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						}, // msr
					},
				},
			}, // spec
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mkeClient, err := FlattenInputConfigModel(d)
	if err != nil {
		return diag.FromErr(err)
	}

	logrusBuffer := &bytes.Buffer{}
	logrus.SetOutput(logrusBuffer)

	if err := mkeClient.Apply(false, false, 10); err != nil {
		return diag.FromErr(fmt.Errorf("%w; %s", err, logrusBuffer.String()))
	}

	if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// if any of the other attributes have changes run create
	if d.HasChangeExcept("last_updated") {
		return resourceConfigCreate(ctx, d, m)
	}
	return resourceConfigRead(ctx, d, m)
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	skip_destroy := d.Get("skip_destroy").(bool)
	if !skip_destroy {
		logrusBuffer := &bytes.Buffer{}
		logrus.SetOutput(logrusBuffer)

		mkeClient, err := FlattenInputConfigModel(d)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := mkeClient.Reset(); err != nil {
			return diag.FromErr(fmt.Errorf("%w; %s", err, logrusBuffer.String()))
		}
		return nil
	} else {
		var diags diag.Diagnostics
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "MKE cluster destruction was skipped!",
		})

		return diags
	}
}
