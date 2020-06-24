package pkg

import (
	"github.com/giantswarm/micrologger"
)

const (
	DaemonSetNotSatisfiedRunbookID        = "daemonset-not-satisfied"
	DaemonSetNotSatisfiedRunbookSourceURL = "https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/daemonset-not-satisfied/"
)

type DaemonSetNotSatisfiedRunbook struct {
	Logger       micrologger.Logger
	Installation string
	ClusterID    string
	Context      map[string]string
}

func (r DaemonSetNotSatisfiedRunbook) GetID() string {
	return DaemonSetNotSatisfiedRunbookID
}

func (r DaemonSetNotSatisfiedRunbook) GetSourceURL() string {
	return DaemonSetNotSatisfiedRunbookSourceURL
}

func (r DaemonSetNotSatisfiedRunbook) Test() (bool, error) {
	return false, nil
}

func (r DaemonSetNotSatisfiedRunbook) Apply() error {
	return nil
}
