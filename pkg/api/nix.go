package api

import (
	"context"

	"connectrpc.com/connect"
	nixv1alpha1 "go.jetify.com/pkg/api/gen/priv/nix/v1alpha1"
)

func (c *Client) GetAWSCredentials(ctx context.Context) (*nixv1alpha1.AWSCredentials, error) {
	r, err := c.nixClient().GetAWSCredentials(ctx, connect.NewRequest(&nixv1alpha1.GetAWSCredentialsRequest{}))
	if err != nil {
		return nil, err
	}
	return r.Msg.GetCredentials(), nil
}

func (c *Client) GetBinCache(
	ctx context.Context,
) (*nixv1alpha1.GetBinCacheResponse, error) {
	response, err := c.nixClient().GetBinCache(
		ctx,
		connect.NewRequest(&nixv1alpha1.GetBinCacheRequest{}),
	)
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}
