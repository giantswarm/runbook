package nodeexporterdown

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/runbook/pkg/problem"
)

func (r *Runbook) getProblemData() (*problemData, error) {
	// endpoint -> node address map
	e2nAddressMap := make(map[string]*string)

	var endpoints corev1.Endpoints
	var staleAddresses []string
	var missingAddresses []string

	// get endpoints CR
	endpointsList, err := r.k8sClient.CoreV1().Endpoints("kube-system").List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// get all endpoint addresses
	if len(endpointsList.Items) > 0 {
		endpoints = endpointsList.Items[0];
		for _, address := range endpointsList.Items[0].Subsets[0].Addresses {
			e2nAddressMap[address.IP] = nil
		}
	}

	// get all nodes
	nodes, err := r.k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// get all node addresses and try to map them to endpoint addresses
	for _, node := range nodes.Items {
		nodeAddress := node.ObjectMeta.Labels["ip"]
		if _, ok := e2nAddressMap[nodeAddress]; ok {
			// found endpoint address for node address
			e2nAddressMap[nodeAddress] = &nodeAddress
		} else {
			// node address is missing in endpoints list
			missingAddresses = append(missingAddresses, nodeAddress)
		}
	}

	// finally, find all remaining endpoint addresses that do not have a corresponding node address
	for endpointAddress, nodeAddress := range e2nAddressMap {
		if nodeAddress == nil {
			// endpoint address does not have a corresponding node address
			staleAddresses = append(staleAddresses, endpointAddress)
		}
	}

	problemData := problemData{
		problems:  []problem.Kind{},
		endpoints: endpoints,
	}

	if len(staleAddresses) > 0 {
		// we need to remove stale endpoints
		problemData.problems = append(problemData.problems, problemStaleEndpoints)
		problemData.staleAddresses = staleAddresses
	}

	if len(missingAddresses) > 0 {
		// we need to trigger endpoint addition
		problemData.problems = append(problemData.problems, problemMissingEndpoints)
	}

	return &problemData, nil
}
