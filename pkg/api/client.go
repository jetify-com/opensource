package api

import (
	"context"
	"sync"

	"go.jetpack.io/pkg/api/gen/priv/members/v1alpha1/membersv1alpha1connect"
	"go.jetpack.io/pkg/api/gen/priv/projects/v1alpha1/projectsv1alpha1connect"
	"go.jetpack.io/pkg/auth/session"
	"golang.org/x/oauth2"
)

// Client manages state for interacting with the JetCloud API, as well as
// communicating with the JetCloud API.
type Client struct {
	membersClient  func() membersv1alpha1connect.MembersServiceClient
	projectsClient func() projectsv1alpha1connect.ProjectsServiceClient
}

func NewClient(ctx context.Context, host string, token *session.Token) *Client {
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken},
	))
	return &Client{
		membersClient: sync.OnceValue(func() membersv1alpha1connect.MembersServiceClient {
			return membersv1alpha1connect.NewMembersServiceClient(
				httpClient,
				host,
			)
		}),
		projectsClient: sync.OnceValue(func() projectsv1alpha1connect.ProjectsServiceClient {
			return projectsv1alpha1connect.NewProjectsServiceClient(
				httpClient,
				host,
			)
		}),
	}
}
