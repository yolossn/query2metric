package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const MONGO = "MONGO"
const SQL = "SQL"

type Metric struct {
	Name       string `yaml:"name"`
	HelpString string `yaml:"helpString"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
	Query      string `yaml:"query"`
	Time       int64  `yaml:"time"`
}

type Config struct {
	Connections []struct {
		Name             string   `yaml:"name"`
		Type             string   `yaml:"type"`
		Subtype          string   `yaml:"subtype"`
		ConnectionString string   `yaml:"connectionStringFromEnv"`
		Metrics          []Metric `yaml:"metrics"`
	} `yaml:"connections"`
}

func FromFile(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading file")
	}

	conf := Config{}
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling config")
	}
	return &conf, nil
}
