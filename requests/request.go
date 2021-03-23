package requests

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/choria-io/go-external/agent"
	"github.com/levigross/grequests"
)

type RequestRequest struct {
	CommonRequest

	URL          string `json:"url"`
	BodyFile     string `json:"body_file"`
	Body         string `json:"body"`
	ExpectedCode int    `json:"statuscode"`
	Method       string `json:"method"`
}

type RequestResponse struct {
	Body       string              `json:"body"`
	StatusCode int                 `json:"statuscode"`
	Headers    map[string][]string `json:"headers"`
	Duration   float64             `json:"duration"`
}

func RequestAction(request *agent.Request, reply *agent.Reply, _ map[string]string) {
	gr := &RequestRequest{}
	if !request.ParseRequestData(gr, reply) {
		return
	}

	if gr.URL == "" {
		reply.Abort("URL is required")
		return
	}
	if gr.Method == "" {
		gr.Method = "GET"
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout-time.Second)
	defer cancel()

	opts, err := requestOptions(ctx, &gr.CommonRequest)
	if reply.AbortIfErr(err, "%s", err) {
		return
	}

	gresp := &RequestResponse{StatusCode: 500}
	reply.Data = gresp

	switch {
	case gr.Body != "":
		opts.RequestBody = strings.NewReader(gr.Body)
	case gr.BodyFile != "":
		f, err := os.Open(gr.BodyFile)
		if reply.AbortIfErr(err, "%s", err) {
			return
		}
		defer f.Close()
		opts.RequestBody = f
	}

	start := time.Now()
	resp, err := grequests.Req(gr.Method, gr.URL, opts)
	if reply.AbortIfErr(err, "Request error: %s", err) {
		return
	}
	if reply.AbortIfErr(resp.Error, "Request error: %s", resp.Error) {
		return
	}

	gresp.StatusCode = resp.StatusCode
	gresp.Headers = resp.Header
	gresp.Body = resp.String()
	gresp.Duration = time.Since(start).Seconds()

	if gr.ExpectedCode > 0 && resp.StatusCode != gr.ExpectedCode {
		reply.Abort("Request failed: code %d", resp.StatusCode)
		return
	}
}
