package poller

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("workload-poller")

// Poller represents the poller with a private key and a handler
type Poller struct {
	identity *utils.Identity
	sampler  *sampler.Sampler
	endpoint string
}

// NewPoller creates a new Poller with the provided private key in hex format
func NewPoller(identity *utils.Identity, sampler *sampler.Sampler, endpoint string) (*Poller, error) {
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return &Poller{
		identity: identity,
		sampler:  sampler,
		endpoint: endpoint,
	}, nil
}

func (p *Poller) Start() {
	go p.periodicPoll()

	// Wait forever
	p.waitForShutdown()
}

func (p *Poller) periodicPoll() {
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

		log.Infof("response: %s", string(body))

		time.Sleep(10 * time.Second)
	}
}

func (l *Poller) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
