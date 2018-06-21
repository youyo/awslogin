package awslogin

import (
	"reflect"
	"testing"

	latest "github.com/tcnksm/go-latest"
)

func TestNewVersions(t *testing.T) {
	v := NewVersions()
	expected := &Versions{}
	if reflect.TypeOf(v) != reflect.TypeOf(expected) {
		t.Errorf("expected %v to eq %v",
			reflect.TypeOf(v), reflect.TypeOf(expected))
	}
}

func TestWithVersion(t *testing.T) {
	input := "test-version"
	v := NewVersions().WithVersion(input)
	expected := input
	if v.Version != expected {
		t.Errorf("expected %v to eq %v",
			v.Version, expected)
	}
}

func TestWithCommitHash(t *testing.T) {
	input := "test-commit-hash"
	v := NewVersions().WithCommitHash(input)
	expected := input
	if v.CommitHash != expected {
		t.Errorf("expected %v to eq %v",
			v.CommitHash, expected)
	}
}

func TestSetGithubTag(t *testing.T) {
	owner := "owner"
	repo := "repo"
	expected := &latest.GithubTag{
		Owner:      owner,
		Repository: repo,
	}
	v := NewVersions().
		WithGithubOwner(owner).
		WithGithubRepository(repo)
	v.SetGithubTag()
	if !reflect.DeepEqual(v.GithubTag, expected) {
		t.Errorf("expected %v to eq %v",
			v.GithubTag, expected)
	}
}

func TestFetchVersionData(t *testing.T) {
	owner := "youyo"
	repo := "awslogin"
	version := "0.0.0"
	v := NewVersions().
		WithGithubOwner(owner).
		WithGithubRepository(repo).
		WithVersion(version)
	v.SetGithubTag()
	res, err := v.FetchVersionData(v.GithubTag, v.Version)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	expected := true
	if res.Outdated != expected {
		t.Errorf("expected %v to eq %v",
			res.Outdated, expected)
	}
	expected = false
	if res.Latest != expected {
		t.Errorf("expected %v to eq %v",
			res.Latest, expected)
	}
}

func TestSetVersionCheckData(t *testing.T) {
	owner := "youyo"
	repo := "awslogin"
	version := "0.0.0"
	v := NewVersions().
		WithGithubOwner(owner).
		WithGithubRepository(repo).
		WithVersion(version)
	v.SetGithubTag()
	res, _ := v.FetchVersionData(v.GithubTag, v.Version)
	v.SetVersionCheckData(res)

	expected := true
	if v.Outdate != expected {
		t.Errorf("expected %v to eq %v",
			v.Outdate, expected)
	}
}
