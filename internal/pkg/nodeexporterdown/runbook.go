package nodeexporterdown

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/runbook/pkg/problem"
	runbookconfig "github.com/giantswarm/runbook/pkg/runbook/config"
)

const (
	runbookID        = "node-exporter-is-down"
	runbookSourceURL = "https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/node-exporter-is-down/"
)

type alertContext map[string]string

type Runbook struct {
	logger    micrologger.Logger
	k8sClient kubernetes.Interface

	context alertContext
}

func NewRunbook(config runbookconfig.RunbookConfig) (*Runbook, error) {
	// dependencies
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	// internals
	if config.Context == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Context must not be empty", config)
	}

	runbook := Runbook{
		logger:    config.Logger,
		k8sClient: config.K8sClient,
		context:   config.Context,
	}

	return &runbook, nil
}

func (r *Runbook) GetID() string {
	return runbookID
}

func (r *Runbook) GetSourceURL() string {
	return runbookSourceURL
}

func (r *Runbook) FindProblem() ([]problem.Kind, error) {
	data, err := r.getProblemData()
	if err != nil {
		return []problem.Kind{problem.Unknown}, microerror.Mask(err)
	}

	return data.problems, nil
}

func (r *Runbook) Test() (bool, error) {
	data, err := r.getProblemData()
	if err != nil {
		return false, microerror.Mask(err)
	}

	problemFound := problem.IsFound(data.problems...)
	return problemFound, nil
}

func (r *Runbook) Apply() error {
	data, err := r.getProblemData()
	if err != nil {
		return microerror.Mask(err)
	}

	for _, p := range data.problems {
		switch p.ID {
		case problemStaleEndpoints.ID:
			err := r.fixStaleEndpoint(data)
			if err != nil {
				return microerror.Mask(err)
			}
		case problemMissingEndpoints.ID:
			err := r.fixMissingEndpoint(data)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}
