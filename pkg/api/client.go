package api

import (
	"context"
	"sync"

	"go.jetify.com/pkg/api/gen/priv/members/v1alpha1/membersv1alpha1connect"
	"go.jetify.com/pkg/api/gen/priv/nix/v1alpha1/nixv1alpha1connect"
	"go.jetify.com/pkg/api/gen/priv/projects/v1alpha1/projectsv1alpha1connect"
	"go.jetify.com/pkg/api/gen/priv/secrets/v1alpha1/secretsv1alpha1connect"
	"go.jetify.com/pkg/api/gen/priv/tokenservice/v1alpha1/tokenservicev1alpha1connect"
	"go.jetify.com/pkg/auth/session"
	"golang.org/x/oauth2"
)

// Client manages state for interacting with the JetCloud API, as well as
// communicating with the JetCloud API.
type Client struct {
	membersClient  func() membersv1alpha1connect.MembersServiceClient
	nixClient      func() nixv1alpha1connect.NixServiceClient
	ProjectsClient func() projectsv1alpha1connect.ProjectsServiceClient
	secretsClient  func() secretsv1alpha1connect.SecretsServiceClient
	tokenClient    func() tokenservicev1alpha1connect.TokenServiceClient
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
		nixClient: sync.OnceValue(func() nixv1alpha1connect.NixServiceClient {
			return nixv1alpha1connect.NewNixServiceClient(
				httpClient,
				host,
			)
		}),
		ProjectsClient: sync.OnceValue(func() projectsv1alpha1connect.ProjectsServiceClient {
			return projectsv1alpha1connect.NewProjectsServiceClient(
				httpClient,
				host,
			)
		}),
		secretsClient: sync.OnceValue(func() secretsv1alpha1connect.SecretsServiceClient {
			return secretsv1alpha1connect.NewSecretsServiceClient(
				httpClient,
				host,
			)
		}),
		tokenClient: sync.OnceValue(func() tokenservicev1alpha1connect.TokenServiceClient {
			return tokenservicev1alpha1connect.NewTokenServiceClient(
				httpClient,
				host,
			)
		}),
	}
}
