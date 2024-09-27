package api

import (
	"context"

	"connectrpc.com/connect"

	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p PlatformClient) CreateNamespace(ctx context.Context, name string) *ApiError {
	request := connect.NewRequest(&platformv1.CreateNamespaceRequest{Name: name})
	response, err := p.Client.CreateNamespace(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "CreateSubgraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "CreateNamespace", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p PlatformClient) RenameNamespace(ctx context.Context, oldName, newName string) *ApiError {
	request := connect.NewRequest(&platformv1.RenameNamespaceRequest{
		Name:    oldName,
		NewName: newName,
	})
	response, err := p.Client.RenameNamespace(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "RenameNamespace", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "RenameNamespace", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p PlatformClient) DeleteNamespace(ctx context.Context, name string) error {
	request := connect.NewRequest(&platformv1.DeleteNamespaceRequest{Name: name})
	response, err := p.Client.DeleteNamespace(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "DeleteNamespace", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeleteNamespace", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p PlatformClient) ListNamespaces(ctx context.Context) ([]*platformv1.Namespace, *ApiError) {
	request := connect.NewRequest(&platformv1.GetNamespacesRequest{})
	response, err := p.Client.GetNamespaces(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "ListNamespaces", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "ListNamespaces", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg.Namespaces, nil
}
