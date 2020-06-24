package nodeexporterdown

import (
	v1 "k8s.io/api/core/v1"

	"github.com/giantswarm/runbook/pkg/problem"
)

var problemStaleEndpoints = problem.Kind{
	ID: "StaleEndpoints",
	Description: "In KVM environments in older k8s versions (<1.12) frequent reason that node-exporter endpoints were" +
		" not updated properly on VM reinstall (e.g. after cluster upgrade). Stale endpoint needs to be removed.",
}

var problemMissingEndpoints = problem.Kind{
	ID: "MissingEndpoints",
	Description: "In KVM environments in older k8s versions (<1.12) frequent reason that node-exporter endpoints were" +
		" not updated properly on VM reinstall (e.g. after cluster upgrade). To trigger endpoint addition, we can" +
		" recreate node-exporter pods by removing them all.",
}

type problemData struct {
	problemKind   problem.Kind
	endpoints     v1.Endpoints
	nodeAddresses []string
}
