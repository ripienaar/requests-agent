package main

import (
	"github.com/choria-io/go-external/agent"

	"github.com/ripienaar/requests-agent/requests"
)

func main() {
	ra := agent.NewAgent("requests")
	defer ra.ProcessRequest()

	ra.MustRegisterAction("download", requests.DownloadAction)
	ra.MustRegisterAction("request", requests.RequestAction)
}
