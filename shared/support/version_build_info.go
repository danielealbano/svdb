package shared_support

import (
	"os"
	"path/filepath"
)

type versionBuildInfo struct {
	version       string
	commit        string
	buildDate     string
	builtBy       string
	goLangVersion string
}

var versionBuildInfoData = versionBuildInfo{
	version:       "0.0.0-dev",
	commit:        "0000000",
	buildDate:     "1970-01-01T00:00:00Z",
	builtBy:       "unknown",
	goLangVersion: "0.0.0",
}

func SetVersionBuildInfo(version, commit, buildDate, builtBy, goLangVersion string) {
	if len(version) > 0 {
		versionBuildInfoData.version = version
	}

	if len(commit) > 0 {
		versionBuildInfoData.commit = commit
	}

	if len(buildDate) > 0 {
		versionBuildInfoData.buildDate = buildDate
	}

	if len(builtBy) > 0 {
		versionBuildInfoData.builtBy = builtBy
	}

	if len(goLangVersion) > 0 {
		versionBuildInfoData.goLangVersion = goLangVersion
	}
}

func GetVersion() string {
	return versionBuildInfoData.version
}

func GetCommit() string {
	return versionBuildInfoData.commit
}

func GetBuildDate() string {
	return versionBuildInfoData.buildDate
}

func GetBuiltBy() string {
	return versionBuildInfoData.builtBy
}

func GetGoLangVersion() string {
	return versionBuildInfoData.goLangVersion
}

func GetExecutableName() string {
	executable, _ := os.Executable()
	executableFilename := filepath.Base(executable)
	return executableFilename
}
