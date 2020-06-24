package nodeexporterdown

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/runbook/pkg/problem"
)

func (r *Runbook) getProblemData() (*problemData, error) {
	var endpoints corev1.Endpoints
	var endpointAddresses []string
	{

		endpointsList, err := r.k8sClient.CoreV1().Endpoints("kube-system").List(metav1.ListOptions{})
		if err != nil {
			return nil, microerror.Mask(err)
		}

		if len(endpointsList.Items) > 0 {
			endpoints = endpointsList.Items[0];
			for _, address := range endpointsList.Items[0].Subsets[0].Addresses {
				endpointAddresses = append(endpointAddresses, address.IP)
			}
		}
	}

	var nodeAddresses []string
	{
		nodes, err := r.k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			return nil, microerror.Mask(err)
		}

		for _, node := range nodes.Items {
			nodeAddress := node.ObjectMeta.Labels["ip"]
			nodeAddresses = append(nodeAddresses, nodeAddress)
		}
	}

	var problemKind problem.Kind

	switch {
	case len(endpointAddresses) > len(nodeAddresses):
		problemKind = problemStaleEndpoints
	case len(endpointAddresses) < len(nodeAddresses):
		problemKind = problemMissingEndpoints
	default:
		problemKind = problem.None
	}

	problemData := problemData{
		problemKind:   problemKind,
		endpoints:     endpoints,
		nodeAddresses: nodeAddresses,
	}

	return &problemData, nil
}
