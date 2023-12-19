package api

import (
	"context"

	"connectrpc.com/connect"
	projectsv1alpha1 "go.jetpack.io/pkg/api/gen/priv/projects/v1alpha1"
	"go.jetpack.io/pkg/api/gen/priv/projects/v1alpha1/projectsv1alpha1connect"
	"go.jetpack.io/pkg/id"
)

func (c *Client) ListProjects(
	ctx context.Context,
	orgID id.OrgID,
) ([]*projectsv1alpha1.Project, error) {
	memberResponse, err := projectsv1alpha1connect.NewProjectsServiceClient(
		c.httpClient,
		c.Host,
	).ListProjects(ctx, connect.NewRequest(
		&projectsv1alpha1.ListProjectsRequest{
			OrgId: orgID.String(),
		},
	))
	if err != nil {
		return nil, err
	}
	return memberResponse.Msg.Projects, nil
}

func (c *Client) CreateProject(
	ctx context.Context,
	orgID id.OrgID,
	repoURL string,
	directory string,
	name string,
) (*projectsv1alpha1.Project, error) {
	memberResponse, err := projectsv1alpha1connect.NewProjectsServiceClient(
		c.httpClient,
		c.Host,
	).CreateProject(ctx, connect.NewRequest(
		&projectsv1alpha1.CreateProjectRequest{
			OrgId: orgID.String(),
			Project: &projectsv1alpha1.Project{
				Repo:      repoURL,
				Directory: directory,
				Name:      name,
			},
		},
	))
	if err != nil {
		return nil, err
	}
	return memberResponse.Msg.Project, nil
}
