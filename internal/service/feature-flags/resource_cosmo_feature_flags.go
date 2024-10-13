package feature_flags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"

	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

type FeatureFlagResource struct {
	client *api.PlatformClient
}

type FeatureFlagResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Namespace            types.String `tfsdk:"namespace"`
	IsEnabled            types.Bool   `tfsdk:"is_enabled"`
	Labels               types.List   `tfsdk:"labels"`
	FeatureSubgraphNames types.List   `tfsdk:"feature_subgraph_names"`
}

func NewFeatureFlagResource() resource.Resource {
	return &FeatureFlagResource{}
}

func (r *FeatureFlagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cosmo_feature_flag"
}

func (r *FeatureFlagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"namespace": schema.StringAttribute{
				Required: true,
			},
			"is_enabled": schema.BoolAttribute{
				Required: true,
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"feature_subgraph_names": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *FeatureFlagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PlatformClient)
	if !ok {
		utils.AddDiagnosticError(resp, ErrUnexpectedDataSourceType, fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = client
}

func (r *FeatureFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeatureFlagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	_, apiError := r.client.CreateFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.IsEnabled.ValueBool(), convertLabels(data.Labels), convertFeatureSubgraphNames(data.FeatureSubgraphNames))
	if apiError != nil {
		resp.Diagnostics.AddError("Error creating feature flag", apiError.Error())
		return
	}

	getResponse, apiError := r.client.GetFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		resp.Diagnostics.AddError("Error retrieving created feature flag", apiError.Error())
		return
	}

	data.Id = types.StringValue(getResponse.FeatureFlag.GetId())
	data.IsEnabled = types.BoolValue(getResponse.FeatureFlag.GetIsEnabled())
	data.Labels = convertToList(getResponse.FeatureFlag.GetLabels())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeatureFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureFlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, apiError := r.client.GetFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		resp.Diagnostics.AddError("Error reading feature flag", apiError.Error())
		return
	}

	data.IsEnabled = types.BoolValue(response.FeatureFlag.GetIsEnabled())
	data.Labels = convertToList(response.FeatureFlag.GetLabels())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeatureFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FeatureFlagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, apiError := r.client.UpdateFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString(), false, convertLabels(data.Labels), convertFeatureSubgraphNames(data.FeatureSubgraphNames))
	if apiError != nil {
		resp.Diagnostics.AddError("Error updating feature flag", apiError.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeatureFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeatureFlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiError := r.client.DeleteFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiError != nil {
		resp.Diagnostics.AddError("Error deleting feature flag", apiError.Error())
		return
	}
}

func (r *FeatureFlagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertLabels(labels types.List) []*platformv1.Label {
	var result []*platformv1.Label
	for _, label := range labels.Elements() {
		result = append(result, &platformv1.Label{Key: label.(types.String).ValueString()})
	}
	return result
}

func convertFeatureSubgraphNames(names types.List) []string {
	var result []string
	for _, name := range names.Elements() {
		result = append(result, name.(types.String).ValueString())
	}
	return result
}

func convertToList(items []*platformv1.Label) types.List {
	var result []attr.Value
	for _, item := range items {
		result = append(result, types.StringValue(item.Key))
	}
	return types.ListValueMust(types.StringType, result)
}
