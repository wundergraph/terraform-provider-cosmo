package monograph

import (
	"context"
	"fmt"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

type MonographResource struct {
	client *api.PlatformClient
}

type MonographResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Namespace              types.String `tfsdk:"namespace"`
	SubscriptionUrl        types.String `tfsdk:"subscription_url"`
	WebsocketSubprotocol   types.String `tfsdk:"websocket_subprotocol"`
	SubscriptionProtocol   types.String `tfsdk:"subscription_protocol"`
	GraphUrl               types.String `tfsdk:"graph_url"`
	RoutingURL             types.String `tfsdk:"routing_url"`
	Readme                 types.String `tfsdk:"readme"`
	AdmissionWebhookURL    types.String `tfsdk:"admission_webhook_url"`
	AdmissionWebhookSecret types.String `tfsdk:"admission_webhook_secret"`
	Schema                 types.String `tfsdk:"schema"`
}

func NewMonographResource() resource.Resource {
	return &MonographResource{}
}

func (r *MonographResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monograph"
}

func (r *MonographResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
A monograph is a resource that represents a single subgraph with GraphQL Federation disabled.

For more information on monographs, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/monograph).
		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the monograph resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the monograph.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace in which the monograph is located.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"graph_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GraphQL endpoint URL of the monograph.",
			},
			"websocket_subprotocol": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The websocket subprotocol for the subgraph.",
				Validators: []validator.String{
					stringvalidator.OneOf(api.GraphQLWebsocketSubprotocolDefault, api.GraphQLWebsocketSubprotocolGraphQLWS, api.GraphQLWebsocketSubprotocolGraphQLTransportWS),
				},
			},
			"routing_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The routing URL for the monograph.",
			},
			"readme": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The readme for the subgraph.",
			},
			"admission_webhook_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The admission webhook URL for the monograph.",
			},
			"admission_webhook_secret": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The admission webhook secret for the monograph.",
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
			"schema": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The schema for the subgraph.",
			},
		},
	}
}

