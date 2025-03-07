package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p *PlatformClient) CreateMonograph(ctx context.Context, name string, namespace string, routingUrl string, graphUrl string, subscriptionUrl *string, readme *string, websocketSubprotocol string, subscriptionProtocol string, admissionWebhookUrl string, admissionWebhookSecret string) (*platformv1.CreateMonographResponse, *ApiError) {
	request := connect.NewRequest(&platformv1.CreateMonographRequest{
		Name:                   name,
		Namespace:              namespace,
		RoutingUrl:             routingUrl,
		GraphUrl:               graphUrl,
		SubscriptionUrl:        subscriptionUrl,
		Readme:                 readme,
		WebsocketSubprotocol:   ResolveWebsocketSubprotocol(websocketSubprotocol),
		SubscriptionProtocol:   ResolveSubscriptionProtocol(subscriptionProtocol),
		AdmissionWebhookURL:    admissionWebhookUrl,
		AdmissionWebhookSecret: &admissionWebhookSecret,
	})
	response, err := p.Client.CreateMonograph(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "CreateMonograph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "CreateMonograph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) UpdateMonograph(ctx context.Context, name string, namespace string, routingUrl string, graphUrl string, subscriptionUrl *string, readme *string, websocketSubprotocol string, subscriptionProtocol string, admissionWebhookUrl string, admissionWebhookSecret string) *ApiError {
	request := connect.NewRequest(&platformv1.UpdateMonographRequest{
		Name:                   name,
		Namespace:              namespace,
		RoutingUrl:             routingUrl,
		GraphUrl:               graphUrl,
		SubscriptionUrl:        subscriptionUrl,
		Readme:                 readme,
		WebsocketSubprotocol:   ResolveWebsocketSubprotocol(websocketSubprotocol),
		SubscriptionProtocol:   ResolveSubscriptionProtocol(subscriptionProtocol),
		AdmissionWebhookURL:    &admissionWebhookUrl,
		AdmissionWebhookSecret: &admissionWebhookSecret,
	})
	response, err := p.Client.UpdateMonograph(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "UpdateMonograph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "UpdateMonograph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p *PlatformClient) DeleteMonograph(ctx context.Context, name string, namespace string) *ApiError {
	request := connect.NewRequest(&platformv1.DeleteMonographRequest{
		Name:      name,
		Namespace: namespace,
	})
	response, err := p.Client.DeleteMonograph(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "DeleteMonograph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeleteMonograph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p *PlatformClient) GetMonograph(ctx context.Context, name string, namespace string) (*platformv1.FederatedGraph, *ApiError) {
	request := connect.NewRequest(&platformv1.GetFederatedGraphByNameRequest{
		Name:      name,
		Namespace: namespace,
	})
	response, err := p.Client.GetFederatedGraphByName(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetMonograph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetMonograph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg.Graph, nil
}

func (p *PlatformClient) GetMonographByID(ctx context.Context, id string) (*platformv1.FederatedGraph, *ApiError) {
	request := connect.NewRequest(&platformv1.GetFederatedGraphByIdRequest{
		Id: id,
	})

	response, err := p.Client.GetFederatedGraphById(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetMonographByID", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetMonographByID", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg.Graph, nil
}

func (p *PlatformClient) PublishMonograph(ctx context.Context, name string, namespace string, schema string) *ApiError {
	request := connect.NewRequest(&platformv1.PublishMonographRequest{
		Name:      name,
		Namespace: namespace,
		Schema:    schema,
	})
	response, err := p.Client.PublishMonograph(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "PublishMonograph", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "PublishMonograph", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}
