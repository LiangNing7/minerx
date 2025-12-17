package app

import (
	"strconv"

	"github.com/blang/semver/v4"
	kversion "k8s.io/apimachinery/pkg/version"

	"github.com/LiangNing7/goutils/pkg/version"
)

func convertVersion(info version.Info) *kversion.Info {
	v, _ := semver.Make(info.GitVersion)
	return &kversion.Info{
		Major:        strconv.FormatUint(v.Major, 10),
		Minor:        strconv.FormatUint(v.Minor, 10),
		GitVersion:   info.GitVersion,
		GitCommit:    info.GitCommit,
		GitTreeState: info.GitTreeState,
		BuildDate:    info.BuildDate,
		GoVersion:    info.GoVersion,
		Compiler:     info.Compiler,
		Platform:     info.Platform,
	}
}
