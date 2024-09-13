package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/wg/cosmo/platform/v1"
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

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to create subgraph: %s", response.Msg)
	}

	return nil
}

func (p PlatformClient) UpdateSubgraph(ctx context.Context,  name, namespace, routingUrl string, labels []*platformv1.Label, headers []string, subscriptionUrl, readme *string, unsetLabels *bool, websocketSubprotocol string, subscriptionProtocol string) error {
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

	_, err := p.Client.UpdateSubgraph(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (p PlatformClient) DeleteSubgraph(ctx context.Context,  name, namespace string) error {
	request := connect.NewRequest(&platformv1.DeleteFederatedSubgraphRequest{
		SubgraphName: name,
		Namespace:    namespace,
	})
	response, err := p.Client.DeleteFederatedSubgraph(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to delete subgraph: %s", response.Msg)
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

	subgraph := &platformv1.Subgraph{
		Id:         response.Msg.Graph.Id,
		Name:       response.Msg.Graph.Name,
		Namespace:  response.Msg.Graph.Namespace,
		RoutingURL: response.Msg.Graph.RoutingURL,
	}

	return subgraph, nil
}
