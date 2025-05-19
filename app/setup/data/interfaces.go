package data

import corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

type secretClient interface {
	corev1.SecretInterface
}
