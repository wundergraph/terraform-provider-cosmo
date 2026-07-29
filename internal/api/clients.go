package api

import (
	"context"

	"connectrpc.com/connect"

	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p *PlatformClient) GetClients(ctx context.Context, graphName, namespace string) ([]*platformv1.ClientInfo, *ApiError) {
	request := connect.NewRequest(&platformv1.GetClientsRequest{
		FedGraphName: graphName,
		Namespace:    namespace,
	})

	response, err := p.Client.GetClients(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetClients", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetClients", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg.Clients, nil
}

func (p *PlatformClient) GetClient(ctx context.Context, graphName, namespace, clientName string) (*platformv1.ClientInfo, *ApiError) {
	clients, apiError := p.GetClients(ctx, graphName, namespace)
	if apiError != nil {
		return nil, apiError
	}

	for _, client := range clients {
		if client.Name == clientName {
			return client, nil
		}
	}

	return nil, &ApiError{Err: ErrNotFound, Reason: "GetClient", Status: common.EnumStatusCode_ERR}
}
