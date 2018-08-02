package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
)

var (
	version           string
	ignorePrevious    bool
	incPattern        string
	versionPrefix     string
	versionPrerelease string
	versionMetadata   string
)

func init() {
	version = os.Getenv("version_base")
	ignorePrevious, _ = strconv.ParseBool(os.Getenv("version_ignorePrevious"))
	incPattern = os.Getenv("version_pattern")
	versionPrerelease = os.Getenv("version_prerelease")
	versionPrefix = os.Getenv("version_prefix")
	versionMetadata = os.Getenv("version_metadata")
}

func getCommitSHA() string {
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Unable to determine current commit SHA %v\n", err)
		fmt.Printf("%s\n", stdErr.String())
		os.Exit(1)
	}
	return strings.TrimSpace(out.String())
}

func getLatestGitTag(sha string) string {
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd := exec.Command("git", "describe", "--tag", "--abbrev=0", fmt.Sprintf("%s", sha))
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return "0.0.0"
	}
	return strings.TrimSpace(out.String())
}

func getTagCommitSha(tag string) string {
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd := exec.Command("git", "rev-list", "-n 1", fmt.Sprintf("%s", tag))
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return "0.0.0"
	}
	return strings.TrimSpace(out.String())
}

func timeSinceLastTag(tag string) string {
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd := exec.Command("git", "log", "--pretty=%cI", "-n 1", fmt.Sprintf("%s", getTagCommitSha(tag)))
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	var topOutput string
	if err != nil {
		cmd = exec.Command("git", "log", "--pretty=%cI", "-n 1", fmt.Sprintf("%s", getCommitSHA()))
		cmd.Stdout = &out
		cmd.Stderr = &stdErr
		cmd.Run()
	}
	topOutput = strings.TrimSpace(out.String())

	t, err := time.Parse(time.RFC3339, topOutput)
	if err != nil {
		fmt.Printf("Git did not return expected time format or there was a failure parsing the time [%s] %v\n", topOutput, err)
		fmt.Printf("%s\n", topOutput)
		os.Exit(1)
	}

	timeSince := time.Since(t)

	return fmt.Sprintf("b%09.f", timeSince.Seconds())
}

// Process entrypoint to do the work and determine when/how to increment.
func Process(v, pattern, prerelease, metadata, prefix string, ignore bool) (string, error) {
	lastTag := getLatestGitTag(getCommitSHA())
	if ignore && v == "" {
		return "", fmt.Errorf("ignore previous is true but no base version provided, please check input")
	} else if !ignore {
		v = lastTag
	}

	if prerelease == "" {
		prerelease = timeSinceLastTag(lastTag)
	}

	semVer, err := semver.NewVersion(v)
	if err != nil {
		return "", fmt.Errorf("version parse failed, [%v]. %s is not a valid SemVer tag", err, v)
	}

	outVersion := incVersion(*semVer, pattern)

	t, err := outVersion.SetPrerelease(prerelease)
	if err != nil {
		return "", fmt.Errorf("error appending pre-release data: %v", err)
	}
	outVersion = t

	if metadata != "" {
		t, err := outVersion.SetMetadata(metadata)
		if err != nil {
			return "", fmt.Errorf("error appending metadata: %v", err)
		}
		outVersion = t
	}

	outString := outVersion.String()

	if prefix != "" {
		outString = fmt.Sprintf("%s%s", prefix, outString)
	}

	return outString, nil
}

func determineFromPattern(pattern string) (major bool, minor bool, patch bool) {
	splitPattern := strings.Split(pattern, ".")

	// loop through the pattern so we can determine which pieces should increment
	for i, p := range splitPattern {
		if p == "^" {
			// determine which number to increment based on the index of the ^
			switch i {
			case 0:
				major = true
				break
			case 1:
				minor = true
				break
			case 2:
				patch = true
				break
			}
		}
	}

	return
}

func incVersion(version semver.Version, pattern string) semver.Version {
	major, minor, patch := determineFromPattern(pattern)

	// no increment desired
	if !major && !minor && !patch {
		minor = true
	}

	if major {
		version = version.IncMajor()
	}

	if minor {
		version = version.IncMinor()
	}

	if patch {
		version = version.IncPatch()
	}

	return version
}

func main() {
	v, err := Process(version, incPattern, versionPrerelease, versionMetadata, versionPrefix, ignorePrevious)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", v)
}
