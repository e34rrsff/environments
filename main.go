package main

import (
	_ "embed"
	"fmt"
	"log"

	"go.yaml.in/yaml/v4"
)

type environment struct {
	Title string `yaml:"environment"`
	//TODO: implement other fields for the YAML
	//Run []Action
}

type environmentList map[string]environment

type environmentsConfig struct {
	SessionUser string   `yaml:"default_username"`
	Environments       environmentList `yaml:"environments"`
}

// Function to interpret the YAML file
func (um *environmentList) UnmarshalYAML(val *yaml.Node) error {
	var slice []environment
	if err := val.Decode(&slice); err != nil {
		return err
	}

	*um = make(environmentList)
	for _, env := range slice {
		(*um)[env.Title] = env
	}

	return nil
}

// TODO: make the yaml file an argument to be specified on the command line, no
// embedding

//go:embed environments.yaml
var environments_yaml []byte
var cfg environmentsConfig

func init() {
	if err := yaml.Unmarshal(environments_yaml, &cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Printf("%+v\n", &cfg)
}
