package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigLDAP(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/ldap/master_config.yaml")
	require.NoError(t, err)

	var expectedCrd configv1.OAuth
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Name = "cluster"
	expectedCrd.Namespace = oauth.OAuthNamespace

	var ldapIDP = &configv1.IdentityProvider{}
	ldapIDP.Name = "my_ldap_provider"
	ldapIDP.Type = "LDAP"
	ldapIDP.MappingMethod = "claim"
	ldapIDP.LDAP = &configv1.LDAPIdentityProvider{}
	ldapIDP.LDAP.Attributes.ID = []string{"dn"}
	ldapIDP.LDAP.Attributes.Email = []string{"mail"}
	ldapIDP.LDAP.Attributes.Name = []string{"cn"}
	ldapIDP.LDAP.Attributes.PreferredUsername = []string{"uid"}
	ldapIDP.LDAP.BindDN = "123"
	ldapIDP.LDAP.BindPassword.Name = "321"
	ldapIDP.LDAP.CA = configv1.ConfigMapNameReference{Name: "ldap-configmap"}
	ldapIDP.LDAP.Insecure = false
	ldapIDP.LDAP.URL = "ldap://ldap.example.com/ou=users,dc=acme,dc=com?uid"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, *ldapIDP)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build ldap provider",
			expectedCrd: &expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oauthResources, err := oauth.Translate(identityProviders, oauth.TokenConfig{})
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, oauthResources.OAuthCRD)
		})
	}
}

func TestLDAPValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate ldap provider",
			requireError: false,
			inputFile:    "testdata/ldap/master_config.yaml",
		},
		{
			name:         "fail on invalid name in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid ids in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-ids-master-config.yaml",
			expectedErr:  errors.New("ID can't be empty"),
		},
		{
			name:         "fail on invalid emails in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-emails-master-config.yaml",
			expectedErr:  errors.New("Email can't be empty"),
		},
		{
			name:         "fail on invalid names in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-names-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid preferred usernames in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-usernames-master-config.yaml",
			expectedErr:  errors.New("Preferred username can't be empty"),
		},
		{
			name:         "fail on invalid url in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-url-master-config.yaml",
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail if key file is present for bind password in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-bpass-master-config.yaml",
			expectedErr:  errors.New("Usage of encrypted files as bind password value is not supported"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			identityProvider, err := cpmatest.LoadIPTestData(tc.inputFile)
			require.NoError(t, err)

			err = oauth.Validate(identityProvider)

			if tc.requireError {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
