package awslogin

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	_, err := NewConfig("../../tests/config")
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}

	_, err = NewConfig("../../tests/not_exist_config")
	expected = errors.New("open ../../tests/not_exist_config: no such file or directory")
	if err == nil {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}

func TestConfigPath(t *testing.T) {
	c := configPath()
	var expected string
	if reflect.TypeOf(c) != reflect.TypeOf(expected) {
		t.Errorf("expected %v to eq %v", reflect.TypeOf(c), reflect.TypeOf(expected))
	}
}

func TestLoadConfig(t *testing.T) {
	_, err := loadConfig("../../tests/config")
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}

	_, err = loadConfig("../../tests/not_exist_config")
	expected = errors.New("open ../../tests/not_exist_config: no such file or directory")
	if err == nil {
		t.Errorf("expected %v to eq %v", err, expected)
	}

}

func TestAvailableArn(t *testing.T) {
	cfg, _ := NewConfig("../../tests/config")
	list := cfg.AvailableArn()
	var expected []string
	if reflect.TypeOf(list) != reflect.TypeOf(expected) {
		t.Errorf("expected %v to eq %v", reflect.TypeOf(list), reflect.TypeOf(expected))
	}
}

func TestSetProfileName(t *testing.T) {
	cfg, _ := NewConfig("../../tests/config")
	cfg.SetProfileName("test")
	expected := "test"
	if cfg.ProfileName != expected {
		t.Errorf("expected %v to eq %v", cfg.ProfileName, expected)
	}
}

func TestSetMfaCode(t *testing.T) {
	cfg, _ := NewConfig("../../tests/config")
	cfg.SetMfaCode("123456")
	expected := "123456"
	if cfg.MfaCode != expected {
		t.Errorf("expected %v to eq %v", cfg.ProfileName, expected)
	}
}

func TestFetchArn(t *testing.T) {
	cfg, _ := NewConfig("../../tests/config")
	cfg.SetProfileName("test")
	err := cfg.FetchArn()
	var expected error
	if err != expected {
		t.Errorf("expected %v to eq %v", err, expected)
	}

	cfg.SetProfileName("no_exist_profile")
	err = cfg.FetchArn()
	expected = errors.New("Could not fetch Arn")
	if err == nil {
		t.Errorf("expected %v to eq %v", err, expected)
	}
}
