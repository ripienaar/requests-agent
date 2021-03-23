package requests

import (
	"context"
	"time"

	"github.com/levigross/grequests"
)

const timeout = 120 * time.Second

type CommonRequest struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Query    map[string]string `json:"query"`
	Headers  map[string]string `json:"headers"`
}

func requestOptions(ctx context.Context, common *CommonRequest) (*grequests.RequestOptions, error) {
	opts := &grequests.RequestOptions{
		RedirectLimit: 10,
		UserAgent:     "Choria Requests Agent/0.0.1",
		Context:       ctx,
	}

	if common.Username != "" {
		opts.Auth = []string{common.Username, common.Password}
	}
	if len(common.Query) > 0 {
		opts.Params = common.Query
	}
	if len(common.Headers) > 0 {
		opts.Headers = common.Headers
	}

	return opts, nil
}
