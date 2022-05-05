package validation

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

type doguValidator struct {
	dogus    context.Dogus
	secret   corev1.Secret
	registry remote.Registry
}

func NewDoguValidator(doguRegistrySecret *corev1.Secret, dogus context.Dogus) (*doguValidator, error) {
	credentials := &core.Credentials{
		Username: doguRegistrySecret.StringData["username"],
		Password: doguRegistrySecret.StringData["password"],
	}

	// TODO ConfigMap für URlSchema (default oder mirrored)
	remoteConfig := &core.Remote{
		Endpoint:  doguRegistrySecret.StringData["endpoint"],
		URLSchema: "",
		CacheDir:  "/tmp",
	}

	registry, err := remote.New(remoteConfig, credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote registry: %w", err)
	}

	return &doguValidator{registry: registry, dogus: dogus}, nil
}

func (dv *doguValidator) ValidateDogus() error {
	for _, installDogu := range dv.dogus.Install {
		dogu, err := dv.getDogu(installDogu)
		if err != nil {
			return err
		}

		err = dv.validateDoguDependencies(dogu.GetDependenciesOfType("dogu"))
		if err != nil {
			return fmt.Errorf("failed to validate dependencies for dogu %s: %w", installDogu, err)
		}
	}

	return nil
}

func (dv *doguValidator) validateDoguDependencies(dependencies []core.Dependency) error {
	for _, dependency := range dependencies {
		if dependency.Type != "dogu" {
			continue
		}

		depName := dependency.Name
		_, err := dv.getDoguFromSelection(depName)
		if err != nil {
			return fmt.Errorf("dogu dependency %s ist not selected", depName)
		}

		_, err = core.ParseVersion(dependency.Version)
		if err != nil {
			return fmt.Errorf("failed to parse version from dependency %s: %w", depName, err)
		}

	}

	return nil
}

func (dv *doguValidator) getDoguFromSelection(dogu string) (string, error) {
	for _, installDogu := range dv.dogus.Install {
		// Works with version?
		if core.GetSimpleDoguName(installDogu) == core.GetSimpleDoguName(dogu) {
			return installDogu, nil
		}
	}

	return "", fmt.Errorf("dogu not found")
}

func (dv *doguValidator) getDogu(doguStr string) (*core.Dogu, error) {
	namespacedName, version, found := strings.Cut(doguStr, ":")
	var dogu *core.Dogu
	var err error
	if found {
		v, vErr := core.ParseVersion(version)
		if vErr != nil {
			return nil, fmt.Errorf("failed to parse dogu version %s: %w", version, err)
		}
		dogu, err = dv.registry.GetVersion(namespacedName, v.Raw)
	} else {
		dogu, err = dv.registry.Get(namespacedName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get dogu %s: %w", doguStr, err)
	}

	return dogu, nil
}
