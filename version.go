package awslogin

import (
	"fmt"

	latest "github.com/tcnksm/go-latest"
)

const Version = "1.1.9"

type Versions struct {
	Version    string
	CommitHash string
	BuildTime  string
	GoVersion  string

	GithubOwner      string
	GithubRepository string
	LatestVersion    string
	Outdate          bool
	GithubTag        *latest.GithubTag
}

func NewVersions() *Versions {
	return new(Versions)
}

func (v *Versions) WithVersion(s string) *Versions {
	v.Version = s
	return v
}

func (v *Versions) WithCommitHash(s string) *Versions {
	v.CommitHash = s
	return v
}

func (v *Versions) WithBuildTime(s string) *Versions {
	v.BuildTime = s
	return v
}

func (v *Versions) WithGoVersion(s string) *Versions {
	v.GoVersion = s
	return v
}

func (v *Versions) WithGithubOwner(s string) *Versions {
	v.GithubOwner = s
	return v
}

func (v *Versions) WithGithubRepository(s string) *Versions {
	v.GithubRepository = s
	return v
}

func (v *Versions) SetLatestVersion(s string) {
	v.LatestVersion = s
}

func (v *Versions) SetOutdate(b bool) {
	v.Outdate = b
}

func (v *Versions) SetGithubTag() {
	v.GithubTag = &latest.GithubTag{
		Owner:      v.GithubOwner,
		Repository: v.GithubRepository,
	}
}

func (v *Versions) FetchVersionData(source latest.Source, target string) (*latest.CheckResponse, error) {
	return latest.Check(source, target)
}

func (v *Versions) SetVersionCheckData(res *latest.CheckResponse) {
	v.SetOutdate(res.Outdated)
	v.SetLatestVersion(res.Current)
}

func (v *Versions) OutputVersionMessage() string {
	m := fmt.Sprintf("Version: %s\nCommitHash: %s\nBuildTime: %s\nGoVersion: %s\n",
		v.Version, v.CommitHash, v.BuildTime, v.GoVersion)
	if v.Outdate {
		m += fmt.Sprintf("\n%s is not latest, you should upgrade to %s\n",
			v.Version, v.LatestVersion)
	}
	return m
}
