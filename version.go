package zylog

import (
	"fmt"
	"runtime"
)

var (
	// Version is populated at compile from ./VERSION.
	Version string

	// GitCommit is populated at compile time.
	GitCommit string

	// GitSummary is populated at compile time.
	GitSummary string

	// GitBranch is current branch name the code is built off.
	GitBranch string

	// BuildDate is RFC3339 formatted UTC date.
	BuildDate string
)

// VersionString returns a version string as set by the Makefile when the
// library was last compiled. If a version cannot be extracted, the string 'NA'
// is returned.
func VersionString() string {
	if Version == "" {
		return "N/A"
	}
	return Version
}

// BuildString returns a build string that provides the git branch upon which
// the build was made, the git commit of that branch, as well as the build
// date. If a build string cannot be constructed, the string 'NA' is returned.
func BuildString() string {
	if GitCommit == "" {
		return "N/A"
	}
	return fmt.Sprintf("%s@%s, %s", GitBranch, GitCommit, BuildDate)
}

// PrintVersions outputs version information for zylog, build details, and Go version.
func PrintVersions() {
	fmt.Printf("zylog version: %s\n", VersionString())
	fmt.Printf("Build: %s\n", BuildString())
	fmt.Printf("Go version: %s\n", runtime.Version())
}
