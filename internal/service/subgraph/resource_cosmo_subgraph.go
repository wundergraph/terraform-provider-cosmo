package subgraph

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

type SubgraphResource struct {
	client *api.PlatformClient
}

type SubgraphResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Namespace  types.String `tfsdk:"namespace"`
	RoutingURL types.String `tfsdk:"routing_url"`
	// TODO: re-enable this once Graph Feature Flags are implementd
	// BaseSubgraphName     types.String `tfsdk:"base_subgraph_name"`
	SubscriptionUrl      types.String `tfsdk:"subscription_url"`
	SubscriptionProtocol types.String `tfsdk:"subscription_protocol"`
	WebsocketSubprotocol types.String `tfsdk:"websocket_subprotocol"`
	Readme               types.String `tfsdk:"readme"`
	IsEventDrivenGraph   types.Bool   `tfsdk:"is_event_driven_graph"`
	IsFeatureSubgraph    types.Bool   `tfsdk:"is_feature_subgraph"`
	UnsetLabels          types.Bool   `tfsdk:"unset_labels"`
	// TBD: This is only used in the update subgraph method and not used atm
	// Headers              types.List   `tfsdk:"headers"`
	Labels types.Map    `tfsdk:"labels"`
	Schema types.String `tfsdk:"schema"`
}

func NewSubgraphResource() resource.Resource {
	return &SubgraphResource{}
}

func (r *SubgraphResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PlatformClient)
	if !ok {
		utils.AddDiagnosticError(resp,
			ErrUnexpectedDataSourceType,
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *SubgraphResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subgraph"
}

func (r *SubgraphResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
This resource handles subgraphs. Each subgraph is responsible for defining its specific segment of the schema and managing the related queries.
		
For more information on subgraphs, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/subgraph).
		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the subgraph resource.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the subgraph.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace in which the subgraph is located.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
			},
			"routing_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The routing URL of the subgraph.",
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
			},
			"readme": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The readme for the subgraph.",
			},
			"websocket_subprotocol": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The websocket subprotocol for the subgraph.",
				Validators: []validator.String{
					stringvalidator.OneOf(api.GraphQLWebsocketSubprotocolDefault, api.GraphQLWebsocketSubprotocolGraphQLWS, api.GraphQLWebsocketSubprotocolGraphQLTransportWS),
				},
			},
			"is_event_driven_graph": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Indicates if the subgraph is event-driven.",
			},
			"is_feature_subgraph": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Indicates if the subgraph is a feature subgraph.",
			},
			// "headers": schema.ListAttribute{
			// 	Optional:            true,
			// 	MarkdownDescription: "Headers for the subgraph.",
			// 	ElementType:         types.StringType,
			// },
			"unset_labels": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Unset labels for the subgraph.",
			},
			"labels": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "Labels for the subgraph.",
				ElementType:         types.StringType,
			},
			"schema": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The schema for the subgraph.",
			},
			// TODO: re-enable this once Graph Feature Flags are implementd
			// "base_subgraph_name": schema.StringAttribute{
			// 	Optional:            true,
			// 	MarkdownDescription: "The base subgraph name.",
			// },
		},
	}
}

