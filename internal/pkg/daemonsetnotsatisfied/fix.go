package daemonsetnotsatisfied

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (r *Runbook) fixIncorrectStatusReportedByKubelet(daemonSet *appsv1.DaemonSet, notReadyPods []corev1.Pod) error {
	var affectedNodes []string
	for _, pod := range notReadyPods {
		affectedNodes = append(affectedNodes, pod.Spec.NodeName)
	}
	// TODO restart the kubelets on the affected nodes
	/*
		Use opsctl ssh to get to the node where the misreported pod is running.
		Restart the kubelet on this node (sudo systemctl restart k8s-kubelet).
		The DaemonSet should be satisfied once this is done for all affected nodes
	*/
	return nil
}

// TODO https://github.com/giantswarm/giantswarm/blob/master/content/docs/support-and-ops/ops-recipes/quay-is-down.md
func (r *Runbook) fixQuayDown(daemonSet *appsv1.DaemonSet, notReadyPods []corev1.Pod) error {
	return nil
}
