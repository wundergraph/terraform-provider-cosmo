package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p PlatformClient) CreateNamespace(ctx context.Context, name string) error {
	request := connect.NewRequest(&platformv1.CreateNamespaceRequest{Name: name})
	response, err := p.Client.CreateNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg == nil {
		return fmt.Errorf("failed to create namespace: %s, the server response is nil", name)
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to create namespace: %s", response.Msg)
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
		return fmt.Errorf("failed to rename the namespace: %s, the server response is nil", oldName)
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to rename the namespace: %s", response.Msg)
	}
	return err
}

func (p PlatformClient) DeleteNamespace(ctx context.Context, name string) error {
	request := connect.NewRequest(&platformv1.DeleteNamespaceRequest{Name: name})
	response, err := p.Client.DeleteNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg == nil {
		return fmt.Errorf("failed to delete namespace: %s, the server response is nil", name)
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to delete namespace: %s", response.Msg)
	}
	return err
}

func (p PlatformClient) ListNamespaces(ctx context.Context) ([]*platformv1.Namespace, error) {
	request := connect.NewRequest(&platformv1.GetNamespacesRequest{})
	response, err := p.Client.GetNamespaces(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg == nil {
		return nil, fmt.Errorf("failed to list namespaces, the server response is nil")
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return nil, fmt.Errorf("failed to list namespaces: %s", response.Msg)
	}

	return response.Msg.Namespaces, nil
}
