package version

import (
	"strings"

	"github.com/blang/semver/v4"
)

var (
	Version = "1.2.8"

	//Vars injected at build-time
	BuildTimestamp = ""
)

const (
	// LatestVersion product version supported
	LatestVersion        = "2.38.0"
	CompactLatestVersion = "2380"

	LatestKubeImage = "quay.io/arkmq-org/activemq-artemis-broker-kubernetes:artemis." + LatestVersion
	LatestInitImage = "quay.io/arkmq-org/activemq-artemis-broker-init:artemis." + LatestVersion
)

func DefaultImageName(archSpecificRelatedImageEnvVarName string) string {
	if strings.Contains(archSpecificRelatedImageEnvVarName, "_Init_") {
		return LatestInitImage
	} else {
		return LatestKubeImage
	}
}

var FullVersionFromCompactVersion map[string]string = map[string]string{
	"2210": "2.21.0",
	"2220": "2.22.0",
	"2230": "2.23.0",
	"2250": "2.25.0",
	"2260": "2.26.0",
	"2270": "2.27.0",
	"2271": "2.27.1",
	"2280": "2.28.0",
	"2290": "2.29.0",
	"2300": "2.30.0",
	"2310": "2.31.0",
	"2312": "2.31.2",
	"2320": "2.32.0",
	"2330": "2.33.0",
	"2340": "2.34.0",
	"2350": "2.35.0",
	"2360": "2.36.0",
	"2370": "2.37.0",
	"2380": "2.38.0",
}

// The yacfg profile to use for a given full version of broker
var YacfgProfileVersionFromFullVersion map[string]string = map[string]string{
	"2.21.0": "2.21.0",
	"2.22.0": "2.21.0",
	"2.23.0": "2.21.0",
	"2.25.0": "2.21.0",
	"2.26.0": "2.21.0",
	"2.27.0": "2.21.0",
	"2.27.1": "2.21.0",
	"2.28.0": "2.21.0",
	"2.29.0": "2.21.0",
	"2.30.0": "2.21.0",
	"2.31.0": "2.21.0",
	"2.31.2": "2.21.0",
	"2.32.0": "2.21.0",
	"2.33.0": "2.21.0",
	"2.34.0": "2.21.0",
	"2.35.0": "2.21.0",
	"2.36.0": "2.21.0",
	"2.37.0": "2.21.0",
	"2.38.0": "2.21.0",
}

var YacfgProfileName string = "artemis"

// Sorted array of supported ActiveMQ Artemis versions
var SupportedActiveMQArtemisVersions = []string{
	"2.21.0",
	"2.22.0",
	"2.23.0",
	"2.25.0",
	"2.26.0",
	"2.27.0",
	"2.27.1",
	"2.28.0",
	"2.29.0",
	"2.30.0",
	"2.31.0",
	"2.31.2",
	"2.32.0",
	"2.33.0",
	"2.34.0",
	"2.35.0",
	"2.36.0",
	"2.37.0",
	"2.38.0",
}

func CompactActiveMQArtemisVersion(version string) string {
	return strings.Replace(version, ".", "", -1)
}

var supportedActiveMQArtemisSemanticVersions []semver.Version

func SupportedActiveMQArtemisSemanticVersions() []semver.Version {
	if supportedActiveMQArtemisSemanticVersions == nil {
		supportedActiveMQArtemisSemanticVersions = make([]semver.Version, len(SupportedActiveMQArtemisVersions))
		for i := 0; i < len(SupportedActiveMQArtemisVersions); i++ {
			supportedActiveMQArtemisSemanticVersions[i] = semver.MustParse(SupportedActiveMQArtemisVersions[i])
		}
		semver.Sort(supportedActiveMQArtemisSemanticVersions)
	}

	return supportedActiveMQArtemisSemanticVersions
}

func IsSupportedActiveMQArtemisVersion(version string) bool {
	for i := 0; i < len(SupportedActiveMQArtemisVersions); i++ {
		if SupportedActiveMQArtemisVersions[i] == version {
			return true
		}
	}
	return false
}
