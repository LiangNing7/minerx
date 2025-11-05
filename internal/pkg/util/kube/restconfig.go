package kube

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	restclient "k8s.io/client-go/rest"

	"github.com/LiangNing7/goutils/pkg/version"
)

const (
	unknowString = "unknow"
)

// buildUserAgent builds a User-Agent string from given args.
func buildUserAgent(command, version, os, arch, commit string) string {
	return fmt.Sprintf(
		"%s/%s (%s/%s) minerx.io/%s", command, version, os, arch, commit)
}

// DefaultMinerXUserAgent returns a User-Agent string build from static global vars.
func DefaultMinerXUserAgent() string {
	return buildUserAgent(
		adjustCommand(os.Args[0]),
		adjustVersion(version.Get().GitVersion),
		runtime.GOOS,
		runtime.GOARCH,
		adjustCommit(version.Get().GitCommit))
}

// SetMinerXDefaults sets default values on the provided client config for accessing the
// MinerX API or returns an error if any of the defaults are impossible or invalid.
func SetMinerXDefaults(config *restclient.Config) {
	if len(config.UserAgent) == 0 {
		config.UserAgent = DefaultMinerXUserAgent()
	}
}

// adjustSourceName returns the name of the source calling the client.
func adjustSourceName(c string) string {
	if c == "" {
		return unknowString
	}
	return c
}

// adjustCommit returns sufficient significant figures of the commit's git hash.
func adjustCommit(c string) string {
	if c == "" {
		return unknowString
	}
	if len(c) > 7 {
		return c[:7]
	}
	return c
}

// adjustVersion strips "alpha", "beta", etc. from version in form
// major.minor.patch-[alpha|beta|etc].
func adjustVersion(v string) string {
	if v == "" {
		return unknowString
	}
	seg := strings.SplitN(v, "-", 2)
	return seg[0]
}

// adjustCommand returns the last component of the
// OS-specific command path for use in User-Agent.
func adjustCommand(p string) string {
	// Unlikely, but better than returning "".
	if p == "" {
		return unknowString
	}
	return filepath.Base(p)
}

func GetUserAgent(userAgent string) string {
	return DefaultMinerXUserAgent() + "/" + adjustSourceName(userAgent)
}

func AddUserAgent(config *restclient.Config, userAgent string) *restclient.Config {
	config.UserAgent = GetUserAgent(userAgent)
	return config
}
