package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/wundergraph/cosmo/connect-go/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/wg/cosmo/platform/v1"
	"github.com/wundergraph/cosmo/connect-go/wg/cosmo/platform/v1/platformv1connect"
)

func CreateNamespace(ctx context.Context, client platformv1connect.PlatformServiceClient, apiKey, name string) error {
	request := connect.NewRequest(&platformv1.CreateNamespaceRequest{Name: name})
	request.Header().Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	response, err := client.CreateNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to create namespace: %s", response.Msg)
	}

	return err
}

func RenameNamespace(ctx context.Context, client platformv1connect.PlatformServiceClient, apiKey, oldName, newName string) error {
	request := connect.NewRequest(&platformv1.RenameNamespaceRequest{
		Name:    oldName,
		NewName: newName,
	})
	request.Header().Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	response, err := client.RenameNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to rename the namespace: %s", response.Msg)
	}
	return err
}

func DeleteNamespace(ctx context.Context, client platformv1connect.PlatformServiceClient, apiKey, name string) error {
	request := connect.NewRequest(&platformv1.DeleteNamespaceRequest{Name: name})
	request.Header().Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	response, err := client.DeleteNamespace(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to delete namespace: %s", response.Msg)
	}
	return err
}

func ListNamespaces(ctx context.Context, client platformv1connect.PlatformServiceClient, apiKey string) ([]*platformv1.Namespace, error) {
	request := connect.NewRequest(&platformv1.GetNamespacesRequest{})
	request.Header().Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	response, err := client.GetNamespaces(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return nil, fmt.Errorf("failed to delete namespace: %s", response.Msg)
	}

	return response.Msg.Namespaces, nil
}
