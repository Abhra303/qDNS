package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ServerConfiguration ConfigFile

type ConfigFile struct {
	// TODO: design config file fields
	Zones []struct {
		ZoneName         string
		ZonefileLocation []string
	}
	// Some configurations
}

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func LoadConfigFile(filePath string) (err error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	exists, err := fileExists(absPath)
	if err != nil {
		return err
	}
	if !exists {
		err = fmt.Errorf("the given path doesn't exist")
		return err
	}

	buf, err := os.ReadFile(absPath)
	if err == nil {
		err = yaml.Unmarshal(buf, &ServerConfiguration)
		if err != nil {
			return err
		}
	}

	return err
}
