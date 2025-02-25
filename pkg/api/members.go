package api

import (
	"context"

	"connectrpc.com/connect"
	membersv1alpha1 "go.jetify.com/pkg/api/gen/priv/members/v1alpha1"
)

func (c *Client) GetMember(
	ctx context.Context,
	id string,
) (*membersv1alpha1.Member, error) {
	memberResponse, err := c.membersClient().GetMember(
		ctx,
		connect.NewRequest(&membersv1alpha1.GetMemberRequest{Id: id}),
	)
	if err != nil {
		return nil, err
	}
	return memberResponse.Msg.Member, nil
}
