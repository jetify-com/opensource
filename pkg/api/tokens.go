package api

import (
	"context"

	"connectrpc.com/connect"
	tokenservicev1alpha1 "go.jetpack.io/pkg/api/gen/priv/tokenservice/v1alpha1"
	"go.jetpack.io/pkg/ids"
)

func (c *Client) GetAccessToken(
	ctx context.Context,
	apiToken ids.APIToken,
) (*tokenservicev1alpha1.GetAccessTokenResponse, error) {
	response, err := c.tokenClient().GetAccessToken(
		ctx,
		connect.NewRequest(&tokenservicev1alpha1.GetAccessTokenRequest{
			Token: apiToken.String(),
		}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}

func (c *Client) CreateToken(
	ctx context.Context,
) (*tokenservicev1alpha1.CreateTokenResponse, error) {
	response, err := c.tokenClient().CreateToken(
		ctx,
		connect.NewRequest(&tokenservicev1alpha1.CreateTokenRequest{}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}
