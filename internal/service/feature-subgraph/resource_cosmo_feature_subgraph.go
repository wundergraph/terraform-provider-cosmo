package feature_subgraph

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

type FeatureSubgraphResource struct {
	client *api.PlatformClient
}

type FeatureSubgraphResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Namespace            types.String `tfsdk:"namespace"`
	RoutingURL           types.String `tfsdk:"routing_url"`
	BaseSubgraphName     types.String `tfsdk:"base_subgraph_name"`
	SubscriptionUrl      types.String `tfsdk:"subscription_url"`
	SubscriptionProtocol types.String `tfsdk:"subscription_protocol"`
	WebsocketSubprotocol types.String `tfsdk:"websocket_subprotocol"`
	Readme               types.String `tfsdk:"readme"`
	Schema               types.String `tfsdk:"schema"`
}

func NewSubgraphResource() resource.Resource {
	return &FeatureSubgraphResource{}

}

func (r *FeatureSubgraphResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PlatformClient)
	if !ok {
		utils.AddDiagnosticError(resp,
			ErrUnexpectedDataSourceType,
			fmt.Sprintf("Expected *api.PlatformClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *FeatureSubgraphResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_subgraph"
}

func (r *FeatureSubgraphResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
this resource handles feature subgraphs. Feature subgraphs are a special type of subgraph that can be used to extend the functionality of the platform.
They require a base subgraph to be specified and can be used to add additional functionality to the base subgraph based on a specialized schema, that is published
to this subgraph.

For more information on feature subgraphs, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/feature-subgraph).
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the feature subgraph.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the feature subgraph.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace to create the feature subgraph in.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
			},
			"routing_url": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{},
				MarkdownDescription: "The routing URL of the feature subgraph.",
			},
			"base_subgraph_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the base subgraph that this feature subgraph extends.",
			},
			"subscription_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The subscription URL for the subgraph.",
			},
			"subscription_protocol": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The subscription protocol for the subgraph.",
				Validators: []validator.String{
					stringvalidator.OneOf(api.GraphQLSubscriptionProtocolWS, api.GraphQLSubscriptionProtocolSSE, api.GraphQLSubscriptionProtocolSSEPost),
				},
				Computed: true,
				Default:  stringdefault.StaticString(api.GraphQLSubscriptionProtocolWS),
			},
			"websocket_subprotocol": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The websocket subprotocol for the subgraph.",
				Validators: []validator.String{
					stringvalidator.OneOf(api.GraphQLWebsocketSubprotocolDefault, api.GraphQLWebsocketSubprotocolGraphQLWS, api.GraphQLWebsocketSubprotocolGraphQLTransportWS),
				},
				Computed: true,
				Default:  stringdefault.StaticString(api.GraphQLWebsocketSubprotocolDefault),
			},
			"readme": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The readme for the subgraph.",
			},
			"schema": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The schema for the subgraph.",
			},
		},
	}
}

