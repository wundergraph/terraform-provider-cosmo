package feature_subgraph

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

var _ datasource.DataSourceWithConfigure = &FeatureSubgraphDataSource{}

type FeatureSubgraphDataSource struct {
	client *api.PlatformClient
}

func NewFeatureSubgraphDataSource() datasource.DataSource {
	return &FeatureSubgraphDataSource{}
}

type FeatureSubgraphDataSourceModel struct {
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

func (d *FeatureSubgraphDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *FeatureSubgraphDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_subgraph"
}

func (d *FeatureSubgraphDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
This data source handles feature subgraphs. Feature subgraphs are a special type of subgraph that can be used to extend the functionality of the platform.
They require a base subgraph to be specified and can be used to add additional functionality to the base subgraph based on a specialized schema, that is published
to this subgraph.

For more information on feature subgraphs, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/feature-subgraph).
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the feature subgraph.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the feature subgraph.",
			},
			"namespace": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The namespace of the feature subgraph.",
			},
			"routing_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The routing URL of the feature subgraph.",
			},
			"base_subgraph_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the base subgraph that this feature subgraph extends.",
			},
			"subscription_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subscription URL for the subgraph.",
			},
			"subscription_protocol": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subscription protocol for the subgraph.",
			},
			"websocket_subprotocol": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The websocket subprotocol for the subgraph.",
			},
			"readme": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The readme for the subgraph.",
			},
			"schema": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The schema for the subgraph.",
			},
		},
	}
}

func (d *FeatureSubgraphDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FeatureSubgraphResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.ValueString() == "" {
		utils.AddDiagnosticError(resp,
			ErrReadingFeatureSubgraph,
			"The 'name' attribute is required.",
		)
		return
	}

	namespace := data.Namespace.ValueString()
	if namespace == "" {
		namespace = "default"
	}

	subgraph, apiErr := d.client.GetSubgraph(ctx, data.Name.ValueString(), namespace)
	if apiErr != nil {
		utils.AddDiagnosticError(resp, ErrRetrievingFeatureSubgraph, apiErr.Error())
		return
	}

	if !subgraph.IsFeatureSubgraph {
		utils.AddDiagnosticError(resp, ErrRetrievingFeatureSubgraph, "The subgraph is not a feature subgraph.")
		return
	}

	subgraphSchema, apiError := d.client.GetSubgraphSchema(ctx, subgraph.Name, subgraph.Namespace)
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
