package daemonsetnotsatisfied

import (
	"github.com/giantswarm/runbook/pkg/problem"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var incorrectStatusReportedByKubelet = problem.Kind{
	ID:          "IncorrectStatusReportedByTheKubelet",
	Description: "In some versions of Kubernetes, the Kubelet failed to update the status of daemonset pods. For more information, see https://github.com/giantswarm/giantswarm/issues/8905.",
}

var quayIsDown = problem.Kind{
	ID:          "QuayIsDown",
	Description: "Quay appears to be down at the moment. Check the status page at https://status.quay.io/ for more details.",
}

var podsStuckInCrashLoopBackOff = problem.Kind{
	ID:          "PodsStuckInCrashLoopBackoff",
	Description: "Some of the pods in the daemon set are stuck in CrashLoopBackoff. As Quay seems to be up, this is unknown territory. Good luck.",
}

var podsCanNotBeScheduled = problem.Kind{
	ID:          "UnschedulablePods",
	Description: "Some of the pods cannot be scheduled to the nodes. Check the resources of the impacted nodes.",
}

var hostPortInConflict = problem.Kind{
	ID: "HostPortInConflict",
	Description: "Some pods such as metrics exporters run with hostPort." +
		" If another pod tries to start using the same host port, it will fail and can lead to this alert if it’s part of a DaemonSet." +
		" A pod can even interfere with itself if it becomes a zombie and doesn’t release its TCP connection." +
		" We find this can happen occasionally for node-exporter and is solved by killing the zombie process." +
		" See https://github.com/giantswarm/giantswarm/issues/10034 for more details",
}

type problemData struct {
	problem   problem.Kind
	daemonSet *appsv1.DaemonSet
	pods      []corev1.Pod
}
