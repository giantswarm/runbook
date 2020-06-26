package runbook

import (
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
)

type Input map[string]string

type Config struct {
	Logger    micrologger.Logger
	K8sClient kubernetes.Interface

	ID    string
	Input Input
}
