package awslogin

import (
	"errors"
	"os"
	"reflect"
	"testing"

	ini "gopkg.in/ini.v1"
)

func TestNewConfig(t *testing.T) {
	pwd, _ := os.Getwd()
	os.Setenv("HOME", pwd+"/tests")
	cfg, err := NewConfig()
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}

	var expectedConfig *Config
	if reflect.TypeOf(cfg) != reflect.TypeOf(expectedConfig) {
		t.Errorf("expected %v to eq %v", reflect.TypeOf(cfg), reflect.TypeOf(expectedConfig))
	}

}

func TestConfigPath(t *testing.T) {
	c := configPath(os.Getenv("HOME"))
	var expected string
	if reflect.TypeOf(c) != reflect.TypeOf(expected) {
		t.Errorf("expected %v to eq %v", reflect.TypeOf(c), reflect.TypeOf(expected))
	}
}

func TestLoadConfig(t *testing.T) {
	pwd, _ := os.Getwd()
	os.Setenv("HOME", pwd+"/tests")
	c := configPath(os.Getenv("HOME"))
	cfg, err := loadConfig(c)
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}

	var expectedConfig *ini.File
	if reflect.TypeOf(cfg) != reflect.TypeOf(expectedConfig) {
		t.Errorf("expected %v to eq %v", reflect.TypeOf(cfg), reflect.TypeOf(expectedConfig))
	}

}

func TestAvailableArn(t *testing.T) {
	pwd, _ := os.Getwd()
	os.Setenv("HOME", pwd+"/tests")
	cfg, _ := NewConfig()
	list := cfg.AvailableArn()
	var expected []string
	if reflect.TypeOf(list) != reflect.TypeOf(expected) {
		t.Errorf("expected %v to eq %v", reflect.TypeOf(list), reflect.TypeOf(expected))
	}
}

func TestSetProfileName(t *testing.T) {
	pwd, _ := os.Getwd()
	os.Setenv("HOME", pwd+"/tests")
	cfg, _ := NewConfig()
	cfg.SetProfileName("test")
	expected := "test"
	if cfg.ProfileName != expected {
		t.Errorf("expected %v to eq %v", cfg.ProfileName, expected)
	}
}

func TestSetDurationSeconds(t *testing.T) {
	pwd, _ := os.Getwd()
	os.Setenv("HOME", pwd+"/tests")
	cfg, _ := NewConfig()
	cfg.SetDurationSeconds(4000)
	expected := 4000
	if cfg.DurationSeconds != expected {
		t.Errorf("expected %v to eq %v", cfg.DurationSeconds, expected)
	}
}

func TestFetchArn(t *testing.T) {
	pwd, _ := os.Getwd()
	os.Setenv("HOME", pwd+"/tests")
	cfg, _ := NewConfig()
	cfg.SetProfileName("test")
	err := cfg.FetchArn()
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}

	var expectedARN string = "arn:aws:iam::xxxxxxxxxxxx:role/test"
	if cfg.ARN != expectedARN {
		t.Errorf("expected %v to eq %v", err, expectedARN)
	}

	var expectedSourceProfile string = "default"
	if cfg.SourceProfile != expectedSourceProfile {
		t.Errorf("expected %v to eq %v", err, expectedSourceProfile)
	}

	var expectedMfaSerial string = "arn:aws:iam::123456789012:mfa/jonsmith"
	if cfg.MfaSerial != expectedMfaSerial {
		t.Errorf("expected %v to eq %v", err, expectedMfaSerial)
	}

	var expectedDurationSeconds int = 43200
	if cfg.DurationSeconds != expectedDurationSeconds {
		t.Errorf("expected %v to eq %v", err, expectedDurationSeconds)
	}

	var expectedRoleSessionName string = "test.xxxxxxxxxxxx@awslogin"
	if cfg.RoleSessionName != expectedRoleSessionName {
		t.Errorf("expected %v to eq %v", err, expectedRoleSessionName)
	}

	cfg.SetProfileName("no_exist_profile")
	err = cfg.FetchArn()
	expected = errors.New("Could not fetch Arn")
	if err == nil {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}

func TestBuildRoleSessionName(t *testing.T) {
	roleSessionName := buildRoleSessionName("arn:aws:iam::xxxxxxxxxxxx:role/test")
	var expected string = "test.xxxxxxxxxxxx@awslogin"
	if roleSessionName != expected {
		t.Errorf("expected %v to eq %v", roleSessionName, expected)
	}
}
