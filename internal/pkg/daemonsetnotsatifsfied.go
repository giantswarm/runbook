package pkg

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	runbookconfig "github.com/giantswarm/runbook/pkg/runbook/config"
)

const (
	DaemonSetNotSatisfiedRunbookID        = "daemonset-not-satisfied"
	DaemonSetNotSatisfiedRunbookSourceURL = "https://intranet.giantswarm.io/docs/support-and-ops/ops-recipes/daemonset-not-satisfied/"
)

func NewDaemonSetNotSatisfiedRunbook(config runbookconfig.RunbookConfig) *DaemonSetNotSatisfiedRunbook {
	return &DaemonSetNotSatisfiedRunbook{
		logger:    config.Logger,
		k8sClient: config.K8sClient,
		context:   config.Context,
	}
}

type DaemonSetNotSatisfiedRunbook struct {
	logger    micrologger.Logger
	k8sClient kubernetes.Interface
	context   map[string]string
}

func (r DaemonSetNotSatisfiedRunbook) GetID() string {
	return DaemonSetNotSatisfiedRunbookID
}

func (r DaemonSetNotSatisfiedRunbook) GetSourceURL() string {
	return DaemonSetNotSatisfiedRunbookSourceURL
}

// We check if the daemonset actually exists or not
func (r DaemonSetNotSatisfiedRunbook) Test() (bool, error) {
	if r.context["name"] == "" {
		return false, microerror.Maskf(invalidConfigError, "DaemonSet name must not be empty")
	}
	if r.context["namespace"] == "" {
		return false, microerror.Maskf(invalidConfigError, "DaemonSet namespace must not be empty")
	}
	name := r.context["name"]
	namespace := r.context["namespace"]

	_, err := r.getDaemonSet(namespace, name)
	if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func (r DaemonSetNotSatisfiedRunbook) Apply() error {
	if r.context["name"] == "" {
		return microerror.Maskf(invalidConfigError, "DaemonSet name must not be empty")
	}
	if r.context["namespace"] == "" {
		return microerror.Maskf(invalidConfigError, "DaemonSet namespace must not be empty")
	}
	name := r.context["name"]
	namespace := r.context["namespace"]

	var err error
	var daemonSet *appsv1.DaemonSet
	{
		daemonSet, err = r.getDaemonSet(namespace, name)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var pods *corev1.PodList
	options := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(daemonSet.Spec.Selector.MatchLabels).String(),
	}
	if pods, err = r.k8sClient.CoreV1().Pods(name).List(options); err != nil {
		return microerror.Mask(err)
	}

	if daemonSet.Status.CurrentNumberScheduled < daemonSet.Status.DesiredNumberScheduled {
		notReadyPods := listUnreadyPods(pods.Items)
		if int32(len(pods.Items)) == daemonSet.Status.DesiredNumberScheduled && len(notReadyPods) > 0 {
			// We are most likely in the context of https://github.com/giantswarm/giantswarm/issues/8905
			err := fixIncorrectStatusReportedByKubelet(daemonSet, notReadyPods)
			if err != nil {
				return microerror.Mask(err)
			}
		} else {
			if isAnyContainerInImagePullBackOff(notReadyPods) {
				if isQuayDown() {
					// TODO https://github.com/giantswarm/giantswarm/blob/master/content/docs/support-and-ops/ops-recipes/quay-is-down.md
				} else {
					// TODO No Clue what to do
				}
			} else if arePodsUnschedulable(notReadyPods) {
				/*
					It is possible that pods cannot be scheduled because of inadequate CPU or memory on nodes.
					Ensure that nodes are big enough to hold all core component pods as well as some workload pods and that priority classes are set appropriately.
				*/
			} else if isHostPortInConflict(notReadyPods) {
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

func (r DaemonSetNotSatisfiedRunbook) getDaemonSet(namespace string, name string) (*appsv1.DaemonSet, error) {
	return r.k8sClient.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
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

func isAnyContainerInImagePullBackOff(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}

func isQuayDown() bool {
	// TODO implement this
	return false
}

func arePodsUnschedulable(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}

func isHostPortInConflict(pods []corev1.Pod) bool {
	// TODO implement this
	return false
}

func fixIncorrectStatusReportedByKubelet(daemonSet *appsv1.DaemonSet, notReadyPods []corev1.Pod) error {
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
