package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

//IdentityProviderGitHub is a Github specific identity provider
type IdentityProviderGitHub struct {
	identityProviderCommon `yaml:",inline"`
	GitHub                 GitHub `yaml:"github"`
}

// GitHub provider specific data
type GitHub struct {
	HostName      string       `yaml:"hostname"`
	CA            CA           `yaml:"ca"`
	ClientID      string       `yaml:"clientID"`
	ClientSecret  ClientSecret `yaml:"clientSecret"`
	Organizations []string     `yaml:"organizations"`
	Teams         []string     `yaml:"teams"`
}

func buildGitHubIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderGitHub, secrets.Secret, error) {
	var idP IdentityProviderGitHub
	var secret *secrets.Secret

	var github configv1.GitHubIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &github)

	idP.Type = "GitHub"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.GitHub.HostName = github.Hostname
	idP.GitHub.CA.Name = github.CA
	idP.GitHub.ClientID = github.ClientID
	idP.GitHub.Organizations = github.Organizations
	idP.GitHub.Teams = github.Teams

	secretName := p.Name + "-secret"
	idP.GitHub.ClientSecret.Name = secretName

	encoded := base64.StdEncoding.EncodeToString([]byte(github.ClientSecret.Value))
	secret, err := secrets.GenSecret(secretName, encoded, "openshift-config", "literal")
	if err != nil {
		return idP, *secret, err
	}

	return idP, *secret, nil
}
