package data

import (
	"context"
	"fmt"
	"reflect"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"

	k8sreg "github.com/cloudogu/k8s-registry-lib/registry"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

// RegistryWriter is responsible to write entries into the registry.
type RegistryWriter interface {
	WriteConfigToRegistry(registryConfig appcontext.CustomKeyValue) error
}

type ConfigurationRegistry interface {
	k8sreg.ConfigurationRegistry
}

type ConfigMapClient interface {
	k8sreg.ConfigMapClient
}

type DoguConfigurationRegistryProvider[T k8sreg.ConfigurationRegistry] func(ctx context.Context, name string) (T, error)

func (crp DoguConfigurationRegistryProvider[T]) GetConfig(ctx context.Context, name string) (T, error) {
	return crp(ctx, name)
}

func NewDoguConfigRegistryProvider[T k8sreg.ConfigurationRegistry](k8sClient ConfigMapClient) DoguConfigurationRegistryProvider[T] {
	return func(ctx context.Context, doguName string) (T, error) {
		reg, _ := k8sreg.NewDoguConfigRegistry(ctx, doguName, k8sClient)

		v, ok := interface{}(reg).(T)
		if !ok {
			panic("Used unsupported interface")
		}

		return v, nil
	}
}

// RegistryConfigurationWriter writes a configuration into the registry.
type RegistryConfigurationWriter struct {
	globalConfig       ConfigurationRegistry
	doguConfigProvider DoguConfigurationRegistryProvider[ConfigurationRegistry]
}

type InternalConfigRegistryProvider interface {
}

// NewRegistryConfigurationWriter creates a new configuration writer.
func NewRegistryConfigurationWriter(globalConfig ConfigurationRegistry, doguConfigProvider DoguConfigurationRegistryProvider[ConfigurationRegistry]) *RegistryConfigurationWriter {
	return &RegistryConfigurationWriter{
		globalConfig:       globalConfig,
		doguConfigProvider: doguConfigProvider,
	}
}

// WriteConfigToRegistry write the given registry config to the registry
func (gcw *RegistryConfigurationWriter) WriteConfigToRegistry(registryConfig appcontext.CustomKeyValue) error {
	for fieldName, fieldMap := range registryConfig {
		err := gcw.writeEntriesForConfig(fieldMap, fieldName)
		if err != nil {
			return errors.Wrap(err, "failed to write entries in registry config")
		}
	}
	return nil
}

func (gcw *RegistryConfigurationWriter) writeEntriesForConfig(entries map[string]interface{}, config string) error {
	var configCtx ConfigurationRegistry

	logrus.Infof("write in %s configuration", config)
	if config == "_global" {
		configCtx = gcw.globalConfig
	} else {
		var err error
		configCtx, err = gcw.doguConfigProvider.GetConfig(context.Background(), config)
		if err != nil {
			return fmt.Errorf("failed to create dogu config: %w", err)
		}
	}

	contextWriter := gcw.newConfigWriter(func(field string, value string) error {
		return configCtx.Set(context.Background(), field, value)
	})
	for fieldName, fieldEntry := range entries {
		err := contextWriter.handleEntry(fieldName, fieldEntry)
		if err != nil {
			return errors.Wrapf(err, "failed to write %s in registry config", fieldName)
		}
	}

	return nil
}

// newConfigWriter returns a new object to write config with a given function.
func (gcw *RegistryConfigurationWriter) newConfigWriter(writer configurationContextWriter) *configWriter {
	return &configWriter{write: writer, delimiter: "/"}
}

type configurationContextWriter = func(field string, value string) error

// configWriter writes a configuration with a given function
type configWriter struct {
	write     configurationContextWriter
	delimiter string
}

// handleEntry writes values into the appcontext implemented in the write process.
func (contextWriter configWriter) handleEntry(field string, value interface{}) (err error) {
	switch value.(type) {
	case string:
		// nolint
		err = contextWriter.write(field, value.(string))
		if err != nil {
			return errors.Wrapf(err, "could not set %s", value)
		}
	case map[string]string:
		// nolint
		for fieldName, fieldEntry := range value.(map[string]string) {
			err = contextWriter.handleEntry(field+contextWriter.delimiter+fieldName, fieldEntry)
			if err != nil {
				break
			}
		}
	case map[string]interface{}:
		// nolint
		for fieldName, fieldEntry := range value.(map[string]interface{}) {
			err = contextWriter.handleEntry(field+contextWriter.delimiter+fieldName, fieldEntry)
			if err != nil {
				break
			}
		}
	default:
		err = errors.Errorf("unexpected type %s for %s", reflect.TypeOf(value).String(), field)
	}
	return err
}
