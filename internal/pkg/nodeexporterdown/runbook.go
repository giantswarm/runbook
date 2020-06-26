package nodeexporterdown

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/runbook/pkg/problem"
)

const (
	RunbookID        = "node-exporter-is-down"
	runbookSourceURL = "https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/node-exporter-is-down/"
)

type Input map[string]string

type Config struct {
	Logger    micrologger.Logger
	K8sClient kubernetes.Interface

	Input Input
}

type Runbook struct {
	logger    micrologger.Logger
	k8sClient kubernetes.Interface
	input     Input
}

func NewRunbook(config Config) (*Runbook, error) {
	// dependencies
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	// internals
	if config.Input == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Input must not be empty", config)
	}

	runbook := Runbook{
		logger:    config.Logger,
		k8sClient: config.K8sClient,
		input:     config.Input,
	}

	return &runbook, nil
}

func (r *Runbook) GetID() string {
	return RunbookID
}

func (r *Runbook) GetSourceURL() string {
	return runbookSourceURL
}

func (r *Runbook) FindProblem() ([]problem.Kind, error) {
	data, err := r.investigate()
	if err != nil {
		return []problem.Kind{problem.Unknown}, microerror.Mask(err)
	}

	return data.problems, nil
}

func (r *Runbook) Test() (bool, error) {
	data, err := r.investigate()
	if err != nil {
		return false, microerror.Mask(err)
	}

	problemFound := problem.IsFound(data.problems...)
	return problemFound, nil
}

func (r *Runbook) Apply() error {
	data, err := r.investigate()
	if err != nil {
		return microerror.Mask(err)
	}

	for _, p := range data.problems {
		switch p.ID {
		case problemStaleEndpoints.ID:
			return r.fixStaleEndpoint(data)
		case problemMissingEndpoints.ID:
			return r.fixMissingEndpoint(data)
		}
	}

	return nil
}
