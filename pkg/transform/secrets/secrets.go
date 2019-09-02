package secrets

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

// HTPasswdFileSecret is an htpasswd secret
type HTPasswdFileSecret struct {
	HTPasswd string `json:"htpasswd"`
}

// KeystoneFileSecret is a keystone secret
type KeystoneFileSecret struct {
	Keystone string `json:"keystone"`
}

// LiteralSecret is a literal secret
type LiteralSecret struct {
	ClientSecret string `json:"clientSecret"`
}

// BasicAuthFileSecret is a basic auth secret
type BasicAuthFileSecret struct {
	BasicAuth string `json:"basicAuth"`
}

// Secret contains a secret
type Secret struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Type       string      `json:"type"`
	Metadata   MetaData    `json:"metadata"`
	Data       interface{} `json:"data"`
}

// MetaData is the Metadata for a secret
type MetaData struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// SecretType is an enumerator for secret types
type SecretType int

const (
	// KeystoneSecretType - keystone type for Secret
	KeystoneSecretType = iota
	// HtpasswdSecretType - htpasswd type for Secret
	HtpasswdSecretType
	// LiteralSecretType - literal type for Secret
	LiteralSecretType
	// BasicAuthSecretType - basicauth type for Secret
	BasicAuthSecretType
)

var typeArray = []string{
	"KeystoneSecretType",
	"HtpasswdSecretType",
	"LiteralSecretType",
	"BasicAuthSecretType",
}

// APIVersion is the apiVersion string
var APIVersion = "v1"

const secretNameError = `Secret name is no valid, make sure it consists of lower case alphanumeric characters, ‘-’ or ‘.’,` +
	`and must start and end with an alphanumeric character (e.g. ‘example.com’, regex used for validation is ‘[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*’)`

// GenTLSSecret generates a TLS secret
func GenTLSSecret(name string, namespace string, cert []byte, key []byte) (*corev1.Secret, error) {
	nameErrors := validation.IsDNS1123Label(name)
	if nameErrors != nil {
		return nil, errors.New(secretNameError)
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.cert": cert,
			"tls.key":  key,
		},
	}

	return secret, nil
}

// GenSecret generates a secret
func GenSecret(name string, secretContent string, namespace string, secretType SecretType) (*Secret, error) {
	nameErrors := validation.IsDNS1123Label(name)

	if nameErrors != nil {
		return nil, errors.New(secretNameError)
	}

	data, err := buildData(secretType, secretContent)
	if err != nil {
		return nil, err
	}

	var secret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: MetaData{
			Name:      name,
			Namespace: namespace,
		},
	}
	return &secret, nil
}

func buildData(secretType SecretType, secretContent string) (interface{}, error) {
	var data interface{}

	switch secretType {
	case KeystoneSecretType:
		data = KeystoneFileSecret{Keystone: secretContent}
	case HtpasswdSecretType:
		data = HTPasswdFileSecret{HTPasswd: secretContent}
	case LiteralSecretType:
		data = LiteralSecret{ClientSecret: secretContent}
	case BasicAuthSecretType:
		data = BasicAuthFileSecret{BasicAuth: secretContent}
	default:
		return nil, errors.New("Not a valid secret type " + secretType.String())
	}

	return data, nil
}

// SecretType.String returns a string representation for SecretType enum
func (secType SecretType) String() string {
	if secType >= KeystoneSecretType && int(secType) < len(typeArray) {
		return typeArray[secType]
	}
	return "unknown"
}
