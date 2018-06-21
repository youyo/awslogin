package awslogin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"
)

type Config struct {
	Path            string
	Data            *ini.File
	ARN             string
	SourceProfile   string
	MfaSerial       string
	ProfileName     string
	RoleSessionName string
	DurationSeconds int
}

func NewConfig() *Config {
	return &Config{
		Path:            filepath.Join(os.Getenv("HOME"), ".aws/config"),
		DurationSeconds: 3600,
		RoleSessionName: "awslogin",
	}
}

func (c *Config) SetPath(p string) {
	c.Path = p
}

func (c *Config) SetData() error {
	i, err := ini.Load(c.Path)
	if err != nil {
		return err
	}
	c.Data = i
	return nil
}

func (c *Config) SetARN(a string) {
	c.ARN = a
}

func (c *Config) SetSourceProfile(s string) {
	c.SourceProfile = s
}

func (c *Config) SetMfaSerial(m string) {
	c.MfaSerial = m
}

func (c *Config) SetProfileName(p string) {
	c.ProfileName = p
}

func (c *Config) SetRoleSessionName(r string) {
	c.RoleSessionName = r
}

func (c *Config) SetDurationSeconds(d int) {
	c.DurationSeconds = d
}

func (c *Config) SelectProfile(optProfile string, optReadFromEnv bool) error {
	if optProfile != "" {
		c.SetProfileName(optProfile)
		return nil
	} else if optReadFromEnv {
		c.SetProfileName(os.Getenv("AWS_PROFILE"))
		return nil
	} else {
		p, err := c.Peco()
		c.SetProfileName(p)
		return err
	}
}

func (c *Config) Peco() (string, error) {
	cmd := exec.Command("peco")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()
	aa := c.AvailableArn()
	list := strings.Join(aa[:], "\n")
	io.WriteString(stdin, list)
	stdin.Close()
	byteProfile, err := cmd.Output()
	if err != nil {
		errorMessage := fmt.Sprintf("'%s' is Required command. Please install it. %s ", "peco", "https://github.com/peco/peco")
		err = errors.Wrap(err, errorMessage)
		return "", err
	}
	return strings.TrimRight(string(byteProfile), "\n"), nil
}

func (c *Config) AvailableArn() []string {
	var list []string
	for _, s := range c.Data.Sections() {
		if s.HasKey("role_arn") {
			n := strings.Replace(s.Name(), "profile ", "", 1)
			list = append(list, n)
		}
	}
	return list
}

func (c *Config) FetchArn() error {
	s := "profile " + c.ProfileName
	d := c.Data.Section(s)
	c.SetARN(d.Key("role_arn").String())
	c.SetSourceProfile(d.Key("source_profile").String())
	if c.Data.Section(s).HasKey("duration_seconds") {
		c.SetDurationSeconds(d.Key("duration_seconds").MustInt())
	}
	c.SetMfaSerial(d.Key("mfa_serial").String())
	if c.ARN == "" {
		return errors.New("Could not fetch Arn")
	}
	return nil
}
