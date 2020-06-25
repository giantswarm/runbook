package daemonsetnotsatisfied

import (
	"errors"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/runbook/pkg/problem"
	runbookconfig "github.com/giantswarm/runbook/pkg/runbook/config"
)

const (
	runbookId        = "daemonset-not-satisfied"
	runbookSourceURL = "https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/daemonset-not-satisfied/"
)

type Runbook struct {
	logger    micrologger.Logger
	k8sClient kubernetes.Interface
	input     runbookconfig.RunbookInput
}

func NewDaemonSetNotSatisfiedRunbook(config runbookconfig.RunbookConfig) (*Runbook, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}

	if config.Input == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Input must not be empty")
	}
	runbook := Runbook{
		logger:    config.Logger,
		k8sClient: config.K8sClient,
		input:     config.Input,
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
	observations, err := r.investigate()
	if err != nil {
		return false, microerror.Mask(err)
	}

	return problem.IsFound(observations.problem), nil
}

func (r *Runbook) Apply() error {
	observations, err := r.investigate()
	if err != nil {
		return microerror.Mask(err)
	}

	switch observations.problem.ID {
	case incorrectStatusReportedByKubelet.ID:
		return r.fixIncorrectStatusReportedByKubelet(observations)
	case quayIsDown.ID:
		return r.fixQuayDown(observations)
	case podsStuckInCrashLoopBackOff.ID:
		return r.fixPodsStuckInCrashLoopBackOff(observations)
	case podsCanNotBeScheduled.ID:
		return r.fixPodsCanNotBeScheduled(observations)
	case hostPortInConflict.ID:
		return r.fixHostPortInConflict(observations)
	case problem.Unknown.ID:
		return errors.New(problem.Unknown.Description)
	default:
		return nil
	}
}
