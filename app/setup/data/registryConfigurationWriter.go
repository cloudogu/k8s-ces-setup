package data

import (
	"reflect"

	"github.com/cloudogu/cesapp-lib/registry"
	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

// RegistryWriter is responsible to write entries into the registry.
type RegistryWriter interface {
	WriteConfigToRegistry(registryConfig appcontext.CustomKeyValue) error
}

// RegistryConfigurationWriter writes a configuration into the registry.
type RegistryConfigurationWriter struct {
	Registry registry.Registry
}

// NewRegistryConfigurationWriter creates a new configuration writer.
func NewRegistryConfigurationWriter(registry registry.Registry) *RegistryConfigurationWriter {
	return &RegistryConfigurationWriter{Registry: registry}
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
	var configCtx registry.ConfigurationContext

	logrus.Infof("write in %s configuration", config)
	if config == "_global" {
		configCtx = gcw.Registry.GlobalConfig()
	} else {
		configCtx = gcw.Registry.DoguConfig(config)
	}

	contextWriter := gcw.newConfigWriter(func(field string, value string) error {
		return configCtx.Set(field, value)
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

// handleEntry writes values into the appcontext implemented in the write.
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
