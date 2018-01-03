package awslogin

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

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
	}
)

func NewConfig() (cfg *Config, err error) {
	cfg = &Config{
		Path: configPath(os.Getenv("HOME")),
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

func (cfg *Config) FetchArn() (err error) {
	s := "profile " + cfg.ProfileName
	cfg.ARN = cfg.Data.Section(s).Key("role_arn").String()
	cfg.SourceProfile = cfg.Data.Section(s).Key("source_profile").String()
	cfg.MfaSerial = cfg.Data.Section(s).Key("mfa_serial").String()
	if cfg.ARN == "" {
		err = errors.New("Could not fetch Arn")
		return
	}
	cfg.RoleSessionName = buildRoleSessionName(cfg.ProfileName, cfg.ARN)
	return
}

func buildRoleSessionName(profileName, arn string) (roleSessionName string) {
	roleSessionName = profileName + "." + strings.Split(arn, ":")[4] + "@awslogin"
	return
}
