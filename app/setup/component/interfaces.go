package component

import "github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"

type componentsClient interface {
	ecosystem.ComponentInterface
}
