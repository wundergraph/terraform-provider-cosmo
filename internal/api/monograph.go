package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/wg/cosmo/platform/v1"
)

func (p PlatformClient) CreateMonograph(ctx context.Context, name string, namespace string, routingUrl string, graphUrl string, subscriptionUrl *string, readme *string, websocketSubprotocol string, subscriptionProtocol string, admissionWebhookUrl string, admissionWebhookSecret string) error {
	request := connect.NewRequest(&platformv1.CreateMonographRequest{
		Name:                   name,
		Namespace:              namespace,
		RoutingUrl:             routingUrl,
		GraphUrl:               graphUrl,
		SubscriptionUrl:        subscriptionUrl,
		Readme:                 readme,
		WebsocketSubprotocol:   resolveWebsocketSubprotocol(websocketSubprotocol),
		SubscriptionProtocol:   resolveSubscriptionProtocol(subscriptionProtocol),
		AdmissionWebhookURL:    admissionWebhookUrl,
		AdmissionWebhookSecret: &admissionWebhookSecret,
	})
	response, err := p.Client.CreateMonograph(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to create monograph: %s", response.Msg)
	}

	return nil
}

func (p PlatformClient) UpdateMonograph(ctx context.Context, name string, namespace string, routingUrl string, graphUrl string, subscriptionUrl *string, readme *string, websocketSubprotocol string, subscriptionProtocol string, admissionWebhookUrl string, admissionWebhookSecret string) error {
	request := connect.NewRequest(&platformv1.UpdateMonographRequest{
		Name:                   name,
		Namespace:              namespace,
		RoutingUrl:             routingUrl,
		GraphUrl:               graphUrl,
		SubscriptionUrl:        subscriptionUrl,
		Readme:                 readme,
		WebsocketSubprotocol:   resolveWebsocketSubprotocol(websocketSubprotocol),
		SubscriptionProtocol:   resolveSubscriptionProtocol(subscriptionProtocol),
		AdmissionWebhookURL:    &admissionWebhookUrl,
		AdmissionWebhookSecret: &admissionWebhookSecret,
	})
	response, err := p.Client.UpdateMonograph(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to update monograph: %s", response.Msg)
	}

	return nil
}

func (p PlatformClient) DeleteMonograph(ctx context.Context, name string, namespace string) error {
	request := connect.NewRequest(&platformv1.DeleteMonographRequest{
		Name:      name,
		Namespace: namespace,
	})
	response, err := p.Client.DeleteMonograph(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to delete monograph: %s", response.Msg)
	}

	return nil
}

func (p PlatformClient) GetMonograph(ctx context.Context, name string, namespace string) (*platformv1.FederatedGraph, error) {
	request := connect.NewRequest(&platformv1.GetFederatedGraphByNameRequest{
		Name:      name,
		Namespace: namespace,
	})
	response, err := p.Client.GetFederatedGraphByName(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return nil, fmt.Errorf("failed to get monograph: %s", response.Msg)
	}

	return response.Msg.Graph, nil
}
