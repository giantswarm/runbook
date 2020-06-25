package daemonsetnotsatisfied

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/runbook/pkg/quay"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (r *Runbook) investigate() error {
	if r.inputs["name"] == "" {
		return microerror.Maskf(invalidConfigError, "DaemonSet name must not be empty")
	}
	if r.inputs["namespace"] == "" {
		return microerror.Maskf(invalidConfigError, "DaemonSet namespace must not be empty")
	}
	name := r.inputs["name"]
	namespace := r.inputs["namespace"]

	var err error
	var daemonSet *appsv1.DaemonSet
	{
		daemonSet, err = r.getDaemonSet(namespace, name)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var pods []corev1.Pod
	{
		podList, err := r.getPodsForDaemonSet(daemonSet)
		if err != nil {
			return microerror.Mask(err)
		} else {
			pods = podList.Items
		}
	}

	if daemonSet.Status.CurrentNumberScheduled < daemonSet.Status.DesiredNumberScheduled {
		notReadyPods := listUnreadyPods(pods)
		if int32(len(pods)) == daemonSet.Status.DesiredNumberScheduled && len(notReadyPods) > 0 {
			// We are most likely in the context of https://github.com/giantswarm/giantswarm/issues/8905
			err := r.fixIncorrectStatusReportedByKubelet(daemonSet, notReadyPods)
			if err != nil {
				return microerror.Mask(err)
			}
		} else {
			if r.isAnyContainerInImagePullBackOff(notReadyPods) {
				isQuayDown, err := quay.IsQuayDown()
				if err != nil {
					return microerror.Mask(err)
				}

				if isQuayDown {
					r.fixQuayDown(daemonSet, notReadyPods)
				} else {
					// TODO No Clue what to do
				}
			} else if r.arePodsUnschedulable(notReadyPods) {
				/*
					It is possible that pods cannot be scheduled because of inadequate CPU or memory on nodes.
					Ensure that nodes are big enough to hold all core component pods as well as some workload pods and that priority classes are set appropriately.
				*/
			} else if r.isHostPortInConflict(notReadyPods) {
				/*
					Port conflicts
					Some pods such as metrics exporters run with hostPort.
					If another pod tries to start using the same host port, it will fail and can lead to this alert if it’s part of a DaemonSet.
					A pod can even interfere with itself if it becomes a zombie and doesn’t release its TCP connection.
					We find this can happen occasionally for node-exporter and is solved by killing the zombie process. See PM.
				*/
			}
			// TODO What to do with unknow territory
			return microerror.Mask(err)
		}
	} else {
		// TODO this is a false alert or it resolved already. I am not sure how to report an issue
		return microerror.Mask(err)
	}
	return nil
}

func (r *Runbook) isAnyContainerInImagePullBackOff(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}

func (r *Runbook) arePodsUnschedulable(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}

func (r *Runbook) isHostPortInConflict(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}

func (r *Runbook) getDaemonSet(namespace string, name string) (*appsv1.DaemonSet, error) {
	return r.k8sClient.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
}

func (r *Runbook) getPodsForDaemonSet(daemonset *appsv1.DaemonSet) (*corev1.PodList, error) {
	options := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(daemonset.Spec.Selector.MatchLabels).String(),
	}

	return r.k8sClient.CoreV1().Pods(daemonset.ObjectMeta.Namespace).List(options)
}

func listUnreadyPods(pods []corev1.Pod) []corev1.Pod {
	var notReadyPods []corev1.Pod
	for _, pod := range pods {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status != "True" {
				notReadyPods = append(notReadyPods, pod)
			}
		}
	}
	return notReadyPods
}
