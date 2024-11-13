package contract

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

func NewContractResource() resource.Resource {
	return &contractResource{}
}

type contractResource struct {
	client *api.PlatformClient
}

type contractResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	SourceGraphName        types.String `tfsdk:"source"`
	Namespace              types.String `tfsdk:"namespace"`
	RoutingURL             types.String `tfsdk:"routing_url"`
	ExcludeTags            types.List   `tfsdk:"exclude_tags"`
	IncludeTags            types.List   `tfsdk:"include_tags"`
	Readme                 types.String `tfsdk:"readme"`
	AdmissionWebhookUrl    types.String `tfsdk:"admission_webhook_url"`
	AdmissionWebhookSecret types.String `tfsdk:"admission_webhook_secret"`
}

func (r *contractResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contract"
}

func (r *contractResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
A contract is a Terraform resource representing a single subgraph with GraphQL Federation enabled, allowing developers to build versatile, multi-audience graphs while simplifying development and ensuring maintainability. 

For more information, refer to the Cosmo Documentation at https://cosmo-docs.wundergraph.com/concepts/schema-contracts.
		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.StringAttribute{
				Required: true,
			},
			"exclude_tags": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"include_tags": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"readme": schema.StringAttribute{
				Optional: true,
			},
			"admission_webhook_url": schema.StringAttribute{
				Optional: true,
			},
			"admission_webhook_secret": schema.StringAttribute{
				Optional: true,
			},
			"routing_url": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *contractResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PlatformClient)
	if !ok {
		utils.AddDiagnosticError(resp,
			ErrUnexpectedResourceType,
			fmt.Sprintf("Expected *api.PlatformClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *contractResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data contractResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, apiError := r.createAndFetchContract(ctx, data, resp)
	if apiError != nil {
		if !api.IsSubgraphCompositionFailedError(apiError) {
			// returning only for a non composition error, as the resource exists otherwise
			return
		}
	}

	graph := response.Graph
	data.Id = types.StringValue(graph.GetId())
	data.Name = types.StringValue(graph.GetName())
	data.Namespace = types.StringValue(graph.GetNamespace())
	data.RoutingURL = types.StringValue(graph.GetRoutingURL())

	var responseExcludeTags []attr.Value
	for _, tag := range graph.Contract.GetExcludeTags() {
		responseExcludeTags = append(responseExcludeTags, types.StringValue(tag))
	}
	data.ExcludeTags = types.ListValueMust(types.StringType, responseExcludeTags)

	var responseIncludeTags []attr.Value
	for _, tag := range graph.Contract.IncludeTags {
		responseIncludeTags = append(responseIncludeTags, types.StringValue(tag))
	}
	data.IncludeTags = types.ListValueMust(types.StringType, responseIncludeTags)

	if graph.Readme != nil {
		data.Readme = types.StringValue(*graph.Readme)
	}

	if graph.AdmissionWebhookUrl != nil {
		data.AdmissionWebhookUrl = types.StringValue(*graph.AdmissionWebhookUrl)
	}

	utils.LogAction(ctx, DebugCreate, data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *contractResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data contractResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if data.Id.IsNull() || data.Id.ValueString() == "" {
		utils.AddDiagnosticError(resp, ErrInvalidResourceID, "Cannot read federated graph without an ID.")
		return
	}

	response, apiError := r.client.GetFederatedGraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrReadingContract,
				apiError.Error(),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp,
			ErrReadingContract,
			apiError.Error(),
		)
		return
	}

	graph := response.Graph
	data.Id = types.StringValue(graph.GetId())
	data.Name = types.StringValue(graph.GetName())
	data.Namespace = types.StringValue(graph.GetNamespace())
	data.RoutingURL = types.StringValue(graph.GetRoutingURL())

	var responseExcludeTags []attr.Value
	for _, tag := range graph.Contract.GetExcludeTags() {
		responseExcludeTags = append(responseExcludeTags, types.StringValue(tag))
	}
	data.ExcludeTags = types.ListValueMust(types.StringType, responseExcludeTags)

	var responseIncludeTags []attr.Value
	for _, tag := range graph.Contract.IncludeTags {
		responseIncludeTags = append(responseIncludeTags, types.StringValue(tag))
	}
	data.IncludeTags = types.ListValueMust(types.StringType, responseIncludeTags)

	if graph.Readme != nil {
		data.Readme = types.StringValue(*graph.Readme)
	}

	if graph.AdmissionWebhookUrl != nil {
		data.AdmissionWebhookUrl = types.StringValue(*graph.AdmissionWebhookUrl)
	}

	utils.LogAction(ctx, "read", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *contractResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data contractResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	excludeTags, err := utils.ConvertLabelMatchers(data.ExcludeTags)
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrUpdatingContract,
			err.Error(),
		)
		return
	}

	includeTags, err := utils.ConvertLabelMatchers(data.IncludeTags)
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingContract,
			"Could not create contract: "+err.Error(),
		)
		return
	}

	if len(includeTags) > 0 && len(excludeTags) > 0 {
		utils.AddDiagnosticError(resp, ErrUpdatingContract, "a contract cannot have both the include and exclude tags")
		return
	}

	_, apiError := r.client.UpdateContract(ctx, data.Name.ValueString(), data.Namespace.ValueString(), excludeTags, includeTags)
	if apiError != nil {
		utils.AddDiagnosticError(resp,
			ErrUpdatingContract,
			apiError.Error(),
		)
		return
	}

	response, apiError := r.client.GetFederatedGraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		utils.AddDiagnosticError(resp,
			ErrRetrievingContract,
			apiError.Error(),
		)
		return
	}

	graph := response.Graph
	data.Id = types.StringValue(graph.GetId())
	data.Name = types.StringValue(graph.GetName())
	data.Namespace = types.StringValue(graph.GetNamespace())
	data.RoutingURL = types.StringValue(graph.GetRoutingURL())

	var responseExcludeTags []attr.Value
	for _, tag := range graph.Contract.GetExcludeTags() {
		responseExcludeTags = append(responseExcludeTags, types.StringValue(tag))
	}
	data.ExcludeTags = types.ListValueMust(types.StringType, responseExcludeTags)

	var responseIncludeTags []attr.Value
	for _, tag := range graph.Contract.IncludeTags {
		responseIncludeTags = append(responseIncludeTags, types.StringValue(tag))
	}
	data.IncludeTags = types.ListValueMust(types.StringType, responseIncludeTags)

	if graph.Readme != nil {
		data.Readme = types.StringValue(*graph.Readme)
	}

	if graph.AdmissionWebhookUrl != nil {
		data.AdmissionWebhookUrl = types.StringValue(*graph.AdmissionWebhookUrl)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *contractResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data contractResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiError := r.client.DeleteContract(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticError(resp,
				ErrDeletingContract,
				apiError.Error(),
			)
			resp.State.RemoveResource(ctx)
		} else {
			utils.AddDiagnosticError(resp,
				ErrDeletingContract,
				apiError.Error(),
			)
			return
		}
	}

	utils.LogAction(ctx, "deleted contract", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())
}

