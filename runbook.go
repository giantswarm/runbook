package runbook

import (
	"github.com/giantswarm/microerror"
	runbooks "github.com/giantswarm/runbook/internal/pkg"
	runbookconfig "github.com/giantswarm/runbook/pkg/runbook/config"
)

func New(ID string, config runbookconfig.RunbookConfig) (Runbook, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}

	if config.Context == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Context must not be empty")
	}

	if ID == runbooks.DaemonSetNotSatisfiedRunbookID {
		return runbooks.NewDaemonSetNotSatisfiedRunbook(config), nil
	}

	return nil, nil
}
