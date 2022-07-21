package version

import "strings"

var (
	Version = "1.0.4"
	// PriorVersion - prior version
	PriorVersion = "1.0.3"
)

const (
	// LatestVersion product version supported
	LatestVersion        = "1.0.7"
	CompactLatestVersion = "107"
	// LastMicroVersion product version supported
	LastMicroVersion = "1.0.3"
	// LastMinorVersion product version supported
	LastMinorVersion = "2.15.0"

	LatestKubeImage = "quay.io/artemiscloud/activemq-artemis-broker-kubernetes:" + LatestVersion
	LatestInitImage = "quay.io/artemiscloud/activemq-artemis-broker-init:" + LatestVersion
)

func DefaultImageName(archSpecificRelatedImageEnvVarName string) string {
	if strings.Contains(archSpecificRelatedImageEnvVarName, "_Init_") {
		return LatestInitImage
	} else {
		return LatestKubeImage
	}
}

// SupportedVersions - product versions this operator supports
var SupportedVersions = []string{LatestVersion, LastMicroVersion, LastMinorVersion}
var SupportedMicroVersions = []string{LatestVersion, LastMicroVersion}
var OperandVersionFromOperatorVersion map[string]string = map[string]string{
	"0.18.0": "2.15.0",
	"0.19.0": "2.16.0",
	"0.20.0": "2.18.0",
	"0.20.1": "2.18.0",
	"1.0.0":  "2.20.0",
	"1.0.1":  "1.0.1",
	"1.0.2":  "1.0.2",
	"1.0.3":  "1.0.5",
	"1.0.4":  "1.0.7",
}
var FullVersionFromMinorVersion map[string]string = map[string]string{
	"150": "2.15.0",
	"160": "2.16.0",
	"180": "2.18.0",
	"200": "2.20.0",
	"100": "2.20.0",
	"101": "1.0.1",
	"102": "1.0.2",
	"103": "1.0.3",
	"104": "1.0.4",
	"107": "1.0.7",
}

var CompactFullVersionFromMinorVersion map[string]string = map[string]string{
	"150": "2150",
	"160": "2160",
	"180": "2180",
	"200": "2200",
	"100": "2200",
	"101": "101",
	"102": "102",
	"103": "103",
	"104": "104",
	"107": "107",
}

var CompactVersionFromVersion map[string]string = map[string]string{
	"2.15.0": "2150",
	"2.16.0": "2160",
	"2.18.0": "2180",
	"2.20.0": "2200",
	"1.0.0":  "2200",
	"1.0.1":  "101",
	"1.0.2":  "102",
	"1.0.3":  "103",
	"1.0.4":  "104",
	"1.0.7":  "107",
}

var FullVersionFromCompactVersion map[string]string = map[string]string{
	"2150": "2.15.0",
	"2160": "2.16.0",
	"2180": "2.18.0",
	"2200": "2.20.0",
	"101":  "1.0.1",
	"102":  "1.0.2",
	"103":  "1.0.3",
	"104":  "1.0.4",
	"107":  "1.0.7",
}

var MinorVersionFromFullVersion map[string]string = map[string]string{
	"2.15.0": "150",
	"2.16.0": "160",
	"2.18.0": "180",
	"2.20.0": "200",
	"1.0.0":  "200",
	"1.0.1":  "01",
	"1.0.2":  "02",
	"1.0.3":  "03",
	"1.0.4":  "04",
	"1.0.7":  "07",
}

//The yacfg profile to use for a given full version of broker
var YacfgProfileVersionFromFullVersion map[string]string = map[string]string{
	"2.15.0": "2.15.0",
	"2.16.0": "2.16.0",
	"2.18.0": "2.18.0",
	"2.20.0": "2.18.0",
	"1.0.1":  "2.21.0",
	"1.0.2":  "2.21.0",
	"1.0.3":  "2.21.0",
	"1.0.4":  "2.21.0",
	"1.0.7":  "2.21.0",
}

var YacfgProfileName string = "artemis"
