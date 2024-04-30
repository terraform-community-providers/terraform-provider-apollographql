package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &VariantResource{}
var _ resource.ResourceWithImportState = &VariantResource{}

func NewVariantResource() resource.Resource {
	return &VariantResource{}
}

type VariantResource struct {
	client *graphql.Client
}

type VariantResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Public  types.Bool   `tfsdk:"public"`
	GraphId types.String `tfsdk:"graph_id"`
}

func (r *VariantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variant"
}

func (r *VariantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Apollo GraphQL graph variant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the variant.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the variant.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"public": schema.BoolAttribute{
				MarkdownDescription: "Whether the variant is public.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"graph_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the graph the variant belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *VariantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *VariantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VariantResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := createVariant(ctx, *r.client, data.GraphId.ValueString(), data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create variant, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a variant")

	variant, err := readVariant(ctx, *r.client, data.GraphId.ValueString(), data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read variant, got error: %s", err))
		return
	}

	if data.Public.ValueBool() {
		err := updatePublic(ctx, *r.client, data)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update variant, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "updated a variant")
	}

	data.Id = types.StringValue(variant.Id)
	data.Name = types.StringValue(variant.Name)
	data.GraphId = types.StringValue(variant.GraphId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VariantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VariantResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	variant, err := readVariant(ctx, *r.client, data.GraphId.ValueString(), data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read variant, got error: %s", err))
		return
	}

	data.Id = types.StringValue(variant.Id)
	data.Name = types.StringValue(variant.Name)
	data.Public = types.BoolValue(variant.IsPublic)
	data.GraphId = types.StringValue(variant.GraphId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VariantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VariantResourceModel
	var state *VariantResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Public.ValueBool() != state.Public.ValueBool() {
		err := updatePublic(ctx, *r.client, data)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update variant, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "updated a variant")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VariantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VariantResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteVariant(ctx, *r.client, data.GraphId.ValueString(), data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete variant, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a variant")
}

func (r *VariantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: graph_id:name. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("graph_id"), parts[0])...)
}

func readVariant(ctx context.Context, client graphql.Client, serviceId string, variantName string) (*Variant, error) {
	response, err := getVariant(ctx, client, serviceId, variantName)

	if err != nil {
		return nil, err
	}

	variant := response.Service.Variant.Variant

	return &variant, nil
}

func updatePublic(ctx context.Context, client graphql.Client, data *VariantResourceModel) error {
	response, err := updateVariantIsPublic(ctx, client, data.GraphId.ValueString(), data.Name.ValueString(), data.Public.ValueBool())

	if err != nil {
		return err
	}

	data.Public = types.BoolValue(response.Service.Variant.UpdateVariantIsPublic.Variant.IsPublic)

	return nil
}
