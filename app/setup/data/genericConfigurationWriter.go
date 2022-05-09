package data

import (
	"reflect"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/pkg/errors"
)

// RegistryWriter is responsible to write entries into the registry.
type RegistryWriter interface {
	WriteConfigToRegistry(registryConfig context.CustomKeyValue) error
}

// GenericConfigurationWriter writes a configuration into the registry.
type GenericConfigurationWriter struct {
	Registry registry.Registry
}

// contextConfigWriter writes a configuration into a
// configurationContext (like _global or the context of a dogu)
type contextConfigWriter struct {
	ConfigCtx registry.ConfigurationContext
}

// NewGenericConfigurationWriter creates a new configuration writer.
func NewGenericConfigurationWriter(registry registry.Registry) *GenericConfigurationWriter {
	return &GenericConfigurationWriter{Registry: registry}
}

func (gcw *GenericConfigurationWriter) WriteConfigToRegistry(registryConfig context.CustomKeyValue) error {
	for fieldName, fieldMap := range registryConfig {
		err := gcw.writeEntriesForConfig(fieldMap, fieldName)
		if err != nil {
			return errors.Wrap(err, "failed to write entries in registry config")
		}
	}
	return nil
}

func (gcw *GenericConfigurationWriter) writeEntriesForConfig(entries map[string]interface{}, config string) error {
	var configCtx registry.ConfigurationContext

	logrus.Infof("write in %s configuration", config)
	if config == "_global" {
		configCtx = gcw.Registry.GlobalConfig()
	} else {
		configCtx = gcw.Registry.DoguConfig(config)
	}

	contextWriter := gcw.newContextConfigWriter(configCtx)
	for fieldName, fieldEntry := range entries {
		err := contextWriter.handleEntry(fieldName, fieldEntry)
		if err != nil {
			return errors.Wrapf(err, "failed to write %s in registry config", fieldName)
		}
	}

	return nil
}

// newContextConfigWriter returns a new object to write config to a specific context.
func (gcw *GenericConfigurationWriter) newContextConfigWriter(configCtx registry.ConfigurationContext) *contextConfigWriter {
	return &contextConfigWriter{configCtx}
}

// handleEntry writes values into the context configured in the writer.
func (contextWriter contextConfigWriter) handleEntry(field string, value interface{}) (err error) {
	switch value.(type) {
	case string:
		err = contextWriter.ConfigCtx.Set(field, value.(string))
		if err != nil {
			return errors.Wrapf(err, "could not set %s", value)
		}
	case map[string]string:
		for fieldName, fieldEntry := range value.(map[string]string) {
			err = contextWriter.handleEntry(field+"/"+fieldName, fieldEntry)
			if err != nil {
				break
			}
		}
	case map[string]interface{}:
		for fieldName, fieldEntry := range value.(map[string]interface{}) {
			err = contextWriter.handleEntry(field+"/"+fieldName, fieldEntry)
			if err != nil {
				break
			}
		}
	default:
		err = errors.Errorf("unexpected type %s for %s", reflect.TypeOf(value).String(), field)
	}
	return err
}
