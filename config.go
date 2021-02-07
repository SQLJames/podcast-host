package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"io/ioutil"
)

// ProgramConfig struct for all the program configs.
type ProgramConfig struct {
	RSS RSSConfig `yaml:"RSS"`
	Static StaticConfig `yaml:"Static"`
}
// RSSConfig struct for all the RSS configs.
type RSSConfig struct {
	SearchFolder    string `yaml:"SearchFolder"`
	PodcastFilename string `yaml:"PodcastFilename"`
}

// StaticConfig struct for all the static file configs.
type StaticConfig struct {
	EpisodeLocation string `yaml:"EpisodeLocation"`
	Favicon         string `yaml:"Favicon"`
	Images          string `yaml:"Images"`
}

func initializeConfig(file string) (ProgramConfig, error) {
	var programconfiguration ProgramConfig
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
		return programconfiguration, err
	}
	err = yaml.Unmarshal(yamlFile, &programconfiguration)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %s\n", err)
		return programconfiguration, err
	}
	return programconfiguration, err
}