func (r *SubgraphResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SubgraphResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	subgraph, err := r.createAndPublishSubgraph(ctx, data, resp)
	if err != nil {
		return
	}

	data.Id = types.StringValue(subgraph.GetId())
	data.Name = types.StringValue(subgraph.GetName())
	data.Namespace = types.StringValue(subgraph.GetNamespace())
	data.RoutingURL = types.StringValue(subgraph.GetRoutingURL())

	utils.LogAction(ctx, "created", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubgraphResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SubgraphResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	subgraph, apiError := r.client.GetSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrSubgraphNotFound,
				fmt.Sprintf("Subgraph '%s' not found will be recreated %s", data.Name.ValueString(), apiError.Error()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp, ErrRetrievingSubgraph, fmt.Sprintf("Could not fetch subgraph '%s': %s", data.Name.ValueString(), apiError.Error()))
		return
	}

	data.Id = types.StringValue(subgraph.GetId())
	data.Name = types.StringValue(subgraph.GetName())
	data.Namespace = types.StringValue(subgraph.GetNamespace())
	data.RoutingURL = types.StringValue(subgraph.GetRoutingURL())

	utils.LogAction(ctx, "read", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubgraphResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SubgraphResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var labels []*platformv1.Label
	for key, value := range data.Labels.Elements() {
		if strValue, ok := value.(types.String); ok {
			labels = append(labels, &platformv1.Label{
				Key:   key,
				Value: strValue.ValueString(),
			})
		}
	}

	if data.Schema.ValueString() != "" {
		apiResponse, err := r.client.PublishSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.Schema.ValueString())
		if err != nil {
			utils.AddDiagnosticError(resp,
				ErrPublishingSubgraph,
				fmt.Sprintf("Could not publish subgraph '%s': %s", data.Name.ValueString(), err.Error()),
			)
			return
		}
		if apiResponse.HasChanged != nil && *apiResponse.HasChanged {
			resp.Diagnostics.AddWarning(
				ErrSubgraphSchemaChanged,
				fmt.Sprintf("The schema for subgraph '%s' has changed and was published.", data.Name.ValueString()),
			)
		}
	}

	var unsetLabels *bool
	if data.UnsetLabels.ValueBool() {
		unsetLabels = &[]bool{true}[0]
	}

	// TBD: This is only used in the update subgraph method and not used atm
	// headers := utils.ConvertHeadersToStringList(data.Headers)
	apiErr := r.client.UpdateSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.RoutingURL.ValueString(), labels, []string{}, data.SubscriptionUrl.ValueStringPointer(), data.Readme.ValueStringPointer(), unsetLabels, data.SubscriptionProtocol.ValueString(), data.WebsocketSubprotocol.ValueString())
	if apiErr != nil {
		if api.IsSubgraphCompositionFailedError(apiErr) {
			utils.AddDiagnosticWarning(resp,
				ErrUpdatingSubgraph,
				fmt.Sprintf("Could not update subgraph '%s': %s", data.Name.ValueString(), apiErr.Error()),
			)
		} else {
			utils.AddDiagnosticError(resp,
				ErrUpdatingSubgraph,
				fmt.Sprintf("Could not update subgraph '%s': %s", data.Name.ValueString(), apiErr.Error()),
			)
			return
		}
	}

	subgraph, err := r.client.GetSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrRetrievingSubgraph,
			fmt.Sprintf("Could not fetch updated subgraph '%s': %s", data.Name.ValueString(), err.Error()),
		)
		return
	}

	data.Id = types.StringValue(subgraph.GetId())
	data.Name = types.StringValue(subgraph.GetName())
	data.Namespace = types.StringValue(subgraph.GetNamespace())
	data.RoutingURL = types.StringValue(subgraph.GetRoutingURL())

	utils.LogAction(ctx, "updated", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubgraphResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SubgraphResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrDeletingSubgraph,
			fmt.Sprintf("Could not delete subgraph '%s': %s", data.Name.ValueString(), err.Error()),
		)
		return
	}

	utils.LogAction(ctx, "deleted", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())
}

func (r *SubgraphResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SubgraphResource) createAndPublishSubgraph(ctx context.Context, data SubgraphResourceModel, resp *resource.CreateResponse) (*platformv1.Subgraph, error) {
	var labels []*platformv1.Label
	for key, value := range data.Labels.Elements() {
		if strValue, ok := value.(types.String); ok {
			labels = append(labels, &platformv1.Label{
				Key:   key,
				Value: strValue.ValueString(),
			})
		}
	}
	utils.DebugAction(ctx, "Labels", data.Name.ValueString(), data.Namespace.ValueString(), map[string]interface{}{"labels": labels})

	apiErr := r.client.CreateSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.RoutingURL.ValueString(), nil, labels, data.SubscriptionUrl.ValueStringPointer(), data.Readme.ValueStringPointer(), data.IsEventDrivenGraph.ValueBoolPointer(), data.IsFeatureSubgraph.ValueBoolPointer(), data.SubscriptionProtocol.ValueString(), data.WebsocketSubprotocol.ValueString())
	if apiErr != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingSubgraph,
			fmt.Sprintf("Could not create subgraph '%s': %s", data.Name.ValueString(), apiErr.Error()),
		)
		return nil, apiErr
	}

	subgraph, apiErr := r.client.GetSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiErr != nil {
		utils.AddDiagnosticError(resp,
			ErrRetrievingSubgraph,
			fmt.Sprintf("Could not fetch created subgraph '%s': %s", data.Name.ValueString(), apiErr.Error()),
		)
		return nil, apiErr
	}

	if data.Schema.ValueString() != "" {
		_, apiErr := r.client.PublishSubgraph(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.Schema.ValueString())
		if apiErr != nil {
			if api.IsSubgraphCompositionFailedError(apiErr) {
				utils.AddDiagnosticWarning(resp,
					ErrPublishingSubgraph,
					fmt.Sprintf("Could not publish subgraph '%s': %s", data.Name.ValueString(), apiErr.Error()),
				)
			} else {
				utils.AddDiagnosticError(resp,
					ErrPublishingSubgraph,
					fmt.Sprintf("Could not publish subgraph '%s': %s", data.Name.ValueString(), apiErr.Error()),
				)
				return nil, apiErr
			}
		}
	}

	return subgraph, nil
}
