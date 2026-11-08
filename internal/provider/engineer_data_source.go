package provider // new file

import (
	"context"

    "terraform-provider-engineering-bootcamp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &engineerDataSource{}
	_ datasource.DataSourceWithConfigure = &engineerDataSource{}
)

func NewEngineerDataSource() datasource.DataSource {
	return &engineerDataSource{}
}

type engineerDataSource struct {
	client *client.Client
}

type engineerDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

func (d *engineerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

func (d *engineerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":    schema.StringAttribute{Required: true},  // user gives the id
			"name":  schema.StringAttribute{Computed: true}, // filled from the API
			"email": schema.StringAttribute{Computed: true}, // filled from the API
		},
	}
}

func (d *engineerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *engineerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config engineerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	e, err := d.client.GetEngineer(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading engineer", err.Error())
		return
	}

	config.Name = types.StringValue(e.Name)
	config.Email = types.StringValue(e.Email)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
