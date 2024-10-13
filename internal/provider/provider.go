package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/utils"

	contract "github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/service/contract"
	feature_flags "github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/service/feature-flags"
	federated_graph "github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/service/federated-graph"
	monograph "github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/service/monograph"
	namespace "github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/service/namespace"
	router_token "github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/service/router-token"
	subgraph "github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/service/subgraph"
)

// Ensure CosmoProvider satisfies various provider interfaces.
var _ provider.Provider = &CosmoProvider{}
var _ provider.ProviderWithFunctions = &CosmoProvider{}

// CosmoProvider defines the provider implementation.
type CosmoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type Provider struct {
	*api.PlatformClient
}

// CosmoProviderModel describes the provider data model.
type CosmoProviderModel struct {
	ApiUrl types.String `tfsdk:"api_url"`
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *CosmoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cosmo"
	resp.Version = p.version
}

func (p *CosmoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
The Cosmo provider allows you to interact with WunderGraph's Cosmo API, managing key resources. 
It supports creating and reading namespaces, federated graphs, subgraphs, router tokens, monographs, and contracts. 

Refer to the official [Cosmo Documentation](https://cosmo-docs.wundergraph.com/) for more details.
		`,
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The Api Url to be used: Leave blank to use: https://cosmo-cp.wundergraph.com or use the %s environment variable", utils.EnvCosmoApiUrl),
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The Api Key to be used: Leave blank to use the %s environment variable", utils.EnvCosmoApiKey),
				Optional:            true,
			},
		},
	}
}

func (p *CosmoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CosmoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cosmoApiKey := data.ApiKey.ValueString()
	cosmoApiUrl := data.ApiUrl.ValueString()

	platformClient, err := api.NewClient(cosmoApiKey, cosmoApiUrl)

	if err != nil {
		utils.AddDiagnosticError(resp, "Error configuring client", err.Error())
		return
	}
	resp.DataSourceData = platformClient
	resp.ResourceData = platformClient
}

func (p *CosmoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		federated_graph.NewFederatedGraphResource,
		namespace.NewNamespaceResource,
		subgraph.NewSubgraphResource,
		monograph.NewMonographResource,
		router_token.NewTokenResource,
		contract.NewContractResource,
		feature_flags.NewFeatureFlagResource,
	}
}

func (p *CosmoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		federated_graph.NewFederatedGraphDataSource,
		subgraph.NewSubgraphDataSource,
		namespace.NewNamespaceDataSource,
		monograph.NewMonographDataSource,
		contract.NewContractDataSource,
		feature_flags.NewFeatureFlagDataSource,
	}
}

func (p *CosmoProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CosmoProvider{
			version: version,
		}
	}
}
