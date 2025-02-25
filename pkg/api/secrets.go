package api

import (
	"go.jetify.com/pkg/api/gen/priv/secrets/v1alpha1/secretsv1alpha1connect"
)

func (c *Client) SecretsService() secretsv1alpha1connect.SecretsServiceClient {
	return c.secretsClient()
}
