package validation

import (
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAdminValidator(t *testing.T) {
	// when
	validator := NewAdminValidator()

	// then
	require.NotNil(t, validator)
}

func Test_adminValidator_ValidateAdmin(t *testing.T) {
	tests := []struct {
		name             string
		user             context.User
		dsType           string
		containsErrorMsg string
		wantErr          assert.ErrorAssertionFunc
		wantErrMsg       assert.ComparisonAssertionFunc
	}{
		{"no admin group set", context.User{}, "", "no admin group set", assert.Error, assert.Contains},
		{"no admin mail set", context.User{AdminGroup: "group"}, dsTypeEmbedded, "no admin mail set", assert.Error, assert.Contains},
		{"invalid admin mail set", context.User{AdminGroup: "group", Mail: "t"}, dsTypeEmbedded, "invalid admin mail", assert.Error, assert.Contains},
		{"no admin username set", context.User{AdminGroup: "group", Mail: "test@test.de"}, dsTypeEmbedded, "no admin username set", assert.Error, assert.Contains},
		{"no admin password set", context.User{AdminGroup: "group", Mail: "test@test.de", Username: "name"}, dsTypeEmbedded, "no admin password set", assert.Error, assert.Contains},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			validator := &adminValidator{}

			// when
			result := validator.ValidateAdmin(tt.user, tt.dsType)

			// then
			tt.wantErr(t, result)
			tt.wantErrMsg(t, result.Error(), tt.containsErrorMsg)
		})
	}

	t.Run("successful active directory server validation", func(t *testing.T) {
		// given
		adminUser := context.User{AdminGroup: "group", Mail: "test@test.de", Username: "name", Password: "password"}
		validator := &adminValidator{}

		// when
		result := validator.ValidateAdmin(adminUser, dsTypeEmbedded)

		// then
		require.NoError(t, result)
	})
}
