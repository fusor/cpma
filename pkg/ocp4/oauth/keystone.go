package oauth

import (
	"encoding/base64"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type IdentityProviderKeystone struct {
	identityProviderCommon `yaml:",inline"`
	Keystone               struct {
		DomainName string `yaml:"domainName"`
		URL        string `yaml:"url"`
		CA         struct {
			Name string `yaml:"name"`
		} `yaml:"ca"`
		TLSClientCert struct {
			Name string `yaml:"name"`
		} `yaml:"tlsClientCert"`
		TLSClientKey struct {
			Name string `yaml:"name"`
		} `yaml:"tlsClientKey"`
	} `yaml:"keystone"`
}

func buildKeystoneIP(serializer *json.Serializer, p ocp3.IdentityProvider) (IdentityProviderKeystone, secrets.Secret, secrets.Secret) {
	var idP IdentityProviderKeystone
	var keystone configv1.KeystonePasswordIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &keystone)

	idP.Type = "Keystone"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.Keystone.DomainName = keystone.DomainName
	idP.Keystone.URL = keystone.URL
	idP.Keystone.CA.Name = keystone.CA

	certSecretName := p.Name + "-client-cert-secret"
	idP.Keystone.TLSClientCert.Name = certSecretName

	// TODO: Fetch cert and key
	certFile := "456This is pretend content"
	encoded := base64.StdEncoding.EncodeToString([]byte(certFile))
	certSecret := secrets.GenSecret(certSecretName, encoded, "openshift-config", "keystone")

	keySecretName := p.Name + "-client-key-secret"
	idP.Keystone.TLSClientKey.Name = keySecretName

	keyFile := "123This is pretend content"
	encoded = base64.StdEncoding.EncodeToString([]byte(keyFile))
	keySecret := secrets.GenSecret(keySecretName, encoded, "openshift-config", "keystone")

	return idP, *certSecret, *keySecret
}
