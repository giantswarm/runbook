package runbookconfig

import (
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
)

type RunbookConfig struct {
	Logger    micrologger.Logger
	K8sClient kubernetes.Interface
	Inputs    map[string]string
}
