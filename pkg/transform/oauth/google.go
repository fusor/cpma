package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

//IdentityProviderGoogle is a Google specific identity provider
type IdentityProviderGoogle struct {
	identityProviderCommon `json:",inline"`
	Google                 Google `json:"google"`
}

// Google provider specific data
type Google struct {
	ClientID     string       `json:"clientID"`
	ClientSecret ClientSecret `json:"clientSecret"`
	HostedDomain string       `json:"hostedDomain,omitempty"`
}

func buildGoogleIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderGoogle, *secrets.Secret, error) {
	var (
		err    error
		idP    = &IdentityProviderGoogle{}
		secret *secrets.Secret
		google legacyconfigv1.GoogleIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &google); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to decode google, see error")
	}

	idP.Type = "Google"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.Google.ClientID = google.ClientID
	idP.Google.HostedDomain = google.HostedDomain

	secretName := p.Name + "-secret"
	idP.Google.ClientSecret.Name = secretName
	secretContent, err := io.FetchStringSource(google.ClientSecret)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to fetch client secret for google, see error")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(secretContent))
	if secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.LiteralSecretType); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to generate client secret for google, see error")
	}

	return idP, secret, nil
}

func validateGoogleProvider(serializer *json.Serializer, p IdentityProvider) error {
	var google legacyconfigv1.GoogleIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &google); err != nil {
		return errors.Wrap(err, "Failed to decode google, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if google.ClientSecret.KeyFile != "" {
		return errors.New("Usage of encrypted files as secret value is not supported")
	}

	if err := validateClientData(google.ClientID, google.ClientSecret); err != nil {
		return err
	}

	return nil
}
