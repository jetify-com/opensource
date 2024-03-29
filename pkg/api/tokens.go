package api

import (
	"context"

	"connectrpc.com/connect"
	tokenservicev1alpha1 "go.jetpack.io/pkg/api/gen/priv/tokenservice/v1alpha1"
	"go.jetpack.io/pkg/ids"
)

func (c *Client) GetAccessToken(
	ctx context.Context,
	pat ids.PersonalAccessToken,
) (*tokenservicev1alpha1.GetAccessTokenResponse, error) {
	response, err := c.tokenClient().GetAccessToken(
		ctx,
		connect.NewRequest(&tokenservicev1alpha1.GetAccessTokenRequest{
			Token: pat.String(),
		}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}

func (c *Client) CreatePAT(
	ctx context.Context,
) (*tokenservicev1alpha1.CreatePATResponse, error) {
	response, err := c.tokenClient().CreatePAT(
		ctx,
		connect.NewRequest(&tokenservicev1alpha1.CreatePATRequest{}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}
