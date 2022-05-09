package validation_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/cesapp-lib/core"

	remoteMocks "github.com/cloudogu/cesapp-lib/remote/mocks"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/cloudogu/k8s-ces-setup/app/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDoguValidator(t *testing.T) {
	t.Parallel()

	t.Run("successfully create validator", func(t *testing.T) {
		// given
		registryMock := &remoteMocks.Registry{}

		// when
		validator := validation.NewDoguValidator(registryMock)

		// then
		assert.NotNil(t, validator)
	})
}

func Test_doguValidator_ValidateDogus(t *testing.T) {
	t.Parallel()

	t.Run("successful validation", func(t *testing.T) {
		// given
		doguA := "official/ldap:1.1.1-2"
		doguB := "official/cas:2.0.0-3"
		doguC := "official/redmine:3.1.2-1"
		doguD := "official/postfix"
		doguList := []string{doguA, doguB, doguC, doguD}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		ldapDogu := &core.Dogu{Name: "official/ldap", Version: "1.1.1-2"}
		casDogu := &core.Dogu{Name: "official/cas", Version: "2.0.0-3", Dependencies: []core.Dependency{{
			Type:    "dogu",
			Name:    "ldap",
			Version: ">1.0.0-0",
		}}}
		redmineDogu := &core.Dogu{Name: "official/redmine", Version: "3.1.2-1", Dependencies: []core.Dependency{
			{
				Type: "dogu",
				Name: "cas",
			}, {
				Type:    "dogu",
				Name:    "postfix",
				Version: "1.0.0-0",
			}}}
		postfixDogu := &core.Dogu{Name: "official/postfix", Version: "1.0.0-0"}
		mockRegistry.On("GetVersion", "official/ldap", "1.1.1-2").Return(ldapDogu, nil)
		mockRegistry.On("GetVersion", "official/cas", "2.0.0-3").Return(casDogu, nil)
		mockRegistry.On("GetVersion", "official/redmine", "3.1.2-1").Return(redmineDogu, nil)
		mockRegistry.On("Get", "official/postfix").Return(postfixDogu, nil)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("invalid postfix dependency", func(t *testing.T) {
		// given
		doguA := "official/ldap:1.1.1-2"
		doguB := "official/cas:2.0.0-3"
		doguC := "official/redmine:3.1.2-1"
		doguD := "official/postfix:0.0.1-1"
		doguList := []string{doguA, doguB, doguC, doguD}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		ldapDogu := &core.Dogu{Name: "official/ldap", Version: "1.1.1-2"}
		casDogu := &core.Dogu{Name: "official/cas", Version: "2.0.0-3", Dependencies: []core.Dependency{{
			Type:    "dogu",
			Name:    "ldap",
			Version: ">1.0.0-0",
		}}}
		redmineDogu := &core.Dogu{Name: "official/redmine", Version: "3.1.2-1", Dependencies: []core.Dependency{
			{
				Type: "dogu",
				Name: "cas",
			}, {
				Type:    "dogu",
				Name:    "postfix",
				Version: "1.0.0-0",
			}}}
		postfixDogu := &core.Dogu{Name: "official/postfix", Version: "0.0.1-1"}
		mockRegistry.On("GetVersion", "official/ldap", "1.1.1-2").Return(ldapDogu, nil)
		mockRegistry.On("GetVersion", "official/cas", "2.0.0-3").Return(casDogu, nil)
		mockRegistry.On("GetVersion", "official/redmine", "3.1.2-1").Return(redmineDogu, nil)
		mockRegistry.On("GetVersion", "official/postfix", "0.0.1-1").Return(postfixDogu, nil)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate dependencies for dogu official/redmine")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("failed to get dogu with version", func(t *testing.T) {
		// given
		doguA := "official/ldap:1.1.1-2"
		doguList := []string{doguA}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		mockRegistry.On("GetVersion", "official/ldap", "1.1.1-2").Return(nil, assert.AnError)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get dogu")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("failed to get dogu", func(t *testing.T) {
		// given
		doguA := "official/ldap"
		doguList := []string{doguA}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		mockRegistry.On("Get", "official/ldap").Return(nil, assert.AnError)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get dogu")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("failed to parse version", func(t *testing.T) {
		// given
		doguA := "official/ldap:1.1.asd.1-2"
		doguList := []string{doguA}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse dogu version")
	})

	t.Run("failed to get dogu from selection", func(t *testing.T) {
		// given
		doguB := "official/cas:2.0.0-3"
		casDogu := &core.Dogu{Name: "official/cas", Version: "2.0.0-3", Dependencies: []core.Dependency{{
			Type:    "dogu",
			Name:    "ldap",
			Version: ">1.0.0-0",
		}}}

		doguList := []string{doguB}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		mockRegistry.On("GetVersion", "official/cas", "2.0.0-3").Return(casDogu, nil)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dependency ldap ist not selected")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("failed to parse version for dependency", func(t *testing.T) {
		// given
		doguLdapID := "official/ldap:1.1.1-1"
		doguLdap := &core.Dogu{Name: "official/ldap", Version: "1.1.1-1", Dependencies: []core.Dependency{{
			Type:    "dogu",
			Name:    "cas",
			Version: "1.1.1-1",
		}}}

		doguCasID := "official/cas:1.1.1-1"
		doguCas := &core.Dogu{Name: "official/cas", Version: "1.1-1-1"}

		doguList := []string{doguLdapID, doguCasID}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		mockRegistry.On("GetVersion", "official/ldap", "1.1.1-1").Return(doguLdap, nil)
		mockRegistry.On("GetVersion", "official/cas", "1.1.1-1").Return(doguCas, nil)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse version")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("failed to parse version operator for dependency", func(t *testing.T) {
		// given
		doguLdapID := "official/ldap:1.1.1-1"
		doguLdap := &core.Dogu{Name: "official/ldap", Version: "1.1.1-1", Dependencies: []core.Dependency{{
			Type:    "dogu",
			Name:    "cas",
			Version: "<<<1.1.1-1",
		}}}

		doguCasID := "official/cas:1.1.1-1"
		doguCas := &core.Dogu{Name: "official/cas", Version: "1.1.1-1"}

		doguList := []string{doguLdapID, doguCasID}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		mockRegistry.On("GetVersion", "official/ldap", "1.1.1-1").Return(doguLdap, nil)
		mockRegistry.On("GetVersion", "official/cas", "1.1.1-1").Return(doguCas, nil)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse operator")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})

	t.Run("failed to parse version operator for dependency", func(t *testing.T) {
		// given
		doguLdapID := "official/ldap:1.1.1-1"
		doguLdap := &core.Dogu{Name: "official/ldap", Version: "1.1.1-1", Dependencies: []core.Dependency{{
			Type:    "dogu",
			Name:    "cas",
			Version: "=>1.1.1-1",
		}}}

		doguCasID := "official/cas:1.1.1-1"
		doguCas := &core.Dogu{Name: "official/cas", Version: "1.1.1-1"}

		doguList := []string{doguLdapID, doguCasID}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cockpit"}
		mockRegistry := &remoteMocks.Registry{}
		mockRegistry.On("Get", "cockpit").Return(nil, nil)
		mockRegistry.On("GetVersion", "official/ldap", "1.1.1-1").Return(doguLdap, nil)
		mockRegistry.On("GetVersion", "official/cas", "1.1.1-1").Return(doguCas, nil)
		doguValidator := validation.NewDoguValidator(mockRegistry)

		// when
		err := doguValidator.ValidateDogus(dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find suitable comperator for '=>'")
		mock.AssertExpectationsForObjects(t, mockRegistry)
	})
}
