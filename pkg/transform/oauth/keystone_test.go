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

func TestTransformMasterConfigKeystone(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/keystone/master_config.yaml")
	require.NoError(t, err)

	var expectedCrd configv1.OAuth
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Name = "cluster"
	expectedCrd.Namespace = oauth.OAuthNamespace

	var keystoneIDP = &configv1.IdentityProvider{}
	keystoneIDP.Type = "Keystone"
	keystoneIDP.Name = "my_keystone_provider"
	keystoneIDP.MappingMethod = "claim"
	keystoneIDP.Keystone = &configv1.KeystoneIdentityProvider{}
	keystoneIDP.Keystone.DomainName = "default"
	keystoneIDP.Keystone.URL = "http://fake.url:5000"
	keystoneIDP.Keystone.CA = configv1.ConfigMapNameReference{Name: "keystone-configmap"}
	keystoneIDP.Keystone.TLSClientCert.Name = "my_keystone_provider-client-cert-secret"
	keystoneIDP.Keystone.TLSClientKey.Name = "my_keystone_provider-client-key-secret"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, *keystoneIDP)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build keystone provider",
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

func TestKeystoneValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate keystone provider",
			requireError: false,
			inputFile:    "testdata/keystone/master_config.yaml",
		},
		{
			name:         "fail on invalid name in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid url in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-url-master-config.yaml",
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail on invalid domain name in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-domainname-master-config.yaml",
			expectedErr:  errors.New("Domain name can't be empty"),
		},
		{
			name:         "fail on invalid key file in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-keyfile-master-config.yaml",
			expectedErr:  errors.New("Key file can't be empty if cert file is specified"),
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
