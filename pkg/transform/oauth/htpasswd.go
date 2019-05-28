package oauth

import (
	"encoding/base64"
	"errors"

	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// IdentityProviderHTPasswd is a htpasswd specific identity provider
type IdentityProviderHTPasswd struct {
	identityProviderCommon `yaml:",inline"`
	HTPasswd               `yaml:"htpasswd"`
}

// HTPasswd contains htpasswd FileData
type HTPasswd struct {
	FileData FileData `yaml:"fileData"`
}

// FileData from htpasswd file
type FileData struct {
	Name string `yaml:"name"`
}

func buildHTPasswdIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderHTPasswd, *secrets.Secret, error) {
	var (
		err      error
		idP      = &IdentityProviderHTPasswd{}
		secret   *secrets.Secret
		htpasswd configv1.HTPasswdPasswordIdentityProvider
	)

	_, _, err = serializer.Decode(p.Provider.Raw, nil, &htpasswd)
	if err != nil {
		return nil, nil, err
	}

	idP.Name = p.Name
	idP.Type = "HTPasswd"
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.HTPasswd.FileData.Name = htpasswd.File

	secretName := p.Name + "-secret"
	idP.HTPasswd.FileData.Name = secretName

	encoded := base64.StdEncoding.EncodeToString(p.HTFileData)

	secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.HtpasswdSecretType)
	if err != nil {
		return nil, nil, err
	}

	return idP, secret, nil
}

func validateHTPasswdProvider(serializer *json.Serializer, p IdentityProvider) error {
	var htpasswd configv1.HTPasswdPasswordIdentityProvider

	_, _, err := serializer.Decode(p.Provider.Raw, nil, &htpasswd)
	if err != nil {
		return err
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if htpasswd.File == "" {
		return errors.New("File can't be empty")
	}

	return nil
}
