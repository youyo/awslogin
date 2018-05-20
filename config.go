package awslogin

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	ini "gopkg.in/ini.v1"
)

type (
	Config struct {
		Path            string
		Data            *ini.File
		ARN             string
		SourceProfile   string
		MfaSerial       string
		ProfileName     string
		RoleSessionName string
		DurationSeconds int
	}
)

func NewConfig() (cfg *Config, err error) {
	cfg = &Config{
		Path:            configPath(os.Getenv("HOME")),
		DurationSeconds: 3600,
	}
	cfg.Data, err = loadConfig(cfg.Path)
	return
}

func configPath(home string) (c string) {
	c = filepath.Join(home, ".aws/config")
	return
}

func loadConfig(configPath string) (cfg *ini.File, err error) {
	cfg, err = ini.Load(configPath)
	return
}

func (cfg *Config) AvailableArn() (list []string) {
	sections := cfg.Data.Sections()
	for _, s := range sections {
		if s.HasKey("role_arn") {
			n := strings.Replace(s.Name(), "profile ", "", 1)
			list = append(list, n)
		}
	}
	return
}

func (cfg *Config) SetProfileName(profileName string) {
	cfg.ProfileName = profileName
}

func (cfg *Config) SetDurationSeconds(durationSeconds int) {
	cfg.DurationSeconds = durationSeconds
}

func (cfg *Config) FetchArn() (err error) {
	s := "profile " + cfg.ProfileName
	cfg.ARN = cfg.Data.Section(s).Key("role_arn").String()
	cfg.SourceProfile = cfg.Data.Section(s).Key("source_profile").String()
	if cfg.Data.Section(s).HasKey("duration_seconds") {
		cfg.DurationSeconds = cfg.Data.Section(s).Key("duration_seconds").MustInt()
	}
	cfg.MfaSerial = cfg.Data.Section(s).Key("mfa_serial").String()
	if cfg.ARN == "" {
		err = errors.New("Could not fetch Arn")
		return
	}
	cfg.RoleSessionName = buildRoleSessionName(cfg.ARN)
	return
}

func buildRoleSessionName(arn string) (roleSessionName string) {
	roleName := strings.Split(arn, "/")[1]
	awsAccountID := strings.Split(arn, ":")[4]
	roleSessionName = roleName + "." + awsAccountID + "@awslogin"
	if utf8.RuneCountInString(roleSessionName) > 64 {
		roleSessionName = awsAccountID + "@awslogin"
	}
	return
}
