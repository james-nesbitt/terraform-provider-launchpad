package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	TestingVersion = "test"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &LaunchpadProvider{}

// LaunchpadProvider defines the provider implementation.
type LaunchpadProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// LaunchpadProviderModel describes the provider data model.
type LaunchpadProviderModel struct {
	testingMode bool
}

func (p *LaunchpadProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "launchpad"
	resp.Version = p.version
}

func (p *LaunchpadProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *LaunchpadProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LaunchpadProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if p.version == TestingVersion {
		data.testingMode = true
	}

	resp.ResourceData = &data
	resp.DataSourceData = &data

	AllLoggingToTFLog()
}

func (p *LaunchpadProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLaunchpadConfigResource,
	}
}

func (p *LaunchpadProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LaunchpadProvider{
			version: version,
		}
	}
}
