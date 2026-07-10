package api

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"

	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p *PlatformClient) PublishPersistedOperations(ctx context.Context, graphName, namespace, clientName string, operations []*platformv1.PersistedOperation) *ApiError {
	request := connect.NewRequest(&platformv1.PublishPersistedOperationsRequest{
		FedGraphName: graphName,
		Namespace:    namespace,
		ClientName:   clientName,
		Operations:   operations,
	})

	response, err := p.Client.PublishPersistedOperations(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "PublishPersistedOperations", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "PublishPersistedOperations", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	var conflicts []string
	for _, operation := range response.Msg.Operations {
		if operation.Status == platformv1.PublishedOperationStatus_CONFLICT {
			conflicts = append(conflicts, operation.Id)
		}
	}
	if len(conflicts) > 0 {
		return &ApiError{
			Err:    fmt.Errorf("operations already exist with different contents for the same id: %s", strings.Join(conflicts, ", ")),
			Reason: "PublishPersistedOperations",
			Status: common.EnumStatusCode_ERR,
		}
	}

	return nil
}

func (p *PlatformClient) GetPersistedOperations(ctx context.Context, graphName, namespace, clientId string) ([]*platformv1.GetPersistedOperationsResponse_Operation, *ApiError) {
	request := connect.NewRequest(&platformv1.GetPersistedOperationsRequest{
		FederatedGraphName: graphName,
		Namespace:          namespace,
		ClientId:           clientId,
	})

	response, err := p.Client.GetPersistedOperations(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetPersistedOperations", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetPersistedOperations", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg.Operations, nil
}

func (p *PlatformClient) DeletePersistedOperation(ctx context.Context, graphName, namespace, clientName, operationId string) *ApiError {
	request := connect.NewRequest(&platformv1.DeletePersistedOperationRequest{
		FedGraphName: graphName,
		Namespace:    namespace,
		ClientName:   clientName,
		OperationId:  operationId,
	})

	response, err := p.Client.DeletePersistedOperation(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "DeletePersistedOperation", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeletePersistedOperation", Status: common.EnumStatusCode_ERR}
	}

	return handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
}
