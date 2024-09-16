package api

import (
	"context"

	"connectrpc.com/connect"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p PlatformClient) CreateSubgraph(ctx context.Context, name string, namespace string, routingUrl string, baseSubgraphName *string, labels []*platformv1.Label, subscriptionUrl *string, readme *string, isEventDrivenGraph *bool, isFeatureSubgraph *bool, subscriptionProtocol string, websocketSubprotocol string) error {
	request := connect.NewRequest(&platformv1.CreateFederatedSubgraphRequest{
		Name:                 name,
		Namespace:            namespace,
		RoutingUrl:           &routingUrl,
		Labels:               labels,
		SubscriptionUrl:      subscriptionUrl,
		Readme:               readme,
		WebsocketSubprotocol: resolveWebsocketSubprotocol(websocketSubprotocol),
		SubscriptionProtocol: resolveSubscriptionProtocol(subscriptionProtocol),
		IsEventDrivenGraph:   isEventDrivenGraph,
		BaseSubgraphName:     baseSubgraphName,
		IsFeatureSubgraph:    isFeatureSubgraph,
	})
	response, err := p.Client.CreateFederatedSubgraph(ctx, request)
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

func (p PlatformClient) UpdateSubgraph(ctx context.Context, name, namespace, routingUrl string, labels []*platformv1.Label, headers []string, subscriptionUrl, readme *string, unsetLabels *bool, websocketSubprotocol string, subscriptionProtocol string) error {
	request := connect.NewRequest(&platformv1.UpdateSubgraphRequest{
		Name:                 name,
		RoutingUrl:           &routingUrl,
		Labels:               labels,
		Headers:              headers,
		SubscriptionUrl:      subscriptionUrl,
		Readme:               readme,
		Namespace:            namespace,
		UnsetLabels:          unsetLabels,
		WebsocketSubprotocol: resolveWebsocketSubprotocol(websocketSubprotocol),
		SubscriptionProtocol: resolveSubscriptionProtocol(subscriptionProtocol),
	})

	response, err := p.Client.UpdateSubgraph(ctx, request)
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

func (p PlatformClient) DeleteSubgraph(ctx context.Context, name, namespace string) error {
	request := connect.NewRequest(&platformv1.DeleteFederatedSubgraphRequest{
		SubgraphName: name,
		Namespace:    namespace,
	})
	response, err := p.Client.DeleteFederatedSubgraph(ctx, request)
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

func (p PlatformClient) GetSubgraph(ctx context.Context, name, namespace string) (*platformv1.Subgraph, error) {
	request := connect.NewRequest(&platformv1.GetSubgraphByNameRequest{
		Name:      name,
		Namespace: namespace,
	})
	response, err := p.Client.GetSubgraphByName(ctx, request)
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

	return response.Msg.GetGraph(), nil
}

func (p PlatformClient) PublishSubgraph(ctx context.Context, name, namespace, schema string) (*platformv1.PublishFederatedSubgraphResponse, error) {
	request := connect.NewRequest(&platformv1.PublishFederatedSubgraphRequest{
		Name:      name,
		Namespace: namespace,
		Schema:    schema,
	})
	response, err := p.Client.PublishFederatedSubgraph(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg == nil {
		return nil, ErrEmptyMsg
	}

	return response.Msg, nil
}
