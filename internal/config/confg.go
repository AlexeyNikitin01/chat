package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Mode       string `yaml:"mode"`
	LogLevel   string `yaml:"logLevel"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	ClientName string `yaml:"clientName"`
}

func Init() (*Config, error) {
	path := parseFlag()
	cfg, err := configYaml(path)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseFlag() *string {
	path := flag.String("p", "etc/server-config.yml", "path to config")
	flag.Parse()
	return path
}

func configYaml(path *string) (*Config, error) {
	content, err := os.ReadFile(*path)
	if err != nil {
		return nil, fmt.Errorf("no correct path or file")
	}
	
	cfg := &Config{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, fmt.Errorf("no unmarshal config")
	}
	
	return cfg, nil
}
