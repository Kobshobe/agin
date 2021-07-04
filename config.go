package agin

import (
	"gopkg.in/yaml.v2"
	"log"
)

func GetConfigFromYAML(configSource []byte, config interface{}) {
	err := yaml.Unmarshal(configSource, config)
	if err != nil {
		log.Fatal(err)
	}
}