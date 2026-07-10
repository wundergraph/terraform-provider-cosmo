package persisted_operations

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
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

var _ resource.Resource = &PersistedOperationsResource{}
var _ resource.ResourceWithImportState = &PersistedOperationsResource{}

type PersistedOperationsResource struct {
	client *api.PlatformClient
}

type PersistedOperationsResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	FederatedGraphName types.String `tfsdk:"federated_graph_name"`
	Namespace          types.String `tfsdk:"namespace"`
	ClientName         types.String `tfsdk:"client_name"`
	Operations         types.Map    `tfsdk:"operations"`
}

func NewPersistedOperationsResource() resource.Resource {
	return &PersistedOperationsResource{}
}

func (r *PersistedOperationsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_persisted_operations"
}

func (r *PersistedOperationsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Manages the full set of persisted operations (safelisted GraphQL operations) registered for a client of a federated graph.

This resource owns every persisted operation of the given client: operations pushed outside of Terraform (e.g. via ` + "`wgc operations push`" + `) to the same client will show up as drift and be removed on the next apply. Destroying the resource deletes the operations it manages but leaves the client record itself in place.

Operations are content-addressed: their identifier is the hex-encoded SHA-256 hash of their contents, matching the behavior of ` + "`wgc operations push`" + `. Publishing is not transactional: if an operation conflicts with an existing one (same id, different contents), the create or update fails but the non-conflicting operations remain registered; re-applying after fixing the conflict converges, since publishing is idempotent.

Existing operations can be imported with the id format ` + "`federated_graph_name:namespace:client_name`" + `.

For more information on persisted operations, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/router/persisted-queries/persisted-operations).
		`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the client the operations are registered for.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"federated_graph_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the federated graph to register the operations on.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace of the federated graph.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the client the operations belong to (e.g. `web`, `ios`). The client is created on first publish.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"operations": schema.MapAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The persisted operations as a map of an arbitrary label to the GraphQL operation contents.",
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (r *PersistedOperationsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func operationID(contents string) string {
	hash := sha256.Sum256([]byte(contents))
	return hex.EncodeToString(hash[:])
}

// buildOperations converts the label -> contents map into persisted operations
// keyed by the SHA-256 hash of their contents. It errors when two labels hold
// identical contents, since both would map to the same operation id.
func buildOperations(operations map[string]string) ([]*platformv1.PersistedOperation, error) {
	labelsById := make(map[string]string, len(operations))
	result := make([]*platformv1.PersistedOperation, 0, len(operations))

	for label, contents := range operations {
		id := operationID(contents)
		if existing, ok := labelsById[id]; ok {
			return nil, fmt.Errorf("operations %q and %q have identical contents and would map to the same operation id %s", existing, label, id)
		}
		labelsById[id] = label
		result = append(result, &platformv1.PersistedOperation{
			Id:       id,
			Contents: contents,
		})
	}

	return result, nil
}

func (r *PersistedOperationsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PersistedOperationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var planOperations map[string]string
	resp.Diagnostics.Append(data.Operations.ElementsAs(ctx, &planOperations, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operations, err := buildOperations(planOperations)
	if err != nil {
		utils.AddDiagnosticError(resp, ErrDuplicatePersistedOperation, err.Error())
		return
	}

	if apiError := r.client.PublishPersistedOperations(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), data.ClientName.ValueString(), operations); apiError != nil {
		utils.AddDiagnosticError(resp, ErrCreatingPersistedOperations, apiError.Error())
		return
	}

	client, apiError := r.client.GetClient(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), data.ClientName.ValueString())
	if apiError != nil {
		utils.AddDiagnosticError(resp, ErrCreatingPersistedOperations, apiError.Error())
		return
	}

	data.Id = types.StringValue(client.Id)

	utils.LogAction(ctx, "created persisted operations", data.Id.ValueString(), data.ClientName.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistedOperationsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PersistedOperationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, apiError := r.client.GetClient(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), data.ClientName.ValueString())
	if apiError != nil {
		if api.IsNotFoundError(apiError) {
			utils.AddDiagnosticWarning(resp, ErrReadingPersistedOperations, apiError.Error())
			resp.State.RemoveResource(ctx)
			return
		}
		utils.AddDiagnosticError(resp, ErrReadingPersistedOperations, apiError.Error())
		return
	}

	remoteOperations, apiError := r.client.GetPersistedOperations(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), client.Id)
	if apiError != nil {
		utils.AddDiagnosticError(resp, ErrReadingPersistedOperations, apiError.Error())
		return
	}

	// An empty remote set is drift, not deletion: the client record is the
	// remote object this resource tracks, and its absence is handled above.
	remoteById := make(map[string]*platformv1.GetPersistedOperationsResponse_Operation, len(remoteOperations))
	for _, operation := range remoteOperations {
		remoteById[operation.Id] = operation
	}

	// Operations are null when the resource is being imported.
	var stateOperations map[string]string
	if !data.Operations.IsNull() {
		resp.Diagnostics.Append(data.Operations.ElementsAs(ctx, &stateOperations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Keep state entries whose contents still exist remotely (matched by id, so
	// server-side normalization of contents does not cause spurious diffs) and
	// surface out-of-band operations under their id so the plan shows their removal.
	operations := make(map[string]string, len(remoteById))
	for label, contents := range stateOperations {
		id := operationID(contents)
		if _, ok := remoteById[id]; ok {
			operations[label] = contents
			delete(remoteById, id)
		}
	}
	for id, operation := range remoteById {
		operations[id] = operation.Contents
	}

	operationsValue, diags := types.MapValueFrom(ctx, types.StringType, operations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(client.Id)
	data.Operations = operationsValue

	utils.LogAction(ctx, "read persisted operations", data.Id.ValueString(), data.ClientName.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistedOperationsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PersistedOperationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	var state PersistedOperationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var planOperations map[string]string
	resp.Diagnostics.Append(data.Operations.ElementsAs(ctx, &planOperations, false)...)

	var stateOperations map[string]string
	resp.Diagnostics.Append(state.Operations.ElementsAs(ctx, &stateOperations, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	operations, err := buildOperations(planOperations)
	if err != nil {
		utils.AddDiagnosticError(resp, ErrDuplicatePersistedOperation, err.Error())
		return
	}

	// Publishing is idempotent for unchanged operations, so all planned
	// operations are pushed in one call rather than only the additions.
	if apiError := r.client.PublishPersistedOperations(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), data.ClientName.ValueString(), operations); apiError != nil {
		utils.AddDiagnosticError(resp, ErrUpdatingPersistedOperations, apiError.Error())
		return
	}

	planContents := make(map[string]bool, len(planOperations))
	for _, contents := range planOperations {
		planContents[contents] = true
	}

	for _, contents := range stateOperations {
		// Identical contents produce identical operation ids, so diffing on
		// contents avoids re-hashing the unchanged operations.
		if planContents[contents] {
			continue
		}
		apiError := r.client.DeletePersistedOperation(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), data.ClientName.ValueString(), operationID(contents))
		if apiError != nil && !api.IsNotFoundError(apiError) {
			utils.AddDiagnosticError(resp, ErrUpdatingPersistedOperations, apiError.Error())
			return
		}
	}

	utils.LogAction(ctx, "updated persisted operations", data.Id.ValueString(), data.ClientName.ValueString(), data.Namespace.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PersistedOperationsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PersistedOperationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var stateOperations map[string]string
	resp.Diagnostics.Append(data.Operations.ElementsAs(ctx, &stateOperations, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, contents := range stateOperations {
		apiError := r.client.DeletePersistedOperation(ctx, data.FederatedGraphName.ValueString(), data.Namespace.ValueString(), data.ClientName.ValueString(), operationID(contents))
		if apiError != nil && !api.IsNotFoundError(apiError) {
			utils.AddDiagnosticError(resp, ErrDeletingPersistedOperations, apiError.Error())
			return
		}
	}

	utils.LogAction(ctx, "deleted persisted operations", data.Id.ValueString(), data.ClientName.ValueString(), data.Namespace.ValueString())
}

func (r *PersistedOperationsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		utils.AddDiagnosticError(resp, ErrInvalidImportID, fmt.Sprintf("expected import id in the format 'federated_graph_name:namespace:client_name', got: %q", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("federated_graph_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("namespace"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("client_name"), parts[2])...)
}
