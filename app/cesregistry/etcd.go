package cesregistry

import (
	"fmt"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-ces-setup/app/setup/component"
)

// CreateEtcd returns a new etcd registry pointing to the specified namespace.
func CreateEtcd(namespace string) (registry.Registry, error) {
	return registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://%s:4001", component.GetNodeMasterFileContent(namespace))},
	})
}
