package namespace

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"
)

var (
	_ resource.Resource                = (*NamespaceResource)(nil)
	_ resource.ResourceWithImportState = (*NamespaceResource)(nil)
)

type NamespaceResource struct {
	client *api.PlatformClient
}

type NamespaceResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewNamespaceResource() resource.Resource {
	return &NamespaceResource{}
}

func (r *NamespaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PlatformClient)
	if !ok {
		utils.AddDiagnosticError(resp, ErrUnexpectedDataSourceType, fmt.Sprintf("Expected *api.PlatformClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = client
}

func (r *NamespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *NamespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Namespaces group federated graphs and subgraphs. Each organization has a default, non-deletable namespace. 

For more information on namespaces, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/namespace).
		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the namespace resource.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the namespace.",
			},
		},
	}
}

func (r *NamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NamespaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.ValueString() == "" {
		utils.AddDiagnosticError(resp, ErrInvalidNamespaceName, "The 'name' attribute is required.")
		return
	}

	apiError := r.client.CreateNamespace(ctx, data.Name.ValueString())
	if apiError != nil {
		utils.AddDiagnosticError(resp,
			ErrCreatingNamespace,
			apiError.Error(),
		)
		return
	}

	namespace, err := getNamespace(ctx, *r.client, data.Id.ValueString(), data.Name.ValueString())
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrReadingNamespace,
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(namespace.Id)
	data.Name = types.StringValue(namespace.Name)

	utils.LogAction(ctx, "created", data.Id.ValueString(), data.Name.ValueString(), "")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	namespace, apiError := getNamespace(ctx, *r.client, data.Id.ValueString(), data.Name.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp, ErrReadingNamespace, apiError.Error())
		return
	}

	data.Id = types.StringValue(namespace.Id)
	data.Name = types.StringValue(namespace.Name)

	utils.LogAction(ctx, "read", data.Id.ValueString(), data.Name.ValueString(), "")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NamespaceResourceModel
	var state NamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.ValueString() != state.Name.ValueString() {
		utils.AddDiagnosticError(resp, ErrUpdatingNamespace, "Changing the namespace name requires recreation.")
		return
	}

	namespace, err := getNamespace(ctx, *r.client, data.Id.ValueString(), data.Name.ValueString())
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrReadingNamespace,
			err.Error(),
		)
		return
	}

	renameApiError := r.client.RenameNamespace(ctx, namespace.Name, data.Name.String())
	if renameApiError != nil {
		utils.AddDiagnosticError(resp,
			ErrUpdatingNamespace,
			renameApiError.Error(),
		)
		return
	}

	utils.LogAction(ctx, "updated", data.Id.ValueString(), data.Name.ValueString(), "")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NamespaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNamespace(ctx, data.Name.ValueString())
	if err != nil {
		utils.AddDiagnosticError(resp,
			ErrDeletingNamespace,
			err.Error(),
		)
		return
	}

	utils.LogAction(ctx, "deleted", data.Id.ValueString(), data.Name.ValueString(), "")
}

func getNamespace(ctx context.Context, client api.PlatformClient, id, name string) (*platformv1.Namespace, *api.ApiError) {
	namespace, err := client.GetNamespace(ctx, id, name)
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

func (r *NamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data NamespaceResourceModel

	id := req.ID

	if err := uuid.Validate(id); err != nil {
		data.Name = types.StringValue(id) // We assume this is the namespace name
	} else {
		data.Id = types.StringValue(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
