package util

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetDaemonSet(client kubernetes.Interface, namespace string, name string) (*appsv1.DaemonSet, error) {
	return client.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
}
