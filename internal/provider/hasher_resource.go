package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type HasherResourceModel struct {
	ID      tftypes.String `tfsdk:"id"`
	InputWO tftypes.String `tfsdk:"input_wo"`
	Output  tftypes.String `tfsdk:"output"`
}

var _ resource.Resource = &HasherResource{}
var _ resource.ResourceWithModifyPlan = &HasherResource{}

func NewHasherResource() resource.Resource {
	return &HasherResource{}
}

type HasherResource struct {
}

func (r *HasherResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hasher"
}

func (r *HasherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Accepts an ephemeral input and returns a hash of the input as a non-ephemeral output.",
		Attributes: map[string]schema.Attribute{
			"input_wo": schema.StringAttribute{
				MarkdownDescription: "Write-only input (not stored in state). The resulting hash is stored in `output` and `id`.",
				Required:            true,
				WriteOnly:           true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Returns sha256 hash of `input_wo`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"output": schema.StringAttribute{
				MarkdownDescription: "Returns sha256 hash of `input_wo`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *HasherResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *HasherResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return // destroy
	}

	var config HasherResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.InputWO.IsNull() || config.InputWO.IsUnknown() {
		return
	}

	newHash := hashInput(config.InputWO.ValueString())

	if req.State.Raw.IsNull() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("id"), tftypes.StringValue(newHash))...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("output"), tftypes.StringValue(newHash))...)
		return
	}

	var state HasherResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.ID.IsNull() && !state.ID.IsUnknown() && state.ID.ValueString() == newHash {
		return
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("id"), tftypes.StringValue(newHash))...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("output"), tftypes.StringValue(newHash))...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HasherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan HasherResourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hash := tftypes.StringValue(hashInput(plan.InputWO.ValueString()))
	plan.ID = hash
	plan.Output = hash

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *HasherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state HasherResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *HasherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan HasherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config HasherResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hash := tftypes.StringValue(hashInput(config.InputWO.ValueString()))
	plan.Output = hash
	plan.ID = hash

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *HasherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	req.State.RemoveResource(ctx)
}

func hashInput(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