func (r *MonographResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonographResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonographResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.ValueString() == "" {
		utils.AddDiagnosticError(resp,
			ErrInvalidMonographName,
			"The 'name' attribute is required.",
		)
		return
	}

	_, apiError := r.client.CreateMonograph(
		ctx,
		data.Name.ValueString(),
		data.Namespace.ValueString(),
		data.RoutingURL.ValueString(),
		data.GraphUrl.ValueString(),
		utils.StringValueOrNil(data.SubscriptionUrl),
		utils.StringValueOrNil(data.Readme),
		data.WebsocketSubprotocol.ValueString(),
		data.SubscriptionProtocol.ValueString(),
		data.AdmissionWebhookURL.ValueString(),
		data.AdmissionWebhookSecret.ValueString(),
	)
	if apiError != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingMonograph,
			apiError.Error(),
		)
		return
	}

	if data.Schema.ValueString() != "" {
		err := r.client.PublishMonograph(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.Schema.ValueString())
		if err != nil {
			if api.IsNotFoundError(err) {
				utils.AddDiagnosticError(resp,
					ErrPublishingMonograph,
					err.Error(),
				)
				resp.State.RemoveResource(ctx)
				return
			} else {
				utils.AddDiagnosticError(resp,
					ErrPublishingMonograph,
					err.Error(),
				)
				return
			}
		}
	}

	monograph, apiError := r.client.GetMonograph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		utils.AddDiagnosticError(resp,
			ErrRetrievingMonograph,
			apiError.Error(),
		)
		return
	}

	data.Id = types.StringValue(monograph.GetId())
	if monograph.Readme != nil {
		data.Readme = types.StringValue(*monograph.Readme)
	}

	utils.LogAction(ctx, "created monograph", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonographResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MonographResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var monograph *platformv1.FederatedGraph
	if data.Name.ValueString() == "" {
		graph, apiError := r.client.GetMonographByID(ctx, data.Id.ValueString())
		if apiError != nil {
			if api.IsNotFoundError(apiError) {
				utils.AddDiagnosticWarning(resp,
					ErrMonographNotFound,
					apiError.Error(),
				)
				resp.State.RemoveResource(ctx)
				return
			}
			utils.AddDiagnosticError(resp, ErrReadingMonograph, apiError.Error())
			return
		}
		monograph = graph
	} else {
		graph, apiError := r.client.GetMonograph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
		if apiError != nil {
			if api.IsNotFoundError(apiError) {
				utils.AddDiagnosticWarning(resp,
					ErrMonographNotFound,
					apiError.Error(),
				)
				resp.State.RemoveResource(ctx)
				return
			}
			utils.AddDiagnosticError(resp,
				ErrRetrievingMonograph,
				apiError.Error(),
			)
			return
		}

		monograph = graph
	}

	subGraph, err := r.client.GetSubgraph(ctx, monograph.GetName(), monograph.GetNamespace())
	if err != nil {
		if api.IsNotFoundError(err) {
			utils.AddDiagnosticError(resp,
				ErrRetrievingMonograph,
				err.Error(),
			)
			return
		} else {
			utils.AddDiagnosticError(resp,
				ErrRetrievingMonograph,
				err.Error(),
			)
			return
		}
	}

	data.Id = types.StringValue(monograph.GetId())
	data.Name = types.StringValue(monograph.GetName())
	data.Namespace = types.StringValue(monograph.GetNamespace())
	data.RoutingURL = types.StringValue(monograph.GetRoutingURL())
	data.GraphUrl = types.StringValue(subGraph.GetRoutingURL())

	if monograph.Readme != nil {
		data.Readme = types.StringValue(*monograph.Readme)
	}

	utils.LogAction(ctx, "read monograph", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonographResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MonographResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateMonograph(
		ctx,
		data.Name.ValueString(),
		data.Namespace.ValueString(),
		data.RoutingURL.ValueString(),
		data.GraphUrl.ValueString(),
		utils.StringValueOrNil(data.SubscriptionUrl),
		utils.StringValueOrNil(data.Readme),
		data.WebsocketSubprotocol.ValueString(),
		data.SubscriptionProtocol.ValueString(),
		data.AdmissionWebhookURL.ValueString(),
		data.AdmissionWebhookSecret.ValueString(),
	)
	if err != nil {
		if api.IsNotFoundError(err) {
			utils.AddDiagnosticError(resp,
				ErrUpdatingMonograph,
				err.Error(),
			)
			resp.State.RemoveResource(ctx)
			return
		} else {
			utils.AddDiagnosticError(resp,
				ErrUpdatingMonograph,
				err.Error(),
			)
			return
		}
	}

	if data.Schema.ValueString() != "" {
		err := r.client.PublishMonograph(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.Schema.ValueString())
		if err != nil {
			if api.IsNotFoundError(err) {
				utils.AddDiagnosticError(resp,
					ErrUpdatingMonograph,
					err.Error(),
				)
				resp.State.RemoveResource(ctx)
				return
			} else {
				utils.AddDiagnosticError(resp,
					ErrUpdatingMonograph,
					err.Error(),
				)
				return
			}
		}
	}

	utils.LogAction(ctx, "updated monograph", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())

	monograph, err := r.client.GetMonograph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrRetrievingMonograph,
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(monograph.GetId())
	data.Name = types.StringValue(monograph.GetName())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonographResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MonographResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiError := r.client.DeleteMonograph(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticError(resp,
				ErrDeletingMonograph,
				apiError.Error(),
			)
			resp.State.RemoveResource(ctx)
		} else {
			utils.AddDiagnosticError(resp,
				ErrDeletingMonograph,
				apiError.Error(),
			)
			return
		}
	}

	utils.LogAction(ctx, "deleted monograph", data.Id.ValueString(), data.Name.ValueString(), data.Namespace.ValueString())
}

func (r *MonographResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
