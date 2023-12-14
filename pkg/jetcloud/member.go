package jetcloud

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	membersv1alpha1 "go.jetpack.io/pkg/api/gen/priv/members/v1alpha1"
	"go.jetpack.io/pkg/api/gen/priv/members/v1alpha1/membersv1alpha1connect"
	"go.jetpack.io/pkg/auth/session"
)

func (c *Client) GetMember(
	ctx context.Context,
	tok *session.Token,
	id string,
) (*membersv1alpha1.Member, error) {
	memberResponse, err := membersv1alpha1connect.NewMembersServiceClient(
		http.DefaultClient,
		c.APIHost,
	).GetMember(ctx, newRequest(&membersv1alpha1.GetMemberRequest{
		Id: id,
	}, tok.AccessToken))
	if err != nil {
		return nil, err
	}
	return memberResponse.Msg.Member, nil
}

func newRequest[T any](message *T, token string) *connect.Request[T] {
	req := connect.NewRequest(message)
	req.Header().Set("Authorization", "Bearer "+token)
	return req
}
