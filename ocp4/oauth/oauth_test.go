package oauth

import (
	"testing"

	"github.com/fusor/cpma/ocp3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslateMasterConfig(t *testing.T) {
	masterConfig := ocp3.ParseMaster("../../test/oauth/bulk-test-master-config.yaml")

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, resCrd.Spec.IdentityProviders[0].(identityProviderBasicAuth).Type, "BasicAuth")
	assert.Equal(t, resCrd.Spec.IdentityProviders[1].(identityProviderGitHub).Type, "GitHub")
	assert.Equal(t, resCrd.Spec.IdentityProviders[2].(identityProviderGitLab).Type, "GitLab")
	assert.Equal(t, resCrd.Spec.IdentityProviders[3].(identityProviderGoogle).Type, "Google")
	assert.Equal(t, resCrd.Spec.IdentityProviders[4].(identityProviderHTPasswd).Type, "HTPasswd")
	assert.Equal(t, resCrd.Spec.IdentityProviders[5].(identityProviderKeystone).Type, "Keystone")
	assert.Equal(t, resCrd.Spec.IdentityProviders[6].(identityProviderLDAP).Type, "LDAP")
	assert.Equal(t, resCrd.Spec.IdentityProviders[7].(identityProviderRequestHeader).Type, "RequestHeader")
	assert.Equal(t, resCrd.Spec.IdentityProviders[8].(identityProviderOpenID).Type, "OpenID")
}

func TestGenYAML(t *testing.T) {
	masterConfig := ocp3.ParseMaster("../../test/oauth/htpasswd-test-master-config.yaml")
	crd, _, err := Translate(masterConfig.OAuthConfig)

	CRD := crd.GenYAML()
	expectedYaml := `apiVersion: config.openshift.io/v1
kind: OAuth
metadata:
  name: cluster
  namespace: openshift-config
spec:
  identityProviders:
  - name: htpasswd_auth
    challenge: true
    login: true
    mappingMethod: claim
    type: HTPasswd
    htpasswd:
      fileData:
        name: htpasswd_auth-secret
`
	require.NoError(t, err)
	assert.Equal(t, expectedYaml, string(CRD))
}
