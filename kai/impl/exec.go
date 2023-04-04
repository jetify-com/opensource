package impl

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	kai_v1 "go.jetpack.io/kai/gen/api/kai/v1"
	"go.jetpack.io/kai/gen/api/kai/v1/kai_v1connect"
)

func Exec(query string) ([]string, error) {
	client := kai_v1connect.NewKaiServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	resp, err := client.GetShellCommand(
		context.Background(),
		connect.NewRequest(&kai_v1.GetShellCommandRequest{
			Prompt: query,
		}),
	)
	if err != nil {
		return []string{}, err
	}

	results := resp.Msg.GetChoices()
	return results, nil
}
