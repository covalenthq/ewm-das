package poller

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/apihandler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/sampler"
	"github.com/covalenthq/das-ipfs-pinner/internal/light-client/utils"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("workload-poller")

// WorkloadPoller represents the poller with a private key and a handler
type WorkloadPoller struct {
	identity *utils.Identity
	sampler  *sampler.Sampler
	api      *apihandler.ApiHandler
}

// NewWorkloadPoller creates a new Poller with the provided private key in hex format
func NewWorkloadPoller(identity *utils.Identity, sampler *sampler.Sampler, api *apihandler.ApiHandler) *WorkloadPoller {
	return &WorkloadPoller{
		identity: identity,
		sampler:  sampler,
		api:      api,
	}
}

func (p *WorkloadPoller) Start() {
	go p.periodicPoll()

	// Wait forever
	p.waitForShutdown()
}

func (p *WorkloadPoller) periodicPoll() {
	for {

		response, err := p.api.GetWorkload()
		if err != nil {
			log.Errorf("failed to get workload: %s", err)
			time.Sleep(60 * time.Second)
			continue
		}

		// Process the workloads

		for _, workload := range response.Workloads {
			log.Debugf("processing workload: %v", workload.GetWorkload().ReadableString())
			challenge, err := Decode(workload.Workload.Challenge)
			if err != nil {
				log.Errorf("failed to decode challenge: %s", err)
			}

			eligible, seed, err := challenge.Solve(workload.Workload, p.identity)
			if err != nil {
				log.Errorf("failed to solve challenge: %s", err)
			}

			log.Debugf("workload is eligible: %v", eligible)
			if eligible {
				log.Infof("workload for CID %s is eligible", workload.GetWorkload().GetCID().String())
				p.sampler.ProcessEvent(workload, seed)
			}
		}

		log.Infof("%d workloads available", len(response.Workloads))

		// unix timestamp to time
		nextUpdate := time.Unix(int64(response.NextUpdateTimestamp), 0)

		log.Infof("waiting for next update: %v in %f seconds", nextUpdate, time.Until(nextUpdate).Seconds())
		time.Sleep(time.Until(nextUpdate))
	}
}

func (l *WorkloadPoller) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
