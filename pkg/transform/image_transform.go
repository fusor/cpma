package transform

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/konveyor/cpma/pkg/decode"
	"github.com/konveyor/cpma/pkg/env"
	"github.com/konveyor/cpma/pkg/io"
	"github.com/konveyor/cpma/pkg/transform/image"
	"github.com/konveyor/cpma/pkg/transform/registries"
	"github.com/konveyor/cpma/pkg/transform/reportoutput"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ImageComponentName is the Image component string
const ImageComponentName = "Image"

// ImageExtraction is image specific extraction
type ImageExtraction struct {
	MasterConfig     legacyconfigv1.MasterConfig
	RegistriesConfig RegistriesExtraction
}

// RegistriesExtraction holds registry information extracted from an OCP3 cluster
type RegistriesExtraction struct {
	Registries map[string]registries.Registries
}

// Registries holds a list of Registries
type Registries struct {
	List []string `toml:"registries"`
}

// ImageTransform is an image specific transform
type ImageTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e ImageExtraction) Transform() ([]Output, error) {
	outputs := []Output{}

	if env.Config().GetBool("Manifests") {
		logrus.Info("ImageTransform::Transform:Manifests")
		manifests, err := e.buildManifestOutput()
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, manifests)
	}

	if env.Config().GetBool("Reporting") {
		logrus.Info("ImageTransform::Transform:Reports")
		e.buildReportOutput()
	}

	return outputs, nil
}

func (e ImageExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	const (
		apiVersion = "config.openshift.io/v1"
		kind       = "Image"
		name       = "cluster"
		annokey    = "release.openshift.io/create-only"
		annoval    = "true"
	)

	metadata := metav1.ObjectMeta{
		Name:        name,
		Annotations: map[string]string{annokey: annoval},
	}

	var imageCR configv1.Image
	imageCR.APIVersion = apiVersion
	imageCR.Kind = kind
	imageCR.ObjectMeta = metadata

	imageCR.Spec.RegistrySources.BlockedRegistries = e.RegistriesConfig.Registries["block"].List
	imageCR.Spec.RegistrySources.InsecureRegistries = e.RegistriesConfig.Registries["insecure"].List

	if err := image.Translate(&imageCR, e.MasterConfig.ImagePolicyConfig); err != nil {
		return nil, err
	}

	imageCRYAML, err := GenYAML(imageCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-image.yaml", CRD: imageCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e ImageExtraction) buildReportOutput() {
	componentReport := reportoutput.ComponentReport{
		Component: ImageComponentName,
	}

	for _, registry := range e.RegistriesConfig.Registries["block"].List {
		componentReport.Reports = append(componentReport.Reports,
			reportoutput.Report{
				Name:       "Blocked",
				Kind:       "Registries",
				Supported:  true,
				Confidence: HighConfidence,
				Comment:    registry,
			})
	}

	for _, registry := range e.RegistriesConfig.Registries["insecure"].List {
		componentReport.Reports = append(componentReport.Reports,
			reportoutput.Report{
				Name:       "Insecure",
				Kind:       "Registries",
				Supported:  true,
				Confidence: HighConfidence,
				Comment:    registry,
			})
	}

	for _, registry := range e.RegistriesConfig.Registries["search"].List {
		componentReport.Reports = append(componentReport.Reports,
			reportoutput.Report{
				Name:       "Search",
				Kind:       "Registries",
				Supported:  false,
				Confidence: NoConfidence,
				Comment:    "Search registries can not be configured in OCP 4: " + registry,
			})
	}

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "AllowedRegistriesForImport",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  true,
			Confidence: HighConfidence,
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "AdditionalTrustedCA",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Each registry must provide its own self-signed CA",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "ExternalRegistryHostname",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  true,
			Confidence: HighConfidence,
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "InternalRegistryHostname",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Set by OCP4 image registry operator",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "DisableScheduledImport",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Not supported by OCP4",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "MaxImagesBulkImportedPerRepository",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Not supported by OCP4",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "MaxScheduledImageImportsPerMinute",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Not supported by OCP4",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "ScheduledImageImportMinimumIntervalSeconds",
			Kind:       "MasterConfig.ImagePolicyConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Not supported by OCP4",
		})

	FinalReportOutput.Report.ComponentReports = append(FinalReportOutput.Report.ComponentReports, componentReport)
}

// Extract collects image configuration information from an OCP3 cluster
func (e ImageTransform) Extract() (Extraction, error) {
	logrus.Info("ImageTransform::Extract")
	var extraction ImageExtraction

	content, err := io.FetchFile(env.Config().GetString("RegistriesConfigFile"))
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(content), &extraction.RegistriesConfig); err != nil {
		return nil, err
	}

	content, err = io.FetchFile(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}
	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		return nil, err
	}
	extraction.MasterConfig = *masterConfig

	return extraction, nil
}

// Validate the data extracted from the OCP3 cluster
func (e ImageExtraction) Validate() error {
	err1 := registries.Validate(e.RegistriesConfig.Registries)
	err2 := image.Validate(e.MasterConfig.ImagePolicyConfig)

	if err1 != 0 && err2 != 0 {
		return errors.New("no configured registries and image detected, not generating CR and/or report")
	}

	return nil
}

// Name returns a human readable name for the transform
func (e ImageTransform) Name() string {
	return "Image"
}
