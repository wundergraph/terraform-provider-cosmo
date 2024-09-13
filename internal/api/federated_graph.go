package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/wg/cosmo/platform/v1"
)

func (p *PlatformClient) CreateFederatedGraph(ctx context.Context, admissionWebhookSecret *string, graph *platformv1.FederatedGraph) (*platformv1.CreateFederatedGraphResponse, error) {
	var admissionWebhookURL string
	if graph.AdmissionWebhookUrl != nil {
		admissionWebhookURL = *graph.AdmissionWebhookUrl
	} else {
		admissionWebhookURL = ""
	}

	request := connect.NewRequest(&platformv1.CreateFederatedGraphRequest{
		Name:                   graph.Name,
		Namespace:              graph.Namespace,
		RoutingUrl:             graph.RoutingURL,
		AdmissionWebhookURL:    admissionWebhookURL,
		AdmissionWebhookSecret: admissionWebhookSecret,
		Readme:                 graph.Readme,
		LabelMatchers:          graph.LabelMatchers,
	})

	response, err := p.Client.CreateFederatedGraph(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return nil, fmt.Errorf("failed to create federated graph: %s", response.Msg.GetResponse().GetDetails())
	}

	return response.Msg, nil
}

func (p *PlatformClient) UpdateFederatedGraph(ctx context.Context, admissionWebhookSecret *string, graph *platformv1.FederatedGraph) (*platformv1.UpdateFederatedGraphResponse, error) {
	var admissionWebhookURL *string
	if graph.AdmissionWebhookUrl != nil {
		admissionWebhookURL = graph.AdmissionWebhookUrl
	}

	request := connect.NewRequest(&platformv1.UpdateFederatedGraphRequest{
		Name:                   graph.Name,
		Namespace:              graph.Namespace,
		RoutingUrl:             graph.RoutingURL,
		AdmissionWebhookURL:    admissionWebhookURL,
		AdmissionWebhookSecret: admissionWebhookSecret,
		LabelMatchers:          graph.LabelMatchers,
	})

	response, err := p.Client.UpdateFederatedGraph(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return nil, fmt.Errorf("failed to update federated graph: %s", response.Msg)
	}

	return response.Msg, nil
}

func (p *PlatformClient) DeleteFederatedGraph(ctx context.Context, name, namespace string) error {
	request := connect.NewRequest(&platformv1.DeleteFederatedGraphRequest{
		Name:      name,
		Namespace: namespace,
	})

	response, err := p.Client.DeleteFederatedGraph(ctx, request)
	if err != nil {
		return err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return fmt.Errorf("failed to delete federated graph: %s", response.Msg)
	}

	return nil
}

func (p *PlatformClient) GetFederatedGraph(ctx context.Context, name, namespace string) (*platformv1.GetFederatedGraphByNameResponse, error) {
	request := connect.NewRequest(&platformv1.GetFederatedGraphByNameRequest{
		Name:      name,
		Namespace: namespace,
	})

	response, err := p.Client.GetFederatedGraphByName(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.Msg.GetResponse().Code != common.EnumStatusCode_OK {
		return nil, fmt.Errorf("failed to get federated graph: %s", response.Msg)
	}

	return response.Msg, nil
}
