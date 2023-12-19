package api

import (
	"context"

	"connectrpc.com/connect"
	membersv1alpha1 "go.jetpack.io/pkg/api/gen/priv/members/v1alpha1"
	"go.jetpack.io/pkg/api/gen/priv/members/v1alpha1/membersv1alpha1connect"
)

func (c *Client) GetMember(
	ctx context.Context,
	id string,
) (*membersv1alpha1.Member, error) {
	memberResponse, err := membersv1alpha1connect.NewMembersServiceClient(
		c.httpClient,
		c.Host,
	).GetMember(ctx, connect.NewRequest(&membersv1alpha1.GetMemberRequest{
		Id: id,
	}))
	if err != nil {
		return nil, err
	}
	return memberResponse.Msg.Member, nil
}
