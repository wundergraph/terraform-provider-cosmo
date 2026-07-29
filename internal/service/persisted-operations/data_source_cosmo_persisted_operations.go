package persisted_operations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

var _ datasource.DataSource = &PersistedOperationsDataSource{}

func NewPersistedOperationsDataSource() datasource.DataSource {
	return &PersistedOperationsDataSource{}
}

type PersistedOperationsDataSource struct {
	client *api.PlatformClient
}

type PersistedOperationsDataSourceModel struct {
	FederatedGraphName types.String                        `tfsdk:"federated_graph_name"`
	Namespace          types.String                        `tfsdk:"namespace"`
	ClientName         types.String                        `tfsdk:"client_name"`
	ClientId           types.String                        `tfsdk:"client_id"`
	Operations         []PersistedOperationDataSourceModel `tfsdk:"operations"`
}

type PersistedOperationDataSourceModel struct {
	Id             types.String   `tfsdk:"id"`
	Contents       types.String   `tfsdk:"contents"`
	CreatedAt      types.String   `tfsdk:"created_at"`
	LastUpdatedAt  types.String   `tfsdk:"last_updated_at"`
	OperationNames []types.String `tfsdk:"operation_names"`
}

func (d *PersistedOperationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_persisted_operations"
}

func (d *PersistedOperationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the persisted operations registered for a client of a federated graph.",

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
			"client_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the client the operations belong to.",
			},
			"client_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the client.",
			},
			"operations": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The persisted operations registered for the client.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the operation (the SHA-256 hash of its contents).",
						},
						"contents": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GraphQL operation contents.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The time the operation was created.",
						},
						"last_updated_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The time the operation was last updated.",
						},
						"operation_names": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "The names of the operations contained in the document.",
						},
					},
				},
			},
		},
	}
}

func (d *PersistedOperationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PersistedOperationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PersistedOperationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Namespace = utils.NamespaceOrDefault(data.Namespace)

	client, apiError := d.client.GetClient(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), data.ClientName.ValueString())
	if apiError != nil {
		utils.AddDiagnosticError(resp, ErrReadingPersistedOperations, apiError.Error())
		return
	}

	operations, apiError := d.client.GetPersistedOperations(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), client.Id)
	if apiError != nil {
		utils.AddDiagnosticError(resp, ErrReadingPersistedOperations, apiError.Error())
		return
	}

	data.ClientId = types.StringValue(client.Id)
	data.Operations = make([]PersistedOperationDataSourceModel, 0, len(operations))
	for _, operation := range operations {
		operationNames := make([]types.String, 0, len(operation.OperationNames))
		for _, name := range operation.OperationNames {
			operationNames = append(operationNames, types.StringValue(name))
		}
		data.Operations = append(data.Operations, PersistedOperationDataSourceModel{
			Id:             types.StringValue(operation.Id),
			Contents:       types.StringValue(operation.Contents),
			CreatedAt:      types.StringValue(operation.CreatedAt),
			LastUpdatedAt:  types.StringValue(operation.LastUpdatedAt),
			OperationNames: operationNames,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