func (r *FeatureSubgraphResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeatureSubgraphResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	featureSubgraph, apiError := r.createAndPublishFeatureSubgraph(ctx, data, resp)
	if apiError != nil {
		if !api.IsSubgraphCompositionFailedError(apiError) {
			return
		}
	}

	subgraphSchema, apiError := r.client.GetSubgraphSchema(ctx, featureSubgraph.Name, featureSubgraph.Namespace)
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrFeatureSubgraphNotFound,
				fmt.Sprintf("Subgraph '%s' not found will be recreated %s", data.Name.ValueString(), apiError.Error()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp, ErrRetrievingFeatureSubgraph, fmt.Sprintf("Could not fetch subgraph '%s': %s", data.Name.ValueString(), apiError.Error()))
		return
	}

	data.Id = types.StringValue(featureSubgraph.GetId())
	data.Name = types.StringValue(featureSubgraph.GetName())
	data.Namespace = types.StringValue(featureSubgraph.GetNamespace())
	data.RoutingURL = types.StringValue(featureSubgraph.GetRoutingURL())
	data.SubscriptionProtocol = types.StringValue(featureSubgraph.GetSubscriptionProtocol())
	data.WebsocketSubprotocol = types.StringValue(featureSubgraph.GetWebsocketSubprotocol())
	data.BaseSubgraphName = types.StringValue(featureSubgraph.GetBaseSubgraphName())

	if featureSubgraph.GetSubscriptionUrl() != "" {
		data.SubscriptionUrl = types.StringValue(featureSubgraph.GetSubscriptionUrl())
	}

	if featureSubgraph.Readme != nil {
		data.Readme = types.StringValue(featureSubgraph.GetReadme())
	}

	if len(subgraphSchema) > 0 {
		data.Schema = types.StringValue(subgraphSchema)
	}

	utils.LogAction(ctx, "created subgraph", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeatureSubgraphResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureSubgraphResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		subgraph *platformv1.Subgraph
		apiErr   *api.ApiError
	)

	var fetchSubgraphFunc func() (*platformv1.Subgraph, *api.ApiError)

	if data.Id.ValueString() == "" && (data.Name.ValueString() == "" || data.Namespace.ValueString() == "") {
		utils.AddDiagnosticError(resp, ErrReadingFeatureSubgraph, "Either 'id' or 'name' and 'namespace' must be set")
		return
	}

	if data.Id.ValueString() != "" {
		fetchSubgraphFunc = func() (*platformv1.Subgraph, *api.ApiError) {
			return r.client.GetSubgraphById(ctx, data.Id.ValueString())
		}
	} else {
		fetchSubgraphFunc = func() (*platformv1.Subgraph, *api.ApiError) {
			return r.client.GetSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
		}
	}

	subgraph, apiErr = fetchSubgraphFunc()
	if apiErr != nil {
		if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticWarning(resp, ErrFeatureSubgraphNotFound, apiErr.Error())

			resp.State.RemoveResource(ctx)
			return
		}

		utils.AddDiagnosticError(resp, ErrRetrievingFeatureSubgraph, apiErr.Error())
		return
	}

	// Feature subgraphs are a subset of subgraphs, so we need to check if the subgraph is a feature subgraph
	if !subgraph.IsFeatureSubgraph {
		utils.AddDiagnosticError(resp, ErrReadingFeatureSubgraph, fmt.Sprintf("Subgraph '%s' is not a feature subgraph", data.Name.ValueString()))
		resp.State.RemoveResource(ctx)

		return
	}

	subgraphSchema, apiError := r.client.GetSubgraphSchema(ctx, subgraph.Name, subgraph.Namespace)
	if apiError != nil {
		if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticWarning(resp, ErrFeatureSubgraphSchemaNotFound, apiErr.Error())

			resp.State.RemoveResource(ctx)
			return
		}

		utils.AddDiagnosticError(resp, ErrRetrievingFeatureSubgraphSchema, apiErr.Error())
		return
	}

	data.Id = types.StringValue(subgraph.GetId())
	data.Name = types.StringValue(subgraph.GetName())
	data.Namespace = types.StringValue(subgraph.GetNamespace())
	data.RoutingURL = types.StringValue(subgraph.GetRoutingURL())
	data.SubscriptionProtocol = types.StringValue(subgraph.GetSubscriptionProtocol())
	data.WebsocketSubprotocol = types.StringValue(subgraph.GetWebsocketSubprotocol())
	data.BaseSubgraphName = types.StringValue(subgraph.GetBaseSubgraphName())

	if subgraph.GetSubscriptionUrl() != "" {
		data.SubscriptionUrl = types.StringValue(subgraph.GetSubscriptionUrl())
	}

	if subgraph.Readme != nil {
		data.Readme = types.StringValue(subgraph.GetReadme())
	}

	if len(subgraphSchema) > 0 {
		data.Schema = types.StringValue(subgraphSchema)
	}

	utils.LogAction(ctx, "read subgraph", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeatureSubgraphResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData FeatureSubgraphResourceModel
	var state FeatureSubgraphResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)

	if resp.Diagnostics.HasError() {
		return
	}
	readme := utils.GetValueOrDefault(planData.Readme.ValueStringPointer(), "")
	subscriptionUrl := utils.GetValueOrDefault(planData.SubscriptionUrl.ValueStringPointer(), "")
	subscriptionProtocol := utils.GetValueOrDefault(planData.SubscriptionProtocol.ValueStringPointer(), api.GraphQLSubscriptionProtocolWS)
	websocketSubprotocol := utils.GetValueOrDefault(planData.WebsocketSubprotocol.ValueStringPointer(), api.GraphQLWebsocketSubprotocolDefault)

	apiErr := r.client.UpdateSubgraph(ctx, &platformv1.UpdateSubgraphRequest{
		Name:                 planData.Name.ValueString(),
		Namespace:            planData.Namespace.ValueString(),
		RoutingUrl:           planData.RoutingURL.ValueStringPointer(),
		Labels:               nil,
		SubscriptionUrl:      &subscriptionUrl,
		SubscriptionProtocol: api.ResolveSubscriptionProtocol(subscriptionProtocol),
		WebsocketSubprotocol: api.ResolveWebsocketSubprotocol(websocketSubprotocol),
		Readme:               &readme,
		Headers:              []string{},
	})

	if apiErr != nil {
		if api.IsSubgraphCompositionFailedError(apiErr) {
			utils.AddDiagnosticWarning(resp,
				ErrFeatureSubgraphCompositionFailed,
				apiErr.Error(),
			)
		} else if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticError(resp,
				ErrUpdatingFeatureSubgraph,
				apiErr.Error(),
			)
			resp.State.RemoveResource(ctx)
			return
		} else {
			utils.AddDiagnosticError(resp,
				ErrUpdatingFeatureSubgraph,
				apiErr.Error(),
			)
			return
		}
	}

	if planData.Schema.ValueString() != "" {
		err := r.publishSubgraphSchema(ctx, planData)
		if err != nil {
			if api.IsNotFoundError(err) {
				utils.AddDiagnosticError(resp,
					ErrUpdatingFeatureSubgraph,
					err.Error(),
				)
				resp.State.RemoveResource(ctx)
				return
			} else if api.IsSubgraphCompositionFailedError(err) {
				utils.AddDiagnosticError(resp, ErrFeatureSubgraphCompositionFailed, err.Error())
			} else {
				utils.AddDiagnosticError(resp, ErrPublishingFeatureSubgraph, err.Error())
				return
			}
		}
	}

	subgraph, err := r.client.GetSubgraph(ctx, planData.Name.ValueString(), planData.Namespace.ValueString())
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrRetrievingFeatureSubgraph,
			err.Error(),
		)
		return
	}

	subgraphSchema, apiError := r.client.GetSubgraphSchema(ctx, subgraph.Name, subgraph.Namespace)
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrFeatureSubgraphSchemaNotFound,
				fmt.Sprintf("Schema from subgraph '%s' not found will be recreated %s", planData.Name.ValueString(), apiError.Error()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp, ErrRetrievingFeatureSubgraphSchema, fmt.Sprintf("Could not fetch sceham from subgraph '%s': %s", planData.Name.ValueString(), apiError.Error()))
		return
	}

	planData.Id = types.StringValue(subgraph.GetId())
	planData.Name = types.StringValue(subgraph.GetName())
	planData.Namespace = types.StringValue(subgraph.GetNamespace())
	planData.RoutingURL = types.StringValue(subgraph.GetRoutingURL())
	planData.SubscriptionProtocol = types.StringValue(subgraph.GetSubscriptionProtocol())
	planData.WebsocketSubprotocol = types.StringValue(subgraph.GetWebsocketSubprotocol())
	planData.BaseSubgraphName = types.StringValue(subgraph.GetBaseSubgraphName())

	if subgraph.GetSubscriptionUrl() != "" {
		planData.SubscriptionUrl = types.StringValue(subgraph.GetSubscriptionUrl())
	}

	if subgraph.Readme != nil {
		planData.Readme = types.StringValue(subgraph.GetReadme())
	}

	if len(subgraphSchema) > 0 {
		planData.Schema = types.StringValue(subgraphSchema)
	}

	utils.LogAction(ctx, "updated", planData.Id.ValueString(), planData.Name.ValueString(), planData.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *FeatureSubgraphResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeatureSubgraphResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiErr := r.client.DeleteSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiErr != nil {
		if api.IsSubgraphCompositionFailedError(apiErr) {
			utils.AddDiagnosticWarning(resp,
				ErrDeletingFeatureSubgraph,
				apiErr.Error(),
			)
		} else if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticError(resp,
				ErrDeletingFeatureSubgraph,
				apiErr.Error(),
			)
			resp.State.RemoveResource(ctx)
		} else {
			utils.AddDiagnosticError(resp,
				ErrDeletingFeatureSubgraph,
				apiErr.Error(),
			)
			return
		}
	}

	utils.LogAction(ctx, "deleted subgraph", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

}

