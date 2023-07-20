package provider

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"gopkg.in/yaml.v2"

	mcc_mke "github.com/Mirantis/mcc/pkg/product/mke"
	mcc_logrus "github.com/sirupsen/logrus"
)

var _ resource.Resource = &LaunchpadConfigResource{}

type LaunchpadConfigResource struct {
	testingMode bool
}

func NewLaunchpadConfigResource() resource.Resource {
	return &LaunchpadConfigResource{}
}

func (r *LaunchpadConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *LaunchpadConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = launchpadSchema14()
}

func (r *LaunchpadConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	lpm, ok := req.ProviderData.(*LaunchpadProviderModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *LaunchpadProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.testingMode = lpm.testingMode
}

func (r *LaunchpadConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var cls *launchpadSchema14Model

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &cls)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cc := cls.ClusterConfig(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Failed to build cluster config from terraform config",
			"Provider failed to convert resource schema to a usable cluster config",
		)
		return
	}

	c := mcc_mke.MKE{ClusterConfig: cc}

	if err := cc.Validate(); err != nil {
		resp.Diagnostics.AddError(
			"Launchpad config validation failed",
			err.Error(),
		)

		return
	}

	logrusBuffer := &bytes.Buffer{}
	mcc_logrus.SetOutput(logrusBuffer)

	if r.testingMode {
		resp.Diagnostics.AddWarning("testing mode warning", "launchpad config resource handler is in testing mode, no installation will be run.")
	} else if err := c.Apply(false, false, 10); err != nil {
		ccout, _ := yaml.Marshal(cc)
		resp.Diagnostics.AddError(
			"Launchpad apply failed",
			fmt.Sprintf("%s \n\n%s; %s", ccout, err.Error(), logrusBuffer.String()),
		)

		return
	}

	cls.Id = cls.Metadata.Name

	if diags := resp.State.Set(ctx, cls); diags != nil {
		resp.Diagnostics.Append(diags...)
	}
}

func (r *LaunchpadConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// launchpad has no good way to discover existing installation, so we don't do anything
}

func (r *LaunchpadConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// if enough attributes have changed then run apply
	var cls launchpadSchema14Model
	var sls launchpadSchema14Model

	if diags := req.Config.Get(ctx, &cls); diags != nil {
		resp.Diagnostics.Append(diags...)
	}
	if diags := req.State.Get(ctx, &sls); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Launchpad config interpret failed",
			"Failed to interpret either the resource config or state",
		)

		return
	}

	if cls.ClusterEqual(sls) {
		return
	}

	cc := cls.ClusterConfig(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Failed to build cluster config from terraform config",
			"Provider failed to convert resource schema to a usable cluster config",
		)
		return
	}

	c := mcc_mke.MKE{ClusterConfig: cc}

	if err := cc.Validate(); err != nil {
		resp.Diagnostics.AddError(
			"Launchpad config validation failed",
			err.Error(),
		)

		return
	}

	logrusBuffer := &bytes.Buffer{}
	mcc_logrus.SetOutput(logrusBuffer)

	if r.testingMode {
		resp.Diagnostics.AddWarning("testing mode warning", "launchpad config resource handler is in testing mode, no update will be run.")
	} else if err := c.Apply(false, false, 10); err != nil {
		resp.Diagnostics.AddError(
			"Launchpad apply failed",
			fmt.Sprintf("%s; %s", err.Error(), logrusBuffer.String()),
		)

		return
	}

	if diags := resp.State.Set(ctx, cls); diags != nil {
		resp.Diagnostics.Append(diags...)
	}
}

func (r *LaunchpadConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var sls launchpadSchema14Model

	diags := req.State.Get(ctx, &sls)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if sls.SkipDestroy.ValueBool() {
		resp.Diagnostics.AddWarning(
			"Cluster destruction was skipped!",
			"The cluster was not actively destroyed, as configuration told us to skip destruction",
		)

		return
	}

	cc := sls.ClusterConfig(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Failed to build cluster config from terraform config",
			"Provider failed to convert resource schema to a usable cluster config",
		)
		return
	}

	c := mcc_mke.MKE{ClusterConfig: cc}

	logrusBuffer := &bytes.Buffer{}
	mcc_logrus.SetOutput(logrusBuffer)

	if r.testingMode {
		resp.Diagnostics.AddWarning("testing mode warning", "launchpad config resource handler is in testing mode, no reset will be run.")
	} else if err := c.Reset(); err != nil {
		resp.Diagnostics.AddError(
			"Launchpad Reset failed",
			fmt.Sprintf("%s; %s", err.Error(), logrusBuffer.String()),
		)

		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r *LaunchpadConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}
