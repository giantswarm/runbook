package daemonsetnotsatisfied

import (
	"errors"
	"fmt"

	"github.com/giantswarm/microerror"
)

func (r *Runbook) fixIncorrectStatusReportedByKubelet(data *problemData) error {
	// See https://github.com/giantswarm/giantswarm/issues/8905
	for _, pod := range data.pods {
		r.logger.Log("level", "info", "message", fmt.Sprintf("Restarting the kubelet for %s", pod.Spec.NodeName))
		err := r.runner.Run(pod.Spec.NodeName, "sudo systemctl restart k8s-kubelet")
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (r *Runbook) fixQuayDown(data *problemData) error {
	// TODO https://github.com/giantswarm/giantswarm/blob/master/content/docs/support-and-ops/ops-recipes/quay-is-down.md
	return errors.New("Not implemented :" + quayIsDown.Description)
}

func (r *Runbook) fixPodsStuckInCrashLoopBackOff(data *problemData) error {
	/*
		This could be multiple reasons:
		- If the installation is in China, is Alyun UP?
		- Is there any network connectivity issue with the registry?
		- Does the tag actually exists?
	*/
	return errors.New("Not implemented :" + podsStuckInCrashLoopBackOff.Description)
}

func (r *Runbook) fixPodsCanNotBeScheduled(data *problemData) error {
	return errors.New("Not implemented :" + podsCanNotBeScheduled.Description)
}

func (r *Runbook) fixHostPortInConflict(data *problemData) error {
	return errors.New("Not implemented :" + hostPortInConflict.Description)
}
