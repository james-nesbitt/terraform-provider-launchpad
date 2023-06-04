package launchpad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"

	mcc_config "github.com/Mirantis/mcc/pkg/config"
)

// ResourceConfig for Launchpad config schema
func ResourceYamlConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceYamlConfigCreate,
		ReadContext:   resourceYamlConfigRead,
		UpdateContext: resourceYamlConfigUpdate,
		DeleteContext: resourceYamlConfigDelete,
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
			"yaml_config": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceYamlConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sConfig := d.Get("yaml_config")
	bConfig := []byte(sConfig.(string))

	logrusBuffer := &bytes.Buffer{}
	logrus.SetOutput(logrusBuffer)

	product, err := mcc_config.ProductFromYAML(bConfig)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w; %s", err, logrusBuffer.String()))
	}

	productConfig, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := product.Apply(false, false, 10); err != nil {
		return diag.FromErr(fmt.Errorf("%w; %s\nProductConfig:\n%s", err, logrusBuffer.String(), string(productConfig)))
	}

	if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func resourceYamlConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func resourceYamlConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// if any of the other attributes have changes run create
	if d.HasChangeExcept("last_updated") {
		return resourceConfigCreate(ctx, d, m)
	}
	return resourceConfigRead(ctx, d, m)
}

func resourceYamlConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	skip_destroy := d.Get("skip_destroy").(bool)
	if !skip_destroy {
		logrusBuffer := &bytes.Buffer{}
		logrus.SetOutput(logrusBuffer)

		sConfig := d.Get("yaml_config")
		bConfig := []byte(sConfig.(string))

		product, err := mcc_config.ProductFromYAML(bConfig)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%w; %s", err, logrusBuffer.String()))
		}

		if err := product.Reset(); err != nil {
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
