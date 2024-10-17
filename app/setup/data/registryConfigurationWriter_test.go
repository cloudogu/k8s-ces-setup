package data_test

import (
	k8sconf "github.com/cloudogu/k8s-registry-lib/config"
	k8sreg "github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"

	"github.com/cloudogu/k8s-ces-setup/app/setup/data"
	"github.com/stretchr/testify/require"
)

func TestNewGenericConfigurationWriter(t *testing.T) {
	t.Run("create new generic configuration Writer", func(t *testing.T) {
		cm := &data.MockConfigMapClient{}
		// given
		globalReg := k8sreg.NewGlobalConfigRepository(cm)
		doguReg := k8sreg.NewDoguConfigRepository(cm)

		// when
		writer := data.NewRegistryConfigurationWriter(
			globalReg,
			doguReg,
		)

		// then
		require.NotNil(t, writer)
	})
}

func TestGenericConfigurationWriter_WriteConfigToRegistry(t *testing.T) {
	t.Run("failed to write to config", func(t *testing.T) {
		cm := &data.MockConfigMapClient{}

		// given
		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"test3": "myTestKey3",
			},
		}

		globalReg := k8sreg.NewGlobalConfigRepository(cm)
		doguReg := k8sreg.NewDoguConfigRepository(cm)

		cm.EXPECT().List(mock.Anything, mock.Anything).Return(nil, assert.AnError)

		writer := data.NewRegistryConfigurationWriter(globalReg, doguReg)

		// when
		err := writer.WriteConfigToRegistry(registryConfig)

		// then
		require.Contains(t, err.Error(), assert.AnError.Error())
	})

	t.Run("set all keys correctly", func(t *testing.T) {
		cm := &data.MockConfigMapClient{}

		// given
		registryConfig := appcontext.CustomKeyValue{
			"_global": map[string]interface{}{
				"test": map[string]string{
					"t1": "myTestt1",
					"t2": "myTestt2",
				},
				"test3": "myTestKey3",
			},
			"cas": map[string]interface{}{
				"ldap": map[string]string{
					"attribute_fullname": "myAttributeFullName",
				},
			},
			"ldap-mapper": map[string]interface{}{
				"backend": map[string]string{
					"type": "external",
				},
				"mapping": map[string]interface{}{
					"group": map[string]string{
						"base_dn": "myGroupBaseDN",
						"member":  "myGroupAttributeMember",
					},
					"user": map[string]string{
						"base_dn":   "myBaseDN",
						"full_name": "myAttributeFullName",
						"surname":   "myAttributeSurname",
					},
				},
			},
		}

		globalReg := k8sreg.NewGlobalConfigRepository(cm)
		doguReg := k8sreg.NewDoguConfigRepository(cm)

		emptyCm := &v1.ConfigMap{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data: map[string]string{
				"config.yaml": "{}",
			},
			BinaryData: nil,
		}
		globalConfig := v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: "global-config",
			},
			Immutable: nil,
			Data: map[string]string{
				"config.yaml": "{}",
			},
			BinaryData: nil,
		}
		// List is called as SingletonList for watches. Get is called at the saveOrMerge step
		cm.EXPECT().List(mock.Anything, metav1.SingleObject(metav1.ObjectMeta{Name: "global-config"})).Return(&v1.ConfigMapList{Items: []v1.ConfigMap{globalConfig}}, nil)
		cm.EXPECT().Get(mock.Anything, "global-config", mock.Anything).Return(&globalConfig, nil)

		casConfig := v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: "ldap-mapper-config",
			},
			Immutable: nil,
			Data: map[string]string{
				"config.yaml": "{}",
			},
			BinaryData: nil,
		}
		cm.EXPECT().List(mock.Anything, metav1.SingleObject(metav1.ObjectMeta{Name: "cas-config"})).Return(&v1.ConfigMapList{Items: []v1.ConfigMap{casConfig}}, nil)
		cm.EXPECT().Get(mock.Anything, "cas-config", mock.Anything).Return(&casConfig, nil)

		ldapMapperConfig := v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: "ldap-mapper-config",
			},
			Immutable: nil,
			Data: map[string]string{
				"config.yaml": "{}",
			},
			BinaryData: nil,
		}
		cm.EXPECT().List(mock.Anything, metav1.SingleObject(metav1.ObjectMeta{Name: "ldap-mapper-config"})).Return(&v1.ConfigMapList{Items: []v1.ConfigMap{ldapMapperConfig}}, nil)
		cm.EXPECT().Get(mock.Anything, "ldap-mapper-config", mock.Anything).Return(&ldapMapperConfig, nil)

		cm.EXPECT().Update(mock.Anything, mock.MatchedBy(func(inputCm *v1.ConfigMap) bool {
			ldapMapperReader, err := (&k8sconf.YamlConverter{}).Read(strings.NewReader("backend:\n    type: external\nldap:\n    attribute_fullname: myAttributeFullName\nmapping:\n    group:\n        base_dn: myGroupBaseDN\n        member: myGroupAttributeMember\n    user:\n        base_dn: myBaseDN\n        full_name: myAttributeFullName\n        surname: myAttributeSurname\ntest:\n    t1: myTestt1\n    t2: myTestt2\ntest3: myTestKey3\n"))
			require.NoError(t, err)
			casReader, err := (&k8sconf.YamlConverter{}).Read(strings.NewReader("ldap:\n    attribute_fullname: myAttributeFullName\ntest:\n    t1: myTestt1\n    t2: myTestt2\ntest3: myTestKey3\n"))
			require.NoError(t, err)
			globalReader, err := (&k8sconf.YamlConverter{}).Read(strings.NewReader("test:\n    t1: myTestt1\n    t2: myTestt2\ntest3: myTestKey3\n"))
			require.NoError(t, err)

			result, err := (&k8sconf.YamlConverter{}).Read(strings.NewReader(inputCm.Data["config.yaml"]))
			require.NoError(t, err)

			for key, value := range result {
				if inputCm.ObjectMeta.Name == "ldap-mapper-config" {
					if !(ldapMapperReader[key] == value) {
						return false
					}
				} else if inputCm.ObjectMeta.Name == "cas-config" {
					if !(casReader[key] == value) {
						return false
					}
				} else if inputCm.ObjectMeta.Name == "global-config" {
					if !(globalReader[key] == value) {
						return false
					}
				} else {
					return false
				}
			}

			return true
		}), mock.Anything).Return(emptyCm, nil)

		writer := data.NewRegistryConfigurationWriter(globalReg, doguReg)

		// when
		err := writer.WriteConfigToRegistry(registryConfig)

		// then
		require.NoError(t, err)
	})
}
