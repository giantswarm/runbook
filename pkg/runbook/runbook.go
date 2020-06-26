package runbook

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/runbook/internal/pkg/daemonsetnotsatisfied"
	"github.com/giantswarm/runbook/internal/pkg/nodeexporterdown"
)

func New(config Config) (Interface, error) {
	var err error

	// dependencies
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	// internals
	if len(config.ID) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.ID must not be empty", config)
	}

	if config.Input == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Input must not be empty", config)
	}

	var runbook Interface

	switch config.ID {
	case nodeexporterdown.RunbookID:
		runbook, err = newNodeExporterRunbook(config)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	case daemonsetnotsatisfied.RunbookID:
		runbook, err = newDaemonSetNotSatisfiedRunbook(config)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	default:
		return nil, microerror.Maskf(invalidConfigError, "Runbook with ID '%s' not found", config.ID)
	}

	return runbook, nil
}

func newNodeExporterRunbook(config Config) (Interface, error) {
	c := nodeexporterdown.Config{
		Logger:    config.Logger,
		K8sClient: config.K8sClient,
		Input:     nodeexporterdown.Input(config.Input),
	}
	runbook, err := nodeexporterdown.NewRunbook(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return runbook, nil
}

func newDaemonSetNotSatisfiedRunbook(config Config) (Interface, error) {
	c := daemonsetnotsatisfied.Config{
		Logger:    config.Logger,
		K8sClient: config.K8sClient,
		Input:     daemonsetnotsatisfied.Input(config.Input),
	}
	runbook, err := daemonsetnotsatisfied.NewDaemonSetNotSatisfiedRunbook(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return runbook, nil
}
