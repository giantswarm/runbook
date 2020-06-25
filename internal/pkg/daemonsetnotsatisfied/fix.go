package daemonsetnotsatisfied

import "errors"

func (r *Runbook) fixIncorrectStatusReportedByKubelet(data *problemData) error {
	var affectedNodes []string
	for _, pod := range data.pods {
		affectedNodes = append(affectedNodes, pod.Spec.NodeName)
	}
	// TODO restart the kubelets on the affected nodes
	/*
		Use opsctl ssh to get to the node where the misreported pod is running.
		Restart the kubelet on this node (sudo systemctl restart k8s-kubelet).
		The DaemonSet should be satisfied once this is done for all affected nodes
	*/
	return errors.New("Not implemented :" + incorrectStatusReportedByKubelet.Description)
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