func (r *contractResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *contractResource) createAndFetchContract(ctx context.Context, data contractResourceModel, resp *resource.CreateResponse) (*platformv1.GetFederatedGraphByNameResponse, *api.ApiError) {
	excludeTags, err := utils.ConvertLabelMatchers(data.ExcludeTags)
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingContract,
			"Could not create contract: "+err.Error(),
		)
		return nil, &api.ApiError{Err: err, Reason: "CreateContract", Status: common.EnumStatusCode_ERR}
	}

	includeTags, err := utils.ConvertLabelMatchers(data.IncludeTags)
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingContract,
			"Could not create contract: "+err.Error(),
		)
		return nil, &api.ApiError{Err: err, Reason: "CreateContract", Status: common.EnumStatusCode_ERR}
	}

	if len(includeTags) > 0 && len(excludeTags) > 0 {
		utils.AddDiagnosticError(resp, ErrCreatingContract, "a contract cannot have both the include and exclude tags")
		return nil, &api.ApiError{Err: err, Reason: "CreateContract", Status: common.EnumStatusCode_ERR}
	}

	utils.DebugAction(ctx, DebugCreate, data.Name.ValueString(), data.Namespace.ValueString(), map[string]interface{}{
		"routing_url": data.RoutingURL.ValueString(),
		"excludeTags": strings.Join(excludeTags, ","),
		"includeTags": strings.Join(includeTags, ","),
	})

	_, apiError := r.client.CreateContract(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.SourceGraphName.ValueString(), data.RoutingURL.ValueString(), data.AdmissionWebhookUrl.ValueString(), data.AdmissionWebhookSecret.ValueString(), excludeTags, includeTags, data.Readme.ValueString())
	if apiError != nil {
		if api.IsContractCompositionFailedError(apiError) || api.IsSubgraphCompositionFailedError(apiError) {
			utils.AddDiagnosticError(resp,
				ErrCreatingContract,
				"Contract composition failed: "+apiError.Error(),
			)
		} else {
			utils.AddDiagnosticError(resp,
				ErrCreatingContract,
				"Could not create contract: "+apiError.Error(),
			)
			return nil, apiError
		}
	}

	response, apiError := r.client.GetFederatedGraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		return nil, apiError
	}

	utils.DebugAction(ctx, DebugCreate, data.Name.ValueString(), data.Namespace.ValueString(), map[string]interface{}{
		"id":    response.Graph.GetId(),
		"graph": response.Graph,
	})

	return response, nil
}
