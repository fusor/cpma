package transform

import (
	"fmt"

	"github.com/konveyor/cpma/pkg/decode"
	"github.com/konveyor/cpma/pkg/env"
	"github.com/konveyor/cpma/pkg/io"
	"github.com/konveyor/cpma/pkg/transform/reportoutput"
	"github.com/konveyor/cpma/pkg/transform/sdn"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

// SDNComponentName is the SDN component string
const SDNComponentName = "SDN"

// SDNExtraction is an SDN specific extraction
type SDNExtraction struct {
	legacyconfigv1.MasterConfig
}

// SDNTransform is an SDN specific transform
type SDNTransform struct {
}

const clusterNetworkComment = `Networks must be configured during installation,
 hostSubnetLength was replaced with hostPrefix in OCP4, default value was set to 23`

// Transform converts data collected from an OCP3 into a useful output
func (e SDNExtraction) Transform() ([]Output, error) {
	outputs := []Output{}

	if env.Config().GetBool("Manifests") {
		logrus.Info("SDNTransform::Transform:Manifests")
		manifests, err := e.buildManifestOutput()
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, manifests)
	}

	if env.Config().GetBool("Reporting") {
		logrus.Info("SDNTransform::Transform:Reports")
		e.buildReportOutput()
	}

	return outputs, nil
}
func (e SDNExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	networkCR, err := sdn.Translate(e.MasterConfig)
	if err != nil {
		return nil, err
	}

	networkCRYAML, err := GenYAML(networkCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-sdn.yaml", CRD: networkCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e SDNExtraction) buildReportOutput() {
	componentReport := reportoutput.ComponentReport{
		Component: SDNComponentName,
	}

	for _, n := range e.MasterConfig.NetworkConfig.ClusterNetworks {
		cidrComment := fmt.Sprintf("Networks must be configured during installation, it's possible to use %s", n.CIDR)
		componentReport.Reports = append(componentReport.Reports,
			reportoutput.Report{
				Name:       "CIDR",
				Kind:       "ClusterNetwork",
				Supported:  true,
				Confidence: ModerateConfidence,
				Comment:    cidrComment,
			})

		componentReport.Reports = append(componentReport.Reports,
			reportoutput.Report{
				Name:       "HostSubnetLength",
				Kind:       "ClusterNetwork",
				Supported:  false,
				Confidence: NoConfidence,
				Comment:    clusterNetworkComment,
			})
	}

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       e.MasterConfig.NetworkConfig.ServiceNetworkCIDR,
			Kind:       "ServiceNetwork",
			Supported:  true,
			Confidence: ModerateConfidence,
			Comment:    "Networks must be configured during installation",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "ExternalIPNetworkCIDRs",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Configuration of ExternalIPNetworkCIDRs is not supported in OCP4",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "IngressIPNetworkCIDR",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of this configuration is not supported, refer to ingress operator configuration for more information",
		})

	FinalReportOutput.Report.ComponentReports = append(FinalReportOutput.Report.ComponentReports, componentReport)
}

// Extract collects SDN configuration information from an OCP3 cluster
func (e SDNTransform) Extract() (Extraction, error) {
	logrus.Info("SDNTransform::Extract")

	content, err := io.FetchFile(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		return nil, err
	}

	var extraction SDNExtraction
	extraction.MasterConfig = *masterConfig

	return extraction, nil
}

// Validate the data extracted from the OCP3 cluster
func (e SDNExtraction) Validate() error {
	if err := sdn.Validate(e.MasterConfig); err != nil {
		return err
	}

	return nil
}

// Name returns a human readable name for the transform
func (e SDNTransform) Name() string {
	return SDNComponentName
}
