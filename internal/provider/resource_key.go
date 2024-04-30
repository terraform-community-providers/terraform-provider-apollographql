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

var _ resource.Resource = &KeyResource{}
var _ resource.ResourceWithImportState = &KeyResource{}

func NewKeyResource() resource.Resource {
	return &KeyResource{}
}

type KeyResource struct {
	client *graphql.Client
}

type KeyResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Role    types.String `tfsdk:"role"`
	GraphId types.String `tfsdk:"graph_id"`
	Token   types.String `tfsdk:"token"`
}

func (r *KeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (r *KeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Apollo GraphQL graph API key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the key.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the key.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Role of the key.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("GRAPH_ADMIN"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"GRAPH_ADMIN",
						"CONSUMER",
						"CONTRIBUTOR",
						"DOCUMENTER",
						"OBSERVER",
						"PERSISTED_QUERY_PUBLISHER",
					),
				},
			},
			"graph_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the graph the key belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token of the key.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *KeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *KeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := createKey(ctx, *r.client, data.GraphId.ValueString(), data.Name.ValueString(), data.Role.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create key, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a key")

	key := response.Service.NewKey.Key

	data.Id = types.StringValue(key.Id)
	data.Name = types.StringValue(key.KeyName)
	data.Role = types.StringValue(key.Role)
	data.Token = types.StringValue(key.Token)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *KeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	key, err := readKey(ctx, *r.client, data.GraphId.ValueString(), data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read key, got error: %s", err))
		return
	}

	data.Id = types.StringValue(key.Id)
	data.Name = types.StringValue(key.KeyName)
	data.Role = types.StringValue(key.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *KeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := updateKey(ctx, *r.client, data.GraphId.ValueString(), data.Id.ValueString(), data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update key, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a key")

	key := response.Service.RenameKey.Key

	data.Id = types.StringValue(key.Id)
	data.Name = types.StringValue(key.KeyName)
	data.Role = types.StringValue(key.Role)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *KeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := deleteKey(ctx, *r.client, data.GraphId.ValueString(), data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete key, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a key")
}

func (r *KeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func readKey(ctx context.Context, client graphql.Client, serviceId string, keyId string) (*Key, error) {
	response, err := listKeys(ctx, client, serviceId)

	if err != nil {
		return nil, err
	}

	for _, key := range response.Service.ApiKeys {
		if key.Id == keyId {
			return &key.Key, nil
		}
	}

	return nil, fmt.Errorf("Unable to find key with id: %s", keyId)
}
