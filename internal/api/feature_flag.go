package api

import (
	"context"

	"connectrpc.com/connect"

	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
	platformv1 "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/platform/v1"
)

type FeatureFlag struct {
	*platformv1.FeatureFlag
	FeatureSubgraphNames []string
}

func (p *PlatformClient) CreateFeatureFlag(ctx context.Context, data *FeatureFlag) *ApiError {

	req := connect.NewRequest(&platformv1.CreateFeatureFlagRequest{
		Name:                 data.Name,
		Namespace:            data.Namespace,
		Labels:               data.Labels,
		FeatureSubgraphNames: data.FeatureSubgraphNames,
		IsEnabled:            data.IsEnabled,
	})

	resp, err := p.Client.CreateFeatureFlag(ctx, req)
	if err != nil {
		return &ApiError{Err: err, Reason: "CreateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if resp.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "CreateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	return handleErrorCodes(resp.Msg.GetResponse().Code, resp.Msg.String())
}

func (p *PlatformClient) GetFeatureFlag(ctx context.Context, name, namespace string) (*FeatureFlag, *ApiError) {
	req := connect.NewRequest(&platformv1.GetFeatureFlagByNameRequest{
		Name:      name,
		Namespace: namespace,
	})

	resp, err := p.Client.GetFeatureFlagByName(ctx, req)
	if err != nil {
		return nil, &ApiError{Err: err, Reason: "GetFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if resp.Msg == nil {
		return nil, &ApiError{Err: ErrEmptyMsg, Reason: "GetFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	apiError := handleErrorCodes(resp.Msg.GetResponse().Code, resp.Msg.String())
	if apiError != nil {
		return nil, apiError
	}

	var featureSubgraphNames []string
	for _, sg := range resp.Msg.GetFeatureSubgraphs() {
		featureSubgraphNames = append(featureSubgraphNames, sg.GetName())
	}

	return &FeatureFlag{
		FeatureFlag:          resp.Msg.GetFeatureFlag(),
		FeatureSubgraphNames: featureSubgraphNames,
	}, nil
}

func (p *PlatformClient) UpdateFeatureFlag(ctx context.Context, data *FeatureFlag) *ApiError {
	req := connect.NewRequest(&platformv1.UpdateFeatureFlagRequest{
		Name:                 data.Name,
		Namespace:            data.Namespace,
		Labels:               data.Labels,
		FeatureSubgraphNames: data.FeatureSubgraphNames,
		UnsetLabels:          len(data.GetLabels()) == 0,
	})

	resp, apiErr := p.Client.UpdateFeatureFlag(ctx, req)
	if apiErr != nil {
		return &ApiError{Err: apiErr, Reason: "UpdateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if resp.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "UpdateFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	return handleErrorCodes(resp.Msg.GetResponse().Code, resp.Msg.String())

}

func (p *PlatformClient) SetFeatureFlagState(ctx context.Context, name, namespace string, enabled bool) *ApiError {

	req := connect.NewRequest(&platformv1.EnableFeatureFlagRequest{
		Name:      name,
		Namespace: namespace,
		Enabled:   enabled,
	})

	resp, apiErr := p.Client.EnableFeatureFlag(ctx, req)
	if apiErr != nil {
		return &ApiError{Err: apiErr, Reason: "EnableFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if resp.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "EnableFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	return handleErrorCodes(resp.Msg.GetResponse().Code, resp.Msg.String())
}

func (p *PlatformClient) DeleteFeatureFlag(ctx context.Context, name, namespace string) *ApiError {
	resp, apiErr := p.Client.DeleteFeatureFlag(ctx, connect.NewRequest(&platformv1.DeleteFeatureFlagRequest{
		Name:      name,
		Namespace: namespace,
	}))

	if apiErr != nil {
		return &ApiError{Err: apiErr, Reason: "DeleteFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	if resp.Msg == nil {
		return &ApiError{Err: ErrEmptyMsg, Reason: "DeleteFeatureFlag", Status: common.EnumStatusCode_ERR}
	}

	return handleErrorCodes(resp.Msg.GetResponse().Code, resp.Msg.String())
}
