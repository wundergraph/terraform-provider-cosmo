package router_token

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

type TokenResource struct {
	client *api.PlatformClient
}

type TokenResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	GraphName types.String `tfsdk:"graph_name"`
	Namespace types.String `tfsdk:"namespace"`
	Token     types.String `tfsdk:"token"`
}

func NewTokenResource() resource.Resource {
	return &TokenResource{}
}

func (r *TokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_router_token"
}

func (r *TokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Generates a token that is limited to a federated graph. This token allows the router to interact with the platform and send metrics to the collectors.

For more information on router tokens, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/router/token/create).
		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the router token.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the router token.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"graph_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the graph to create the token for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace to create the token in.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The token to be used for the router.",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *TokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PlatformClient)
	if !ok {
		utils.AddDiagnosticError(resp, ErrUnexpectedDataSourceType, fmt.Sprintf("Expected *api.PlatformClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = client
}

func (r *TokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, apiError := r.client.CreateToken(ctx, data.Name.ValueString(), data.GraphName.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrCreatingToken,
				apiError.Error(),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp,
			ErrCreatingToken,
			apiError.Error(),
		)
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("%s-%s-%s", data.GraphName.ValueString(), data.Namespace.ValueString(), data.Name.ValueString()))
	data.Token = types.StringValue(apiResponse)
	data.Name = types.StringValue(data.Name.ValueString())
	data.GraphName = types.StringValue(data.GraphName.ValueString())
	data.Namespace = types.StringValue(data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// check if the token exists
	_, apiError := r.client.GetToken(ctx, data.Name.ValueString(), data.GraphName.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp, ErrReadingToken, apiError.Error())
		return
	}

	utils.LogAction(ctx, "read", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	utils.AddDiagnosticError(resp, ErrUpdatingToken, "Token update should never be called, please delete and recreate the token")
}

func (r *TokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	apiError := r.client.DeleteToken(ctx, data.Name.ValueString(), data.GraphName.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		utils.AddDiagnosticError(resp,
			ErrDeletingToken,
			apiError.Error(),
		)
		return
	}

	utils.LogAction(ctx, "deleted", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())
}
