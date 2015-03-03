package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindingConfigDir(t *testing.T) {
	//ensure env is clean
	os.Setenv("BIGV_CONFIG_DIR", "")

	config := NewConfig("")
	expected := filepath.Join(os.Getenv("HOME"), "/.go-bigv")
	if config.Dir != expected {
		t.Errorf("HOME test failed, expected %s, got %s", expected, config.Dir)
	}

	expected = "/tmp"
	os.Setenv("BIGV_CONFIG_DIR", expected)

	config = NewConfig("")
	if config.Dir != expected {
		t.Errorf("ENV test failed, expected %s, got %s", expected, config.Dir)
	}

	expected = "/home"

	config = NewConfig(expected)
	if config.Dir != expected {
		t.Errorf("--config-dir test failed, expected %s, got %s", expected, config.Dir)
	}
}
