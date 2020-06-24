package runbook

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	runbooks "github.com/giantswarm/runbook/internal/pkg"
)

type RunbookConfig struct {
	Logger       micrologger.Logger
	Installation string
	ClusterID    string
	Context      map[string]string
}

func New(ID string, config RunbookConfig) (Runbook, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Installation must not be empty")
	}

	if config.ClusterID == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ClusterID must not be empty")
	}

	if ID == runbooks.DaemonSetNotSatisfiedRunbookID {
		return &runbooks.DaemonSetNotSatisfiedRunbook{
			Logger:       config.Logger,
			Installation: config.Installation,
			ClusterID:    config.ClusterID,
			Context:      config.Context,
		}, nil
	}

	return nil, nil
}
