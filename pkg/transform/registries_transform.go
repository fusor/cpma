package transform

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// RegistriesComponentName is the registry component string
const RegistriesComponentName = "Registries"

// RegistriesExtraction holds registry information extracted from an OCP3 cluster
type RegistriesExtraction struct {
	Registries map[string]Registries
}

// Registries holds a list of Registries
type Registries struct {
	List []string `toml:"registries"`
}

// ImageCR is an Image Cluster Resource
type ImageCR struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Metadata   Metadata  `json:"metadata"`
	Spec       ImageSpec `json:"spec"`
}

// Metadata is the Metadata for an Image Cluster Resource
type Metadata struct {
	Name        string
	Annotations map[string]string `json:"annotations"`
}

// ImageSpec is a Spec for an ImageCR
type ImageSpec struct {
	RegistrySources RegistrySources `json:"registrySources"`
}

// RegistrySources holds lists of blocked and insecure registries from an OCP3 cluster
type RegistrySources struct {
	BlockedRegistries  []string `json:"blockedRegistries,omitempty"`
	InsecureRegistries []string `json:"insecureRegistries,omitempty"`
}

// RegistriesTransform is a registry specific transform
type RegistriesTransform struct {
}

// Transform contains registry configuration collected from an OCP3 into a useful output
func (e RegistriesExtraction) Transform() ([]Output, error) {
	logrus.Info("RegistriesTransform::Extraction")
	manifests, err := e.buildManifestOutput()
	if err != nil {
		return nil, err
	}
	reports, err := e.buildReportOutput()
	if err != nil {
		return nil, err
	}
	outputs := []Output{manifests, reports}
	return outputs, nil
}

func (e RegistriesExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	const (
		apiVersion = "config.openshift.io/v1"
		kind       = "Image"
		name       = "cluster"
		annokey    = "release.openshift.io/create-only"
		annoval    = "true"
	)

	var imageCR ImageCR
	imageCR.APIVersion = apiVersion
	imageCR.Kind = kind
	imageCR.Metadata.Name = name
	imageCR.Metadata.Annotations = make(map[string]string)
	imageCR.Metadata.Annotations[annokey] = annoval
	imageCR.Spec.RegistrySources.BlockedRegistries = e.Registries["block"].List
	imageCR.Spec.RegistrySources.InsecureRegistries = e.Registries["insecure"].List

	imageCRYAML, err := yaml.Marshal(&imageCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-registries.yaml", CRD: imageCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e RegistriesExtraction) buildReportOutput() (Output, error) {
	reportOutput := ReportOutput{
		Component: RegistriesComponentName,
	}

	for _, registry := range e.Registries["block"].List {
		reportOutput.Reports = append(reportOutput.Reports,
			Report{
				Name:       registry,
				Kind:       "Blocked",
				Supported:  true,
				Confidence: HighConfidence,
			})
	}

	for _, registry := range e.Registries["insecure"].List {
		reportOutput.Reports = append(reportOutput.Reports,
			Report{
				Name:       registry,
				Kind:       "Insecure",
				Supported:  true,
				Confidence: HighConfidence,
			})
	}

	for _, registry := range e.Registries["search"].List {
		reportOutput.Reports = append(reportOutput.Reports,
			Report{
				Name:       registry,
				Kind:       "Search",
				Supported:  false,
				Confidence: NoConfidence,
				Comment:    "Search registries can not be configured in OCP 4",
			})
	}

	return reportOutput, nil
}

// Extract collects registry information from an OCP3 cluster
func (e RegistriesTransform) Extract() (Extraction, error) {
	logrus.Info("RegistriesTransform::Extract")
	content, err := io.FetchFile(env.Config().GetString("RegistriesConfigFile"))
	if err != nil {
		return nil, err
	}

	var extraction RegistriesExtraction
	if _, err := toml.Decode(string(content), &extraction); err != nil {
		return nil, err
	}
	return extraction, nil

}

// Validate registry data collected from an OCP3 cluster
func (e RegistriesExtraction) Validate() error {
	if len(e.Registries["block"].List) == 0 && len(e.Registries["insecure"].List) == 0 && len(e.Registries["search"].List) == 0 {
		return errors.New("no configured registries detected, not generating a cr or report")
	}
	return nil
}

// Name returns a human readable name for the transform
func (e RegistriesTransform) Name() string {
	return RegistriesComponentName
}
