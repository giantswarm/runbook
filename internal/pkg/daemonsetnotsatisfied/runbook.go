package daemonsetnotsatisfied

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	runbookconfig "github.com/giantswarm/runbook/pkg/runbook/config"
)

const (
	runbookId        = "daemonset-not-satisfied"
	runbookSourceURL = "https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/daemonset-not-satisfied/"
)

type Runbook struct {
	logger    micrologger.Logger
	k8sClient kubernetes.Interface
	inputs    map[string]string
}

func NewDaemonSetNotSatisfiedRunbook(config runbookconfig.RunbookConfig) (*Runbook, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}

	if config.Inputs == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Inputs must not be empty")
	}
	runbook := Runbook{
		logger:    config.Logger,
		k8sClient: config.K8sClient,
		inputs:    config.Inputs,
	}

	return &runbook, nil
}

func (r *Runbook) GetID() string {
	return runbookId
}

func (r *Runbook) GetSourceURL() string {
	return runbookSourceURL
}

// We check if the daemonset actually exists or not
func (r *Runbook) Test() (bool, error) {
	r.investigate()
	return true, nil
}

func (r *Runbook) Apply() error {
	r.investigate()
	return nil
}
