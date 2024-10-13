package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

func (p *PlatformClient) CreateFeatureFlag(ctx context.Context, name, namespace string, enabled bool, labels []*platformv1.Label, featureSubgraphNames []string) (*platformv1.CreateFeatureFlagResponse, *ApiError) {
	request := connect.NewRequest(&platformv1.CreateFeatureFlagRequest{
		Name:                 name,
		Namespace:            namespace,
		IsEnabled:            enabled,
		Labels:               labels,
		FeatureSubgraphNames: featureSubgraphNames,
	})

	response, err := p.Client.CreateFeatureFlag(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "CreateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "CreateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) UpdateFeatureFlag(ctx context.Context, name, namespace string, unsetLabels bool, labels []*platformv1.Label, featureSubgraphNames []string) (*platformv1.UpdateFeatureFlagResponse, *ApiError) {
	request := connect.NewRequest(&platformv1.UpdateFeatureFlagRequest{
		Name:                 name,
		Namespace:            namespace,
		Labels:               labels,
		FeatureSubgraphNames: featureSubgraphNames,
		UnsetLabels:          unsetLabels,
	})

	response, err := p.Client.UpdateFeatureFlag(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "UpdateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "UpdateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}

func (p *PlatformClient) DeleteFeatureFlag(ctx context.Context, name, namespace string) *ApiError {
	request := connect.NewRequest(&platformv1.DeleteFeatureFlagRequest{
		Name:      name,
		Namespace: namespace,
	})

	response, err := p.Client.DeleteFeatureFlag(ctx, request)
	if err != nil {
		return &ApiError{Err: err, Reason: "DeleteFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeleteFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return apiError
	}

	return nil
}

func (p *PlatformClient) GetFeatureFlag(ctx context.Context, name, namespace string) (*platformv1.GetFeatureFlagByNameResponse, *ApiError) {
	request := connect.NewRequest(&platformv1.GetFeatureFlagByNameRequest{
		Name:      name,
		Namespace: namespace,
	})

	response, err := p.Client.GetFeatureFlagByName(ctx, request)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if response.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(response.Msg.GetResponse().Code, response.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	return response.Msg, nil
}
