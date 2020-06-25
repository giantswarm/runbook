package nodeexporterdown

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
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
	_, err := r.k8sClient.CoreV1().Endpoints("kube-system").Update(e)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Runbook) fixMissingEndpoint(data *problemData) error {
	return nil
}

func contains (a []string, s string) bool {
	for _, val := range a {
		if val == s {
			return true
		}
	}

	return false
}
