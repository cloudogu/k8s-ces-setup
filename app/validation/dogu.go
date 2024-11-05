package validation

import (
	ctx "context"
	"fmt"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-ces-setup/app/retry"
	"strings"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-ces-setup/app/context"
)

var maxTries = 20

type doguValidator struct {
	secret     corev1.Secret
	Repository cescommons.RemoteDoguDescriptorRepository
}

// NewDoguValidator creates a new validator for the dogu region of the setup configuration.
func NewDoguValidator(repository cescommons.RemoteDoguDescriptorRepository) *doguValidator {
	return &doguValidator{Repository: repository}
}

// ValidateDogus check whether the configured dogu has no invalid or unmet dependencies.
func (dv *doguValidator) ValidateDogus(ctx ctx.Context, dogus context.Dogus) error {
	doguList, err := dv.parseDoguStrToDoguList(ctx, dogus.Install)
	if err != nil {
		return err
	}

	isDeafultDoguValid := false
	for _, dogu := range dogus.Install {
		if strings.Contains(dogu, dogus.DefaultDogu) {
			isDeafultDoguValid = true
		}
	}

	if !isDeafultDoguValid {
		return fmt.Errorf("invalid value for default dogu [%s]", dogus.DefaultDogu)
	}

	for _, installDogu := range doguList {
		err = dv.validateDoguDependencies(doguList, installDogu.GetDependenciesOfType("dogu"))
		if err != nil {
			return fmt.Errorf("failed to validate dependencies for dogu %s: %w", installDogu.Name, err)
		}
	}

	return nil
}

func (dv *doguValidator) parseDoguStrToDoguList(ctx ctx.Context, dogus []string) ([]*core.Dogu, error) {
	var doguList = make([]*core.Dogu, len(dogus))
	for i, doguStr := range dogus {
		dogu, err := dv.getDoguFromVersionStr(ctx, doguStr)
		if err != nil {
			return nil, err
		}

		doguList[i] = dogu
	}

	return doguList, nil
}

func (dv *doguValidator) validateDoguDependencies(dogus []*core.Dogu, dependencies []core.Dependency) error {
	for _, dependency := range dependencies {
		depName := dependency.Name
		if depName == "nginx" || depName == "registrator" {
			continue
		}
		dependentDogu, err := dv.getDoguFromSelection(dogus, depName)
		if err != nil {
			return fmt.Errorf("dogu dependency %s ist not selected", depName)
		}

		if dependency.Version == "" {
			continue
		}

		allows, err := isDependencyVersionAllowed(dependentDogu, dependency)
		if err != nil {
			return err
		}
		if !allows {
			return errors.Errorf("%s parsed Version does not fulfill version requirement of %s dogu %s", dependentDogu.Version, dependency.Version, dependency.Name)
		}
	}

	return nil
}

func isDependencyVersionAllowed(dependentDogu *core.Dogu, dependency core.Dependency) (bool, error) {
	doguVersion, err := core.ParseVersion(dependentDogu.Version)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse version of dependency %s", dependency.Name)
	}

	comparator, err := core.ParseVersionComparator(dependency.Version)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse ParseVersionComparator of version %s for doguDependency %s", dependency.Version, dependency.Name)
	}

	allows, err := comparator.Allows(doguVersion)
	if err != nil {
		return false, errors.Wrapf(err, "An error occurred when comparing the versions")
	}

	return allows, nil
}

func (dv *doguValidator) getDoguFromSelection(dogus []*core.Dogu, doguName string) (*core.Dogu, error) {
	for _, installDogu := range dogus {
		if installDogu.GetSimpleName() == doguName {
			return installDogu, nil
		}
	}

	return nil, fmt.Errorf("dogu not found")
}

func (dv *doguValidator) getDoguFromVersionStr(ctx ctx.Context, doguStr string) (*core.Dogu, error) {
	namespacedName, version, found := strings.Cut(doguStr, ":")
	namespace, name, _ := strings.Cut(namespacedName, "/")
	var dogu *core.Dogu
	var err error

	qualifiedDoguName := cescommons.QualifiedDoguName{
		SimpleName: cescommons.SimpleDoguName(name),
		Namespace:  cescommons.DoguNamespace(namespace),
	}
	if found {
		v, vErr := core.ParseVersion(version)
		if vErr != nil {
			return nil, fmt.Errorf("failed to parse dogu version %s: %w", version, err)
		}
		qualifiedDoguVersion := cescommons.QualifiedDoguVersion{
			Name:    qualifiedDoguName,
			Version: v,
		}
		err := retry.OnError(maxTries, isConnectionError, func() error {
			var err error
			dogu, err = dv.Repository.Get(ctx, qualifiedDoguVersion)
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get version of dogu [%s] [%s]: %w", qualifiedDoguName, v.Raw, err)
		}
	} else {
		err := retry.OnError(maxTries, isConnectionError, func() error {
			var err error
			dogu, err = dv.Repository.GetLatest(ctx, qualifiedDoguName)
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version of dogu [%s]: %w", qualifiedDoguName, err)
		}
	}

	return dogu, nil
}
func isConnectionError(err error) bool {
	return strings.Contains(err.Error(), cescommons.ConnectionError.Error())
}