func (r *FeatureSubgraphResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FeatureSubgraphResource) createAndPublishFeatureSubgraph(ctx context.Context, data FeatureSubgraphResourceModel, resp *resource.CreateResponse) (*platformv1.Subgraph, *api.ApiError) {
	routingURL := data.RoutingURL.ValueString()
	isFeatureSubgraph := true

	requestData := &platformv1.CreateFederatedSubgraphRequest{
		Name:                 data.Name.ValueString(),
		Namespace:            data.Namespace.ValueString(),
		RoutingUrl:           &routingURL,
		BaseSubgraphName:     data.BaseSubgraphName.ValueStringPointer(),
		IsFeatureSubgraph:    &isFeatureSubgraph,
		SubscriptionUrl:      data.SubscriptionUrl.ValueStringPointer(),
		Readme:               data.Readme.ValueStringPointer(),
		SubscriptionProtocol: api.ResolveSubscriptionProtocol(data.SubscriptionProtocol.ValueString()),
		WebsocketSubprotocol: api.ResolveWebsocketSubprotocol(data.WebsocketSubprotocol.ValueString()),
	}

	apiErr := r.client.CreateSubgraph(ctx, requestData)
	if apiErr != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingFeatureSubgraph,
			apiErr.Error(),
		)
		return nil, apiErr
	}

	if data.Schema.ValueString() != "" {
		apiError := r.publishSubgraphSchema(ctx, data)
		if apiError != nil {
			if api.IsNotFoundError(apiError) {
				utils.AddDiagnosticError(resp,
					ErrUpdatingFeatureSubgraph,
					apiError.Error(),
				)
				resp.State.RemoveResource(ctx)
				return nil, apiError
			} else if api.IsSubgraphCompositionFailedError(apiError) {
				utils.AddDiagnosticError(resp, ErrFeatureSubgraphCompositionFailed, apiError.Error())
			} else {
				utils.AddDiagnosticError(resp, ErrPublishingFeatureSubgraph, apiError.Error())
				return nil, apiError
			}
		}
	}

	subgraph, apiErr := r.client.GetSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiErr != nil {
		return nil, apiErr
	}

	return subgraph, nil
}

func (r *FeatureSubgraphResource) publishSubgraphSchema(ctx context.Context, data FeatureSubgraphResourceModel) *api.ApiError {
	_, apiError := r.client.PublishSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.Schema.ValueString())
	return apiError
}
