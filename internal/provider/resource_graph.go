package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &GraphResource{}
var _ resource.ResourceWithImportState = &GraphResource{}

func NewGraphResource() resource.Resource {
	return &GraphResource{}
}

type GraphResource struct {
	client *graphql.Client
}

type GraphResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Title                  types.String `tfsdk:"title"`
	OnboardingArchitecture types.String `tfsdk:"onboarding_architecture"`
	OrganizationId         types.String `tfsdk:"organization_id"`
	Description            types.String `tfsdk:"description"`
}

func (r *GraphResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph"
}

func (r *GraphResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Apollo GraphQL graph.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the graph.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title of the graph.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"onboarding_architecture": schema.StringAttribute{
				MarkdownDescription: "Onboarding architecture of the graph.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("MONOLITH"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("MONOLITH", "SUPERGRAPH"),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the organization the graph belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the graph.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *GraphResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GraphResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *GraphResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := createService(ctx, *r.client, data.Id.ValueString(), data.Title.ValueString(), data.OnboardingArchitecture.ValueString(), data.OrganizationId.ValueString(), data.Description.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create graph, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a graph")

	service := response.NewService.Service

	data.Id = types.StringValue(service.Id)
	data.Title = types.StringValue(service.Title)
	data.OnboardingArchitecture = types.StringValue(service.OnboardingArchitecture)
	data.OrganizationId = types.StringValue(service.AccountId)
	data.Description = types.StringValue(service.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GraphResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GraphResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getService(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read graph, got error: %s", err))
		return
	}

	service := response.Service.Service

	data.Id = types.StringValue(service.Id)
	data.Title = types.StringValue(service.Title)
	data.OnboardingArchitecture = types.StringValue(service.OnboardingArchitecture)
	data.OrganizationId = types.StringValue(service.AccountId)
	data.Description = types.StringValue(service.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GraphResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *GraphResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Allow updating title and description

	// input := ProjectUpdateInput{
	// 	Name:        data.Name.ValueString(),
	// 	Description: data.Description.ValueString(),
	// 	IsPublic:    !data.Private.ValueBool(),
	// 	PrDeploys:   data.HasPrDeploys.ValueBool(),
	// }

	// resp.Diagnostics.Append(data.DefaultVariant.As(ctx, &defaultVariantData, basetypes.ObjectAsOptions{})...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// response, err := updateProject(ctx, *r.client, data.Id.ValueString(), input)

	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update graph, got error: %s", err))
	// 	return
	// }

	// tflog.Trace(ctx, "updated a graph")

	// service := response.ProjectUpdate.Project

	// data.Id = types.StringValue(service.Id)
	// data.Title = types.StringValue(service.Title)
	// data.OnboardingArchitecture = types.StringValue(service.OnboardingArchitecture)
	// data.OrganizationId = types.StringValue(service.AccountId)
	// data.Description = types.StringValue(service.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GraphResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GraphResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteService(ctx, *r.client, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete graph, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a graph")
}

func (r *GraphResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
