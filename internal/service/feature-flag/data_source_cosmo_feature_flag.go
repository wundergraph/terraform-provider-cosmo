package feature_flag

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

var _ datasource.DataSourceWithConfigure = &FeatureFlagDataSource{}

type FeatureFlagDataSource struct {
	client *api.PlatformClient
}

func NewFeatureFlagDataSource() datasource.DataSource {
	return &FeatureFlagDataSource{}
}

type FeatureFlagDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Namespace        types.String `tfsdk:"namespace"`
	FeatureSubgraphs types.Set    `tfsdk:"feature_subgraphs"`
	Labels           types.Map    `tfsdk:"labels"`
	IsEnabled        types.Bool   `tfsdk:"is_enabled"`
	CreatedBy        types.String `tfsdk:"created_by"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

func (d *FeatureFlagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cosmo Feature Flag Data Source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the feature flag.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the feature flag.",
				Required:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace of the feature flag. Defaults to the default namespace.",
				Optional:            true,
			},
			"feature_subgraphs": schema.SetAttribute{
				MarkdownDescription: "The list of feature subgraphs associated with the feature flag.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: `The labels associated with the feature flag. These labels indicate which 
federated graphs can be associated with the feature flag to enabled calls against the corresponding feature subgraph.`,
				ElementType: types.StringType,
				Computed:    true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the feature flag is enabled.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The user who created the feature flag.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the feature flag was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the feature flag was last updated.",
				Computed:            true,
			},
		},
	}
}

func (d *FeatureFlagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FeatureFlagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag"
}

func (d *FeatureFlagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FeatureFlagDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.ValueString() == "" {
		utils.AddDiagnosticError(resp,
			ErrInvalidFeatureFlagName,
			"The 'name' attribute is required.",
		)
		return
	}

	namespace := data.Namespace.ValueString()
	if namespace == "" {
		namespace = "default"
	}

	ff, apiError := d.client.GetFeatureFlag(ctx, data.Name.ValueString(), namespace)
	if apiError != nil {
		utils.AddDiagnosticError(resp,
			ErrRetrievingFeatureFlag,
			apiError.Error(),
		)
		return
	}

	labelMap := map[string]attr.Value{}

	for _, label := range ff.Labels {
		labelMap[label.Key] = types.StringValue(label.Value)
	}

	var diags diag.Diagnostics
	data.Labels, diags = types.MapValue(types.StringType, labelMap)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sgNames []attr.Value

	for _, name := range ff.FeatureSubgraphNames {
		sgNames = append(sgNames, types.StringValue(name))
	}

	data.FeatureSubgraphs, diags = types.SetValue(types.StringType, sgNames)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(ff.GetId())
	data.Name = types.StringValue(ff.GetName())
	data.Namespace = types.StringValue(ff.GetNamespace())
	data.IsEnabled = types.BoolValue(ff.GetIsEnabled())
	data.CreatedBy = types.StringValue(ff.GetCreatedBy())
	data.CreatedAt = types.StringValue(ff.GetCreatedAt())
	data.UpdatedAt = types.StringValue(ff.GetUpdatedAt())

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
