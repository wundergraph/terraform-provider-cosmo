package contract

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	SourceGraphName        types.String `tfsdk:"source"`
	Namespace              types.String `tfsdk:"namespace"`
	ExcludeTags            types.List   `tfsdk:"exclude_tags"`
	Readme                 types.String `tfsdk:"readme"`
	AdmissionWebhookUrl    types.String `tfsdk:"admission_webhook_url"`
	AdmissionWebhookSecret types.String `tfsdk:"admission_webhook_secret"`
	RoutingUrl             types.String `tfsdk:"routing_url"`
}

func (r *contractResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contract"
}

func (r *contractResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"source": schema.StringAttribute{
				Required: true,
			},
			"namespace": schema.StringAttribute{
				Required: true,
			},
			"exclude_tags": schema.ListAttribute{
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

	excludeTags, err := utils.ConvertLabelMatchers(data.ExcludeTags)
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingContract,
			"Could not create contract: "+err.Error(),
		)
		return
	}
	_, apiError := r.client.CreateContract(ctx, data.Name.ValueString(), data.Namespace.ValueString(), data.SourceGraphName.ValueString(), data.RoutingUrl.ValueString(), data.AdmissionWebhookUrl.ValueString(), data.AdmissionWebhookSecret.ValueString(), excludeTags, data.Readme.ValueString())
	if apiError != nil {
		if api.IsContractCompositionFailedError(apiError) || api.IsSubgraphCompositionFailedError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrCreatingContract,
				"Contract composition failed: "+apiError.Error(),
			)
		} else {
			utils.AddDiagnosticError(resp,
				ErrCreatingContract,
				"Could not create contract: "+apiError.Error(),
			)
			return
		}
	}

	data.ID = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *contractResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data contractResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

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
			"Could not update contract: "+err.Error(),
		)
		return
	}

	_, apiError := r.client.UpdateContract(ctx, data.Name.ValueString(), data.Namespace.ValueString(), excludeTags)
	if apiError != nil {
		if api.IsContractCompositionFailedError(apiError) || api.IsSubgraphCompositionFailedError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrUpdatingContract,
				"Contract composition failed: "+apiError.Error(),
			)
		} else {
			utils.AddDiagnosticError(resp,
				ErrUpdatingContract,
				"Could not update contract: "+apiError.Error(),
			)
			return
		}
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
		if api.IsContractCompositionFailedError(apiError) || api.IsSubgraphCompositionFailedError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrDeletingContract,
				"Contract composition failed: "+apiError.Error(),
			)
		} else if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticWarning(resp,
				ErrDeletingContract,
				"Contract composition failed: "+apiError.Error(),
			)
			resp.State.RemoveResource(ctx)
		} else {
			utils.AddDiagnosticError(resp,
				ErrDeletingContract,
				"Could not delete contract: "+apiError.Error(),
			)
			return
		}
	}
}

func (r *contractResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
