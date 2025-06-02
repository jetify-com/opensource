package api

import (
	"context"

	"connectrpc.com/connect"
	projectsv1alpha1 "go.jetify.com/pkg/api/gen/priv/projects/v1alpha1"
	"go.jetify.com/pkg/ids"
)

func (c *Client) ListProjects(
	ctx context.Context,
	orgID ids.OrgID,
) ([]*projectsv1alpha1.Project, error) {
	memberResponse, err := c.ProjectsClient().ListProjects(
		ctx, connect.NewRequest(
			&projectsv1alpha1.ListProjectsRequest{
				OrgId: orgID.String(),
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return memberResponse.Msg.Projects, nil
}

func (c *Client) CreateProject(
	ctx context.Context,
	orgID ids.OrgID,
	repoURL string,
	directory string,
	name string,
) (*projectsv1alpha1.Project, error) {
	memberResponse, err := c.ProjectsClient().CreateProject(
		ctx, connect.NewRequest(
			&projectsv1alpha1.CreateProjectRequest{
				OrgId: orgID.String(),
				Project: &projectsv1alpha1.Project{
					Repo:      repoURL,
					Directory: directory,
					Name:      name,
				},
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return memberResponse.Msg.Project, nil
}
