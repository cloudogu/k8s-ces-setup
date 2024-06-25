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

// RegistryConfigurationWriter writes a configuration into the registry.
type RegistryConfigurationWriter struct {
	globalConfig       k8sreg.ConfigurationWriter
	doguConfigProvider k8sreg.DoguConfigRegistryProvider
}

// NewRegistryConfigurationWriter creates a new configuration writer.
func NewRegistryConfigurationWriter(globalConfig k8sreg.ConfigurationWriter, doguConfigProvider k8sreg.DoguConfigRegistryProvider) *RegistryConfigurationWriter {
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
	var configCtx k8sreg.ConfigurationWriter

	logrus.Infof("write in %s configuration", config)
	if config == "_global" {
		configCtx = gcw.globalConfig
	} else {
		var err error
		configCtx, err = gcw.doguConfigProvider.GetDoguConfig(context.Background(), config)
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
