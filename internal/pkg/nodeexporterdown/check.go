package nodeexporterdown

import (
	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/runbook/pkg/problem"
)

func (r *Runbook) investigate() (*problemData, error) {
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
	r.logger.Log("level", "debug", "message", "checking existing endpoint IP addresses")
	if len(endpointsList.Items) > 0 {
		// endpoints = endpointsList.Items[0]
		for _, endpoints := range endpointsList.Items {
			if endpoints.ObjectMeta.Name != "node-exporter" {
				continue
			}

			for _, subset := range endpoints.Subsets {
				for _, address := range subset.Addresses {
					e2nAddressMap[address.IP] = nil
					r.logger.Log("level", "debug", "foundEndpointIP", address.IP)
				}
			}
		}
	}

	// get all nodes
	r.logger.Log("level", "debug", "message", "checking node IP addresses")
	nodes, err := r.k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// get all node addresses and try to map them to endpoint addresses
	for _, node := range nodes.Items {
		nodeAddress := node.ObjectMeta.Labels["ip"]
		r.logger.Log("level", "debug", "foundNodeIP", nodeAddress)

		if _, ok := e2nAddressMap[nodeAddress]; ok {
			// found endpoint address for node address
			e2nAddressMap[nodeAddress] = &nodeAddress
			r.logger.Log("level", "debug", "correctEndpointIP", nodeAddress)
		} else {
			// node address is missing in endpoints list
			missingAddresses = append(missingAddresses, nodeAddress)
			r.logger.Log("level", "debug", "missingEndpointIP", nodeAddress)
		}
	}

	// finally, find all remaining endpoint addresses that do not have a corresponding node address
	for endpointAddress, nodeAddress := range e2nAddressMap {
		if nodeAddress == nil {
			// endpoint address does not have a corresponding node address
			staleAddresses = append(staleAddresses, endpointAddress)
			r.logger.Log("level", "debug", "staleEndpointIP", endpointAddress)
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
		r.logger.Log("level", "debug", "problem", problemStaleEndpoints.ID)
	}

	if len(missingAddresses) > 0 {
		// we need to trigger endpoint addition
		problemData.problems = append(problemData.problems, problemMissingEndpoints)
		r.logger.Log("level", "debug", "problem", problemMissingEndpoints.ID)
	}

	return &problemData, nil
}
