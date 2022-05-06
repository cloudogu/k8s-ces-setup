package validation

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/remote"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	corev1 "k8s.io/api/core/v1"
)

type doguValidator struct {
	dogus    context.Dogus
	secret   corev1.Secret
	Registry remote.Registry
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
		return nil, fmt.Errorf("failed to create remote Registry: %w", err)
	}

	return &doguValidator{Registry: registry, dogus: dogus}, nil
}

func (dv *doguValidator) parseDoguStrToDoguList(dogus []string) ([]*core.Dogu, error) {
	var doguList = make([]*core.Dogu, len(dogus))
	for i, doguStr := range dogus {
		dogu, err := dv.getDoguFromVersionStr(doguStr)
		if err != nil {
			return nil, err
		}

		doguList[i] = dogu
	}

	return doguList, nil
}

func (dv *doguValidator) ValidateDogus() error {
	doguList, err := dv.parseDoguStrToDoguList(dv.dogus.Install)
	if err != nil {
		return err
	}

	for _, installDogu := range doguList {
		err = dv.validateDoguDependencies(doguList, installDogu.GetDependenciesOfType("dogu"))
		if err != nil {
			return fmt.Errorf("failed to validate dependencies for dogu %s: %w", installDogu.Name, err)
		}
	}

	return nil
}

func (dv *doguValidator) validateDoguDependencies(dogus []*core.Dogu, dependencies []core.Dependency) error {
	for _, dependency := range dependencies {
		depName := dependency.Name
		dependentDogu, err := dv.getDoguFromSelection(dogus, depName)
		if err != nil {
			return fmt.Errorf("dogu dependency %s ist not selected", depName)
		}

		if dependency.Version != "" {
			doguVersion, err := core.ParseVersion(dependentDogu.Version)
			if err != nil {
				return errors.Wrapf(err, "failed to parse version of dependency %s", dependency.Name)
			}

			comparator, err := core.ParseVersionComparator(dependency.Version)
			if err != nil {
				return errors.Wrapf(err, "failed to parse ParseVersionComparator of version %s for doguDependency %s", dependency.Version, dependency.Name)
			}

			allows, err := comparator.Allows(doguVersion)
			if err != nil {
				return errors.Wrapf(err, "An error occurred when comparing the versions")
			}
			if !allows {
				return errors.Errorf("%s parsed Version does not fulfill version requirement of %s dogu %s", dependentDogu.Version, dependency.Version, dependency.Name)
			}
		}
	}

	return nil
}

func (dv *doguValidator) getDoguFromSelection(dogus []*core.Dogu, doguName string) (*core.Dogu, error) {
	for _, installDogu := range dogus {
		if installDogu.GetSimpleName() == doguName {
			return installDogu, nil
		}
	}

	return nil, fmt.Errorf("dogu not found")
}

func (dv *doguValidator) getDoguFromVersionStr(doguStr string) (*core.Dogu, error) {
	namespacedName, version, found := strings.Cut(doguStr, ":")
	var dogu *core.Dogu
	var err error
	if found {
		v, vErr := core.ParseVersion(version)
		if vErr != nil {
			return nil, fmt.Errorf("failed to parse dogu version %s: %w", version, err)
		}
		dogu, err = dv.Registry.GetVersion(namespacedName, v.Raw)
	} else {
		dogu, err = dv.Registry.Get(namespacedName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get dogu %s: %w", doguStr, err)
	}

	return dogu, nil
}
