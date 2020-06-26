package nodeexporterdown

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/giantswarm/microerror"
)

const (
	nodeExporterNamespace = "kube-system"
	appLabel              = "app"
)

func (r *Runbook) fixStaleEndpoint(data *problemData) error {
	staleAddresses := data.staleAddresses
	var endpointAddresses []corev1.EndpointAddress
	e := &data.endpoints

	for _, endpointAddress := range e.Subsets[0].Addresses {
		if !contains(staleAddresses, endpointAddress.IP) {
			// keep only not stale addresses
			endpointAddresses = append(endpointAddresses, endpointAddress)
		}
	}

	e.Subsets[0].Addresses = endpointAddresses
	_, err := r.k8sClient.CoreV1().Endpoints(nodeExporterNamespace).Update(e)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Runbook) fixMissingEndpoint(data *problemData) error {
	appLabelValue := data.endpoints.Labels[appLabel]
	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			appLabel: appLabelValue, // app: node-exporter
		},
	}

	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector.MatchLabels).String(),
	}

	// by deleting node-exporter pods, which will then be recreated, we are trigger endpoint addition
	err := r.k8sClient.CoreV1().Pods(nodeExporterNamespace).DeleteCollection(nil, listOptions)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func contains(a []string, s string) bool {
	for _, val := range a {
		if val == s {
			return true
		}
	}

	return false
}
