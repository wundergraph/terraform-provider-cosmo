package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p PlatformClient) CreateSubgraph(ctx context.Context, data *platformv1.CreateFederatedSubgraphRequest) *ApiError {
	request := connect.NewRequest(data)
	response, err := p.Client.CreateFederatedSubgraph(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "CreateSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "CreateSubgraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p PlatformClient) UpdateSubgraph(ctx context.Context, data *platformv1.UpdateSubgraphRequest) *ApiError {
	request := connect.NewRequest(data)

	response, err := p.Client.UpdateSubgraph(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "UpdateSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "UpdateSubgraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p PlatformClient) DeleteSubgraph(ctx context.Context, name, namespace string) *ApiError {
	request := connect.NewRequest(&platformv1.DeleteFederatedSubgraphRequest{
		SubgraphName: name,
		Namespace:    namespace,
	})
	response, err := p.Client.DeleteFederatedSubgraph(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "DeleteSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeleteSubgraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p PlatformClient) GetSubgraph(ctx context.Context, name, namespace string) (*platformv1.Subgraph, *ApiError) {
	request := connect.NewRequest(&platformv1.GetSubgraphByNameRequest{
		Name:      name,
		Namespace: namespace,
	})
	response, err := p.Client.GetSubgraphByName(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetSubgraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg.GetGraph(), nil
}

func (p PlatformClient) GetSubgraphs(ctx context.Context, namespace string) ([]*platformv1.Subgraph, *ApiError) {
	request := connect.NewRequest(&platformv1.GetSubgraphsRequest{
		Namespace: namespace,
	})

	response, err := p.Client.GetSubgraphs(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetSubgraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg.GetGraphs(), nil
}

func (p PlatformClient) GetSubgraphSchema(ctx context.Context, name, namespace string) (string, *ApiError) {
	request := connect.NewRequest(&platformv1.GetLatestSubgraphSDLRequest{
		Name:      name,
		Namespace: namespace,
	})

	response, err := p.Client.GetLatestSubgraphSDL(ctx, request)
	if err != nil {
		return "", &ApiError{Err: err, Reason: "GetSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return "", &ApiError{Err: ErrEmptyMsg, Reason: "GetSubgraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return "", apiError
	}

	return response.Msg.GetSdl(), nil
}

func (p PlatformClient) PublishSubgraph(ctx context.Context, name, namespace, schema string) (*platformv1.PublishFederatedSubgraphResponse, *ApiError) {
	request := connect.NewRequest(&platformv1.PublishFederatedSubgraphRequest{
		Name:      name,
		Namespace: namespace,
		Schema:    schema,
	})
	response, err := p.Client.PublishFederatedSubgraph(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "PublishSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "PublishSubgraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}
