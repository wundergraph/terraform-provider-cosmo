package clients

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

var _ datasource.DataSource = &ClientsDataSource{}

func NewClientsDataSource() datasource.DataSource {
	return &ClientsDataSource{}
}

type ClientsDataSource struct {
	client *api.PlatformClient
}

type ClientsDataSourceModel struct {
	FederatedGraphName types.String            `tfsdk:"federated_graph_name"`
	Namespace          types.String            `tfsdk:"namespace"`
	Clients            []ClientDataSourceModel `tfsdk:"clients"`
}

type ClientDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	CreatedAt                types.String `tfsdk:"created_at"`
	LastUpdatedAt            types.String `tfsdk:"last_updated_at"`
	CreatedBy                types.String `tfsdk:"created_by"`
	LastUpdatedBy            types.String `tfsdk:"last_updated_by"`
	PersistedOperationsCount types.Int64  `tfsdk:"persisted_operations_count"`
}

func (d *ClientsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clients"
}

func (d *ClientsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the clients registered on a federated graph. Clients group persisted operations and are created implicitly when operations are first published for them.",

		Attributes: map[string]schema.Attribute{
			"federated_graph_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the federated graph.",
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace of the federated graph. Defaults to `default`.",
				Optional:            true,
				Computed:            true,
			},
			"clients": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The clients registered on the federated graph.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the client.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the client.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The time the client was created.",
						},
						"last_updated_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The time the client was last updated.",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user that created the client.",
						},
						"last_updated_by": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user that last updated the client.",
						},
						"persisted_operations_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of persisted operations registered for the client.",
						},
					},
				},
			},
		},
	}
}

func (d *ClientsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PlatformClient)
	if !ok {
		utils.AddDiagnosticError(resp, ErrUnexpectedDataSourceType, fmt.Sprintf("Expected *api.PlatformClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	d.client = client
}

func (d *ClientsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClientsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Namespace = utils.NamespaceOrDefault(data.Namespace)

	clients, apiError := d.client.GetClients(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		utils.AddDiagnosticError(resp, ErrReadingClients, apiError.Error())
		return
	}

	data.Clients = make([]ClientDataSourceModel, 0, len(clients))
	for _, client := range clients {
		model := ClientDataSourceModel{
			Id:            types.StringValue(client.Id),
			Name:          types.StringValue(client.Name),
			CreatedAt:     types.StringValue(client.CreatedAt),
			LastUpdatedAt: types.StringValue(client.LastUpdatedAt),
			CreatedBy:     types.StringValue(client.CreatedBy),
			LastUpdatedBy: types.StringValue(client.LastUpdatedBy),
		}
		if client.PersistedOperationsCount != nil {
			model.PersistedOperationsCount = types.Int64Value(int64(*client.PersistedOperationsCount))
		} else {
			model.PersistedOperationsCount = types.Int64Null()
		}
		data.Clients = append(data.Clients, model)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
