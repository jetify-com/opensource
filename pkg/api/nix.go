package api

import (
	"context"

	"connectrpc.com/connect"
	nixv1alpha1 "go.jetpack.io/pkg/api/gen/priv/nix/v1alpha1"
)

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
