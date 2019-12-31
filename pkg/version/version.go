package version

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
)

var (
	// wallet's version info
	vCtlMajor, vCtlMinor, vCtlPatch int64 = 0, 1, 0
	// GitHash Value will be set during build
	GitHash = "Not provided"
	// BuildTime Value will be set during build
	BuildTime = "Not provided"
)

// walletVer version of wallet
var walletVer = semver.Version{
	Major: vCtlMajor,
	Minor: vCtlMinor,
	Patch: vCtlPatch,
}

// LogAppInfo 打印版本信息
func LogAppInfo() {
	fmt.Printf("AppVersion: %d\nApiVersion: %d\nGitHash: %s\nBuildTime: %s\n\n",
		walletVer.Major, walletVer.Minor, GitHash, BuildTime)
}
