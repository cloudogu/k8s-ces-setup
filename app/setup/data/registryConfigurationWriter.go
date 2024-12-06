package data

import (
	"context"
	"fmt"
	k8sconf "github.com/cloudogu/k8s-registry-lib/config"
	k8serror "github.com/cloudogu/k8s-registry-lib/errors"
	k8sreg "github.com/cloudogu/k8s-registry-lib/repository"
	"reflect"

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
	globalConfig *k8sreg.GlobalConfigRepository
	doguConfig   *k8sreg.DoguConfigRepository
}

// NewRegistryConfigurationWriter creates a new configuration writer.
func NewRegistryConfigurationWriter(globalConfig *k8sreg.GlobalConfigRepository, doguConfigProvider *k8sreg.DoguConfigRepository) *RegistryConfigurationWriter {
	return &RegistryConfigurationWriter{
		globalConfig: globalConfig,
		doguConfig:   doguConfigProvider,
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
	var contextWriter *configWriter
	ctx := context.TODO()

	logrus.Infof("write in %s configuration", config)
	if config == "_global" {
		contextWriter = gcw.newConfigWriter(func(field string, value string) error {
			c, err := gcw.globalConfig.Get(ctx)
			if err != nil && k8serror.IsNotFoundError(err) {
				c, err = gcw.globalConfig.Create(ctx, k8sconf.CreateGlobalConfig(make(k8sconf.Entries)))
				if err != nil {
					return fmt.Errorf("failed to create global config: %w", err)
				}
			} else if err != nil {
				return fmt.Errorf("failed to get global config: %w", err)
			}

			c.Config, err = c.Set(k8sconf.Key(field), k8sconf.Value(value))
			if err != nil {
				return fmt.Errorf("failed to set key '%s' in global config to '%s': %w", field, value, err)
			}

			c, err = gcw.globalConfig.SaveOrMerge(ctx, c)
			return err
		})
	} else {
		contextWriter = gcw.newConfigWriter(func(field string, value string) error {
			c, err := gcw.doguConfig.Get(ctx, k8sconf.SimpleDoguName(config))
			if err != nil && k8serror.IsNotFoundError(err) {
				c, err = gcw.doguConfig.Create(ctx, k8sconf.CreateDoguConfig(k8sconf.SimpleDoguName(config), make(k8sconf.Entries)))
				if err != nil {
					return fmt.Errorf("failed to create dogu config for '%s': %w", config, err)
				}
			} else if err != nil {
				return fmt.Errorf("failed to get dogu config for '%s': %w", config, err)
			}

			c.Config, err = c.Set(k8sconf.Key(field), k8sconf.Value(value))
			if err != nil {
				return fmt.Errorf("failed to set key '%s' in dogu config for '%s' to '%s': %w", config, field, value, err)
			}

			_, err = gcw.doguConfig.SaveOrMerge(ctx, c)
			return err
		})
	}

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
	switch castedConfigValue := value.(type) {
	case string:
		err = contextWriter.write(field, castedConfigValue)
		if err != nil {
			return errors.Wrapf(err, "could not set %s", value)
		}
	case map[string]string:
		for fieldName, fieldEntry := range castedConfigValue {
			err = contextWriter.handleEntry(field+contextWriter.delimiter+fieldName, fieldEntry)
			if err != nil {
				break
			}
		}
	case map[string]interface{}:
		for fieldName, fieldEntry := range castedConfigValue {
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
