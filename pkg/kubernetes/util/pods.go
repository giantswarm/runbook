package util

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const ImagePullBackOff = "ImagePullBackOff"

type podPredicate func(corev1.Pod) bool

func GetPodsMatchingLabelsDaemonSet(client kubernetes.Interface, namespace string, labelSelector string) ([]corev1.Pod, error) {
	options := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	podList, err := client.CoreV1().Pods(namespace).List(options)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return podList.Items, nil
}

func filterPods(pods []corev1.Pod, predicate podPredicate) []corev1.Pod {
	var newPods []corev1.Pod
	for _, pod := range pods {
		if predicate(pod) {
			newPods = append(newPods, pod)
		}
	}
	return newPods
}

func isPodReady(pod corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == "True"
		}
	}
	return false
}

func isPodNotReady(pod corev1.Pod) bool {
	return !isPodReady(pod)
}

func isPodNotScheduled(pod corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodScheduled {
			return condition.Status != "True"
		}
	}
	return false
}

func FilterNotReadyPods(pods []corev1.Pod) []corev1.Pod {
	return filterPods(pods, isPodNotReady)
}

func isAnyContainerInImagePullBackoff(pod corev1.Pod) bool {
	for _, containerStatuses := range pod.Status.ContainerStatuses {
		if containerStatuses.State.Waiting != nil && containerStatuses.State.Waiting.Reason == ImagePullBackOff {
			return true
		}
	}
	for _, containerStatuses := range pod.Status.InitContainerStatuses {
		if containerStatuses.State.Waiting != nil && containerStatuses.State.Waiting.Reason == ImagePullBackOff {
			return true
		}
	}
	// We ignore ephemeral containers for now
	return false
}

func FilterPodsWithContainersInImagePullBackOff(pods []corev1.Pod) []corev1.Pod {
	return filterPods(pods, isAnyContainerInImagePullBackoff)
}

func FilterUnschedulablePods(pods []corev1.Pod) []corev1.Pod {
	return filterPods(pods, isPodNotScheduled)
}

func FilterPodsWithHostPortConflict(pods []corev1.Pod) []corev1.Pod {
	return []corev1.Pod{} // TODO implement this
}
