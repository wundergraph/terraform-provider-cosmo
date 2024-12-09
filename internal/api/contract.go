package api

import (
	"context"

	"connectrpc.com/connect"
	common "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p *PlatformClient) CreateContract(ctx context.Context, data *platformv1.CreateContractRequest) (*platformv1.CreateContractResponse, *ApiError) {
	request := connect.NewRequest(data)

	response, err := p.Client.CreateContract(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "CreateContract", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "CreateContract", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) UpdateContract(ctx context.Context, data *platformv1.UpdateContractRequest) (*platformv1.UpdateContractResponse, *ApiError) {
	request := connect.NewRequest(data)

	response, err := p.Client.UpdateContract(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "UpdateContract", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "UpdateContract", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) DeleteContract(ctx context.Context, name, namespace string, supportsFederation bool) *ApiError {
	if supportsFederation {
		return p.DeleteFederatedGraph(ctx, name, namespace)
	} else {
		return p.DeleteMonograph(ctx, name, namespace)
	}
}

func (p *PlatformClient) GetContract(ctx context.Context, name, namespace string) (*platformv1.GetFederatedGraphByNameResponse, *ApiError) {
	return p.GetFederatedGraph(ctx, name, namespace)
}
