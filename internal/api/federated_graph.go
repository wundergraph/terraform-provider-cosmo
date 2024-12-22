package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p *PlatformClient) CreateFederatedGraph(ctx context.Context, admissionWebhookSecret *string, graph *platformv1.FederatedGraph) (*platformv1.CreateFederatedGraphResponse, *ApiError) {
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
		return nil, &ApiError{Err: err, Reason: "CreateFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "CreateFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) UpdateFederatedGraph(ctx context.Context, admissionWebhookSecret *string, graph *platformv1.FederatedGraph) (*platformv1.UpdateFederatedGraphResponse, *ApiError) {
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
		Readme:                 graph.Readme,
	})

	response, err := p.Client.UpdateFederatedGraph(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "UpdateFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "UpdateFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) DeleteFederatedGraph(ctx context.Context, name, namespace string) *ApiError {
	request := connect.NewRequest(&platformv1.DeleteFederatedGraphRequest{
		Name:      name,
		Namespace: namespace,
	})

	response, err := p.Client.DeleteFederatedGraph(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "DeleteFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeleteFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p *PlatformClient) GetFederatedGraph(ctx context.Context, name, namespace string) (*platformv1.GetFederatedGraphByNameResponse, *ApiError) {
	request := connect.NewRequest(&platformv1.GetFederatedGraphByNameRequest{
		Name:      name,
		Namespace: namespace,
	})

	response, err := p.Client.GetFederatedGraphByName(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) GetFederatedGraphById(ctx context.Context, id string) (*platformv1.GetFederatedGraphByIdResponse, *ApiError) {
	request := connect.NewRequest(&platformv1.GetFederatedGraphByIdRequest{
		Id: id,
	})

	response, err := p.Client.GetFederatedGraphById(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetFederatedGraph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}
