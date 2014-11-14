package main

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"path"
)

// RawBox is the data type for a box in the wercker.yml
type RawBox string

// RawServices is a list of auxilliary boxes to boot in the wercker.yml
type RawServices []RawBox

// RawBuild is the data type for builds in the wercker.yml
type RawBuild struct {
	RawSteps []interface{} `yaml:"steps"`
}

// RawConfig is the data type for wercker.yml
type RawConfig struct {
	SourceDir   string      `yaml:"source-dir"`
	RawBox      *RawBox     `yaml:"box"`
	RawServices RawServices `yaml:"services"`
	RawBuild    *RawBuild   `yaml:"build"`
}

// RawStep is the data type for a step in wercker.yml
type RawStep map[string]RawStepData

// RawStepData is the data type for the contents of a step in wercker.yml
type RawStepData map[string]string

func findYaml(searchDirs []string) (string, error) {
	possibleYaml := []string{"ewok.yml", "wercker.yml", ".wercker.yml"}

	for _, v := range searchDirs {
		for _, y := range possibleYaml {
			possibleYaml := path.Join(y, v)
			ymlExists, err := exists(possibleYaml)
			if err != nil {
				return "", err
			}
			if !ymlExists {
				continue
			}
			return possibleYaml, nil
		}
	}
	return "", fmt.Errorf("No wercker.yml found")
}

// ReadWerckerYaml will try to find a wercker.yml file and return its bytes.
// TODO(termie): If allowDefault is true it will try to generate a
// default yaml file by inspecting the project.
func ReadWerckerYaml(searchDirs []string, allowDefault bool) ([]byte, error) {
	foundYaml, err := findYaml(searchDirs)
	if err != nil {
		return nil, err
	}

	// TODO(termie): If allowDefault, we'd generate something here
	// if !allowDefault && !found {
	//   return nil, errors.New("No wercker.yml found and no defaults allowed.")
	// }

	return ioutil.ReadFile(foundYaml)
}

// ConfigFromYaml reads a []byte as yaml and turn it into a RawConfig object
func ConfigFromYaml(file []byte) (*RawConfig, error) {
	var m RawConfig

	err := yaml.Unmarshal(file, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}