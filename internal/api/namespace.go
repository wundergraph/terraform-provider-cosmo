package api

import (
	"context"

	"connectrpc.com/connect"

	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p PlatformClient) CreateNamespace(ctx context.Context, name string) error {
	request := connect.NewRequest(&platformv1.CreateNamespaceRequest{Name: name})
	response, err := p.Client.CreateNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg == nil {
		return ErrEmptyMsg
	}

	err = handleErrorCodes(response.Msg.GetResponse().Code)
	if err != nil {
		return err
	}

	return err
}

func (p PlatformClient) RenameNamespace(ctx context.Context, oldName, newName string) error {
	request := connect.NewRequest(&platformv1.RenameNamespaceRequest{
		Name:    oldName,
		NewName: newName,
	})
	response, err := p.Client.RenameNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg == nil {
		return ErrEmptyMsg
	}

	err = handleErrorCodes(response.Msg.GetResponse().Code)
	if err != nil {
		return err
	}

	return nil
}

func (p PlatformClient) DeleteNamespace(ctx context.Context, name string) error {
	request := connect.NewRequest(&platformv1.DeleteNamespaceRequest{Name: name})
	response, err := p.Client.DeleteNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg == nil {
		return ErrEmptyMsg
	}

	err = handleErrorCodes(response.Msg.GetResponse().Code)
	if err != nil {
		return err
	}

	return nil
}

func (p PlatformClient) ListNamespaces(ctx context.Context) ([]*platformv1.Namespace, error) {
	request := connect.NewRequest(&platformv1.GetNamespacesRequest{})
	response, err := p.Client.GetNamespaces(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg == nil {
		return nil, ErrEmptyMsg
	}

	err = handleErrorCodes(response.Msg.GetResponse().Code)
	if err != nil {
		return nil, err
	}

	return response.Msg.Namespaces, nil
}
