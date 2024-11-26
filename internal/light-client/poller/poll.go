package poller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("workload-poller")

// WorkloadPoller represents the poller with a private key and a handler
type WorkloadPoller struct {
	identity *utils.Identity
	sampler  *sampler.Sampler
	endpoint string
}

// NewWorkloadPoller creates a new Poller with the provided private key in hex format
func NewWorkloadPoller(identity *utils.Identity, sampler *sampler.Sampler, endpoint string) (*WorkloadPoller, error) {
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return &WorkloadPoller{
		identity: identity,
		sampler:  sampler,
		endpoint: endpoint,
	}, nil
}

func (p *WorkloadPoller) Start() {
	go p.periodicPoll()

	// Wait forever
	p.waitForShutdown()
}

func (p *WorkloadPoller) periodicPoll() {
	for {

		// Poll the endpoint
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", p.endpoint, nil)
		if err != nil {
			log.Errorf("failed to create request: %s", err)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Errorf("failed to poll endpoint: %s", err)
		}

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("failed to read response: %s", err)
		}

		var response internal.WorkloadResponse
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			log.Errorf("failed to unmarshal response: %s", err)
		}

		// Process the workloads
		for _, workload := range response.Workloads {
			log.Infof("processing workload: %v", workload)
			challenge, err := Decode(workload.Challenge)
			if err != nil {
				log.Errorf("failed to decode challenge: %s", err)
			}

			eligible, err := challenge.Solve(&workload, p.identity)
			if err != nil {
				log.Errorf("failed to solve challenge: %s", err)
			}

			log.Infof("workload is eligible: %v", eligible)
			if eligible {
				p.sampler.ProcessEvent2(&workload)
			}
		}

		time.Sleep(time.Until(response.NextUpdate))
	}
}

func (l *WorkloadPoller) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
