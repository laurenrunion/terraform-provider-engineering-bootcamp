package provider
//the resource itself

import (
	"context"
	"errors"

    "terraform-provider-engineering-bootcamp/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &engineerResource{}
	_ resource.ResourceWithConfigure   = &engineerResource{}
	_ resource.ResourceWithImportState = &engineerResource{}
)

func NewEngineerResource() resource.Resource {
	return &engineerResource{}
}

type engineerResource struct {
	client *client.Client
}

// Terraform-side model (types.String), NOT client.Engineer
type engineerResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

//resource "engineering-bootcamp_engineer" "ryan" name of resource called in tf
func (r *engineerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer" // becomes engineering-bootcamp_engineer
}

//what the resource should look like
func (r *engineerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, // the API creates it
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // keep the id on updates
				},
			},
			"name":  schema.StringAttribute{Required: true},
			"email": schema.StringAttribute{Required: true},
		},
	}
}

// Configure receives the *client.Client that provider.go built.
func (r *engineerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type", "expected *client.Client")
		return
	}
	r.client = c
}

func (r *engineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan engineerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateEngineer(client.Engineer{
		Name:  plan.Name.ValueString(),
		Email: plan.Email.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating engineer", err.Error())
		return
	}

	plan.ID = types.StringValue(created.ID) // fill the computed id
	plan.Name = types.StringValue(created.Name)
	plan.Email = types.StringValue(created.Email)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *engineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state engineerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	eng, err := r.client.GetEngineer(state.ID.ValueString())
	if errors.Is(err, client.ErrNotFound) {
		resp.State.RemoveResource(ctx) // gone from the API, so drop it from state
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading engineer", err.Error())
		return
	}

	state.Name = types.StringValue(eng.Name)
	state.Email = types.StringValue(eng.Email)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *engineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state engineerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdateEngineer(state.ID.ValueString(), client.Engineer{
		Name:  plan.Name.ValueString(),
		Email: plan.Email.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating engineer", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(updated.Name)
	plan.Email = types.StringValue(updated.Email)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *engineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state engineerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEngineer(state.ID.ValueString())
	if err != nil && !errors.Is(err, client.ErrNotFound) { // already gone is fine
		resp.Diagnostics.AddError("Error deleting engineer", err.Error())
	}
}

// terraform import engineering-bootcamp_engineer.x <id>
func (r *engineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
