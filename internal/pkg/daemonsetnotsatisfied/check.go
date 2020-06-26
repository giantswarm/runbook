package daemonsetnotsatisfied

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/runbook/pkg/kubernetes/util"
	"github.com/giantswarm/runbook/pkg/problem"
	"github.com/giantswarm/runbook/pkg/quay"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (r *Runbook) investigate() (*problemData, error) {
	if r.input["name"] == "" {
		return nil, microerror.Maskf(invalidConfigError, "DaemonSet name must not be empty")
	}
	if r.input["namespace"] == "" {
		return nil, microerror.Maskf(invalidConfigError, "DaemonSet namespace must not be empty")
	}
	name := r.input["name"]
	namespace := r.input["namespace"]

	var err error
	var daemonSet *appsv1.DaemonSet
	{
		daemonSet, err = util.GetDaemonSet(r.k8sClient, namespace, name)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var pods []corev1.Pod
	pods, err = util.GetPodsMatchingLabelsDaemonSet(r.k8sClient, daemonSet.ObjectMeta.Namespace, labels.SelectorFromSet(daemonSet.Spec.Selector.MatchLabels).String())
	if err != nil {
		return nil, microerror.Mask(err)
	}

	problemData := problemData{
		problem:   problem.Kind{},
		daemonSet: daemonSet,
	}

	r.logger.Log("level", "debug", "message", "checking if the number of currently scheduled pods is inferior to the number of desired pods")
	if daemonSet.Status.CurrentNumberScheduled < daemonSet.Status.DesiredNumberScheduled {
		notReadyPods := util.FilterNotReadyPods(pods)

		r.logger.Log("level", "debug", "message", "checking if the kubelet is correctly reporting its status")
		if int32(len(pods)) == daemonSet.Status.DesiredNumberScheduled && len(notReadyPods) > 0 {
			problemData.problem = incorrectStatusReportedByKubelet
			problemData.pods = notReadyPods
			return &problemData, nil
		}

		r.logger.Log("level", "debug", "message", "checking if any pod is in CrashLoopBackOff")
		if podsWithContainersInImagePullBackOff := util.FilterPodsWithContainersInImagePullBackOff(notReadyPods); len(podsWithContainersInImagePullBackOff) > 0 {

			isQuayDown, err := quay.IsQuayDown()
			if err != nil {
				return nil, microerror.Mask(err)
			}

			if isQuayDown {
				problemData.problem = quayIsDown
			} else {
				problemData.problem = podsStuckInCrashLoopBackOff
			}
			problemData.pods = podsWithContainersInImagePullBackOff
			return &problemData, nil
		}

		r.logger.Log("level", "debug", "message", "checking if any pod is Unschedulable")
		if unschedulablePods := util.FilterUnschedulablePods(notReadyPods); len(unschedulablePods) > 0 {
			problemData.problem = podsCanNotBeScheduled
			problemData.pods = unschedulablePods
			return &problemData, nil
		}

		r.logger.Log("level", "debug", "message", "checking if any pod has any host port conflict")
		if podsWithHostPortConflicts := util.FilterPodsWithHostPortConflict(notReadyPods); r.isHostPortInConflict(notReadyPods) {
			problemData.problem = hostPortInConflict
			problemData.pods = podsWithHostPortConflicts
			return &problemData, nil
		}

		r.logger.Log("level", "debug", "message", "The source of the problem is unknown")
		problemData.problem = problem.Unknown
		return &problemData, nil
	}

	r.logger.Log("level", "debug", "message", "No problem was detected")
	problemData.problem = problem.None
	return &problemData, nil
}

func (r *Runbook) areSomePodsNotScheduled(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}

func (r *Runbook) isHostPortInConflict(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}
