package feature_flag

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

var _ interface {
	resource.ResourceWithConfigure
	resource.ResourceWithImportState
} = &FeatureFlagResource{}

type FeatureFlagResource struct {
	client *api.PlatformClient
}

func NewFeatureFlagResource() resource.Resource {
	return &FeatureFlagResource{}
}

type FeatureFlagResourceModel struct {
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

func (r *FeatureFlagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_flag"
}

func (r *FeatureFlagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `A feature flag is a group of one or more feature subgraphs. 
Each feature subgraph represents a replacement of a specific base subgraph that composes a federated graph. 
Feature flags define labels that dictate the federated graphs to which they will apply when enabled. 
Setting the corresponding feature-flag header or cookie value allows different graph constellations to be served to clients.

For more information on feature flags, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/feature-flags).
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the feature flag.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the feature flag.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace of the feature flag.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"feature_subgraphs": schema.SetAttribute{
				MarkdownDescription: `The list of feature subgraphs associated with the feature flag. 
At least one feature subgraph must be provided.`,
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: `The labels associated with the feature flag. These labels indicate which 
federated graphs can be associated with the feature flag to enabled calls against the corresponding feature subgraph.`,
				ElementType: types.StringType,
				Optional:    true,
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the feature flag is enabled.",
				Optional:            true,
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

func (r *FeatureFlagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FeatureFlagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeatureFlagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var featureSubgraphs []string
	for _, val := range data.FeatureSubgraphs.Elements() {
		if strVal, ok := val.(types.String); ok {
			featureSubgraphs = append(featureSubgraphs, strVal.ValueString())
		}
	}

	ffLabels := make([]*platformv1.Label, 0, len(data.Labels.Elements()))
	for key, value := range data.Labels.Elements() {
		if strVal, ok := value.(types.String); ok {
			ffLabels = append(ffLabels, &platformv1.Label{
				Key:   key,
				Value: strVal.ValueString(),
			})
		}
	}

	apiErr := r.client.CreateFeatureFlag(ctx, &api.FeatureFlag{
		FeatureFlag: &platformv1.FeatureFlag{
			Name:      data.Name.ValueString(),
			Namespace: data.Namespace.ValueString(),
			Labels:    ffLabels,
			IsEnabled: data.IsEnabled.ValueBool(),
		},
		FeatureSubgraphNames: featureSubgraphs,
	})

	if apiErr != nil {
		utils.AddDiagnosticError(resp, ErrFeatureFlagCreate, apiErr.Error())
		return
	}

	ff, apiErr := r.client.GetFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiErr != nil {
		if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticWarning(resp, ErrFeatureFlagGet, fmt.Sprintf("Feature flag %s not found: %s", data.Name, apiErr.Error()))
			resp.State.RemoveResource(ctx)
			return
		}

		utils.AddDiagnosticError(resp, ErrFeatureFlagGet, fmt.Sprintf("Failed to retrieve created feature flag after creation: %s", apiErr.Error()))
		return
	}

	resp.Diagnostics.Append(mapFeatureFlagToResourceModel(ff, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FeatureFlagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeatureFlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.ValueString() == "" || data.Namespace.ValueString() == "" {
		utils.AddDiagnosticError(resp, ErrInvalidFeatureFlagName, "The 'name' and 'namespace' attributes are required.")
		return
	}

	ff, apiErr := r.client.GetFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiErr != nil {
		if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticWarning(resp, ErrFeatureFlagGet, fmt.Sprintf("Feature flag %s not found: %s", data.Name, apiErr.Error()))
			resp.State.RemoveResource(ctx)
			return
		}

		utils.AddDiagnosticError(resp, ErrFeatureFlagGet, apiErr.Error())
		return
	}

	resp.Diagnostics.Append(mapFeatureFlagToResourceModel(ff, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "Read feature flag resource", map[string]interface{}{
		"name":      data.Name.ValueString(),
		"namespace": data.Namespace.ValueString(),
	})
}

func (r *FeatureFlagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FeatureFlagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ffLabels := make([]*platformv1.Label, 0, len(data.Labels.Elements()))

	for key, value := range data.Labels.Elements() {
		if strVal, ok := value.(types.String); ok {
			ffLabels = append(ffLabels, &platformv1.Label{
				Key:   key,
				Value: strVal.ValueString(),
			})
		}
	}

	var featureSubgraphNames []string
	for _, name := range data.FeatureSubgraphs.Elements() {
		if strVal, ok := name.(types.String); ok {
			featureSubgraphNames = append(featureSubgraphNames, strVal.ValueString())
		}
	}

	apiErr := r.client.UpdateFeatureFlag(ctx, &api.FeatureFlag{
		FeatureFlag: &platformv1.FeatureFlag{
			Name:      data.Name.ValueString(),
			Namespace: data.Namespace.ValueString(),
			Labels:    ffLabels,
		},
		FeatureSubgraphNames: featureSubgraphNames,
	})

	if apiErr != nil {
		if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticWarning(resp, ErrFeatureFlagUpdate, apiErr.Error())
			resp.State.RemoveResource(ctx)
			return
		}

		utils.AddDiagnosticError(resp, ErrFeatureFlagUpdate, apiErr.Error())
		return
	}

	ff, apiErr := r.client.GetFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiErr != nil {
		if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticWarning(resp, ErrFeatureFlagGet, fmt.Sprintf("Feature flag %s not found: %s", data.Name, apiErr.Error()))
			resp.State.RemoveResource(ctx)
			return
		}

		utils.AddDiagnosticError(resp, ErrFeatureFlagGet, fmt.Sprintf("Failed to retrieve created feature flag after creation: %s", apiErr.Error()))
		return
	}

	if ff.IsEnabled != data.IsEnabled.ValueBool() {
		apiErr = r.client.SetFeatureFlagState(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.IsEnabled.ValueBool())
		if apiErr != nil {
			if api.IsNotFoundError(apiErr) {
				utils.AddDiagnosticWarning(resp, ErrFeatureFlagUpdate, apiErr.Error())
				resp.State.RemoveResource(ctx)
				return
			}

			utils.AddDiagnosticError(resp, ErrFeatureFlagUpdate, apiErr.Error())
			return
		}

		ff.IsEnabled = data.IsEnabled.ValueBool()
	}

	resp.Diagnostics.Append(mapFeatureFlagToResourceModel(ff, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "Updated feature flag resource", map[string]interface{}{
		"name":      data.Name.ValueString(),
		"namespace": data.Namespace.ValueString(),
	})
}

func (r *FeatureFlagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeatureFlagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiErr := r.client.DeleteFeatureFlag(ctx, data.Name.ValueString(), data.Namespace.ValueString())
	if apiErr != nil {

		if api.IsNotFoundError(apiErr) {
			utils.AddDiagnosticWarning(resp, ErrFeatureFlagDelete, apiErr.Error())
			resp.State.RemoveResource(ctx)
			return
		}

		utils.AddDiagnosticError(resp, ErrFeatureFlagDelete, apiErr.Error())
		return
	}

	tflog.Trace(ctx, "Deleted feature flag resource", map[string]interface{}{
		"name":      data.Name.ValueString(),
		"namespace": data.Namespace.ValueString(),
	})
}

func (r *FeatureFlagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data FeatureFlagResourceModel

	// For feature flags we cannot get the UUID, therefore we need to use a format '<namespace>.<name>'
	// '.' is not a valid token for a namespace name, so we can use it as a separator
	id := req.ID

	before, after, found := strings.Cut(id, ".")
	if !found {
		utils.AddDiagnosticError(resp, ErrInvalidImportState, "The ID must be in the format '<namespace>.<name>'")
		return
	}

	data.Name = types.StringValue(after)
	data.Namespace = types.StringValue(before)

	// We can't use SetState here as we need to initialize sub-typed fields like Set and Map
	// Therefore we need to set the individual fields
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), data.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("namespace"), data.Namespace)...)
}

func mapFeatureFlagToResourceModel(ff *api.FeatureFlag, res *FeatureFlagResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var tt []attr.Value

	for _, subgraphName := range ff.FeatureSubgraphNames {
		tt = append(tt, types.StringValue(subgraphName))
	}

	fsg, d := types.SetValue(types.StringType, tt)
	diags = append(diags, d...)

	labelTypeMap := make(map[string]attr.Value)

	for _, l := range ff.Labels {
		labelTypeMap[l.Key] = types.StringValue(l.Value)
	}

	labels, d := types.MapValue(types.StringType, labelTypeMap)
	diags = append(diags, d...)

	res.Id = types.StringValue(ff.Id)
	res.Name = types.StringValue(ff.Name)
	res.Namespace = types.StringValue(ff.Namespace)
	res.FeatureSubgraphs = fsg
	res.Labels = labels
	res.IsEnabled = types.BoolValue(ff.IsEnabled)
	res.CreatedBy = types.StringValue(ff.CreatedBy)
	res.CreatedAt = types.StringValue(ff.CreatedAt)
	res.UpdatedAt = types.StringValue(ff.UpdatedAt)

	return diags
}
