// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

    "terraform-provider-engineering-bootcamp/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ScaffoldingProvider{}
var _ provider.ProviderWithFunctions = &ScaffoldingProvider{}
var _ provider.ProviderWithEphemeralResources = &ScaffoldingProvider{}
var _ provider.ProviderWithActions = &ScaffoldingProvider{}

// ScaffoldingProvider defines the provider implementation.
//the struct for the provider with defined methods(5 metadata,etc.) below
type ScaffoldingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
//changed to host could be anything. Go side schema. the model
//Was endpoint was better but oh well.
type ScaffoldingProviderModel struct {
	Host types.String `tfsdk:"host"`
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "engineering-bootcamp" // changed: was "scaffolding". This is the prefix for resource names
	resp.Version = p.version
}

func (p *ScaffoldingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{ //change to host
				MarkdownDescription: "URL of the DevOps API. Defaults to http://localhost:8080",
				Optional:            true,
			},
		},
	}
}

func (p *ScaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// removed: the "Configuration values are now available" comment block

	// added: turn the Terraform string into a plain Go string (nil if host wasn't set)
	var host *string
	if !data.Host.IsNull() && !data.Host.IsUnknown() {
		h := data.Host.ValueString()
		host = &h
	}

	// changed: was "client := http.DefaultClient"
	// now we build OUR client from internal/client/client.go
	c, err := client.NewClient(host)
	if err != nil {
		// added: show the error to the user if the client can't be built
		resp.Diagnostics.AddError("Unable to create DevOps API client", err.Error())
		return
	}

	resp.DataSourceData = c // changed: was the http client, now our client
	resp.ResourceData = c   // changed: was the http client, now our client
}

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
        NewEngineerResource, // changed: was NewExampleResource
	}
}

func (p *ScaffoldingProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewExampleEphemeralResource,
	}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
        NewEngineerDataSource, // changed: was NewExampleDataSource
	}
}

func (p *ScaffoldingProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func (p *ScaffoldingProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		NewExampleAction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
