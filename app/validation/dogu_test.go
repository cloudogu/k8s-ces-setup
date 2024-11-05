package validation

import (
	ctx "context"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/cloudogu/cesapp-lib/core"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDoguValidator(t *testing.T) {
	t.Parallel()

	t.Run("successfully create validator", func(t *testing.T) {
		// given
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)

		// when
		validator := NewDoguValidator(remoteDoguRepo)

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
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
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
		ldapVersion, _ := core.ParseVersion("1.1.1-2")
		casVersion, _ := core.ParseVersion("2.0.0-3")
		redmineVersion, _ := core.ParseVersion("3.1.2-1")
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "ldap",
			},
			Version: ldapVersion,
		}).Return(ldapDogu, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "cas",
			},
			Version: casVersion,
		}).Return(casDogu, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "redmine",
			},
			Version: redmineVersion,
		}).Return(redmineDogu, nil)
		remoteDoguRepo.EXPECT().GetLatest(mock.Anything, cescommons.QualifiedDoguName{
			Namespace:  "official",
			SimpleName: "postfix",
		}).Return(postfixDogu, nil)

		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
	})

	t.Run("invalid postfix dependency", func(t *testing.T) {
		// given
		doguA := "official/ldap:1.1.1-2"
		doguB := "official/cas:2.0.0-3"
		doguC := "official/redmine:3.1.2-1"
		doguD := "official/postfix:0.0.1-1"
		doguList := []string{doguA, doguB, doguC, doguD}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
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

		ldapVersion, _ := core.ParseVersion("1.1.1-2")
		casVersion, _ := core.ParseVersion("2.0.0-3")
		redmineVersion, _ := core.ParseVersion("3.1.2-1")
		postfixVersion, _ := core.ParseVersion("0.0.1-1")
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "ldap",
			},
			Version: ldapVersion,
		}).Return(ldapDogu, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "cas",
			},
			Version: casVersion,
		}).Return(casDogu, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "redmine",
			},
			Version: redmineVersion,
		}).Return(redmineDogu, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "postfix",
			},
			Version: postfixVersion,
		}).Return(postfixDogu, nil)

		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to validate dependencies for dogu official/redmine")
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
	})

	t.Run("failed to get dogu with version", func(t *testing.T) {
		// given
		doguA := "official/ldap:1.1.1-2"
		doguList := []string{doguA}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)

		ldapVersion, _ := core.ParseVersion("1.1.1-2")
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "ldap",
			},
			Version: ldapVersion,
		}).Return(&core.Dogu{}, assert.AnError)
		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get version of dogu [{ldap official}] [1.1.1-2]")
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
	})

	t.Run("failed to get dogu", func(t *testing.T) {
		// given
		doguA := "official/ldap"
		doguList := []string{doguA}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		remoteDoguRepo.EXPECT().GetLatest(mock.Anything, cescommons.QualifiedDoguName{
			Namespace:  "official",
			SimpleName: "ldap",
		}).Return(&core.Dogu{}, assert.AnError)

		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get latest version of dogu [{ldap official}]")
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
	})

	t.Run("failed to parse version", func(t *testing.T) {
		// given
		doguA := "official/ldap:1.1.asd.1-2"
		doguList := []string{doguA}
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

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
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		casVersion, _ := core.ParseVersion("2.0.0-3")
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "cas",
			},
			Version: casVersion,
		}).Return(casDogu, nil)

		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dependency ldap ist not selected")
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
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
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		ldapVersion, _ := core.ParseVersion("1.1.1-1")
		casVersion, _ := core.ParseVersion("1.1.1-1")
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "ldap",
			},
			Version: ldapVersion,
		}).Return(doguLdap, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "cas",
			},
			Version: casVersion,
		}).Return(doguCas, nil)

		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse version")
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
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
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		ldapVersion, _ := core.ParseVersion("1.1.1-1")
		casVersion, _ := core.ParseVersion("1.1.1-1")
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "ldap",
			},
			Version: ldapVersion,
		}).Return(doguLdap, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "cas",
			},
			Version: casVersion,
		}).Return(doguCas, nil)

		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse operator")
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
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
		dogus := context.Dogus{Install: doguList, DefaultDogu: "cas"}
		remoteDoguRepo := newMockRemoteDoguDescriptorRepository(t)
		ldapVersion, _ := core.ParseVersion("1.1.1-1")
		casVersion, _ := core.ParseVersion("1.1.1-1")
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "ldap",
			},
			Version: ldapVersion,
		}).Return(doguLdap, nil)
		remoteDoguRepo.EXPECT().Get(mock.Anything, cescommons.QualifiedDoguVersion{
			Name: cescommons.QualifiedDoguName{
				Namespace:  "official",
				SimpleName: "cas",
			},
			Version: casVersion,
		}).Return(doguCas, nil)
		doguValidator := NewDoguValidator(remoteDoguRepo)

		// when
		err := doguValidator.ValidateDogus(ctx.TODO(), dogus)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find suitable comperator for '=>'")
		mock.AssertExpectationsForObjects(t, remoteDoguRepo)
	})
}
