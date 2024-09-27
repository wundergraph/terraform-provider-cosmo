package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p PlatformClient) CreateToken(ctx context.Context, name, graphName, namespace string) (string, *ApiError) {
	request := connect.NewRequest(&platformv1.CreateFederatedGraphTokenRequest{
		GraphName: graphName,
		Namespace: namespace,
		TokenName: name,
	})

	response, err := p.Client.CreateFederatedGraphToken(ctx, request)
	if err != nil {
		return "", &ApiError{Err: err, Reason: "CreateToken", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return "", &ApiError{Err: ErrEmptyMsg, Reason: "CreateToken", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return "", &ApiError{Err: fmt.Errorf("failed to create token: %s", response.Msg.GetResponse().GetDetails()), Reason: "CreateToken", Status: common.EnumStatusCode_ERR}
	}

	return response.Msg.Token, nil
}

func (p PlatformClient) DeleteToken(ctx context.Context, tokenName, graphName, namespace string) *ApiError {
	request := connect.NewRequest(&platformv1.DeleteRouterTokenRequest{
		TokenName: tokenName,
		Namespace: namespace,
	})

	response, err := p.Client.DeleteRouterToken(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "DeleteToken", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeleteToken", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return &ApiError{Err: fmt.Errorf("failed to delete token: %s", response.Msg), Reason: "DeleteToken", Status: common.EnumStatusCode_ERR}
	}

	return nil
}
