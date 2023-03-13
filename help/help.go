package help

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type Content struct {
}

type YAMLHelp struct {
	UUID     string   `yaml:"UUID"`
	Keywords []string `yaml:"Keywords""`
	Content  string   `yaml:"Content"`
}

func buildFromPath(path string) YAMLHelp {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Help path error")
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("Error reading help YAML file %s", absPath)
	}
	yamlHelp := YAMLHelp{}
	err = yaml.Unmarshal(data, &yamlHelp)
	if err != nil {
		log.Fatalf("Error marshal help file")
	}
	if yamlHelp.UUID == "" {
		log.Fatalf("Missing UUID: %s", path)
	}
	if len(yamlHelp.Keywords) == 0 {
		log.Fatalf("Missing Keywords: %s", path)
	}
	if yamlHelp.Content == "" {
		log.Fatalf("Missing Content: %s", path)
	}
	return yamlHelp
}

func Build() map[string]YAMLHelp {
	temp := buildFromPath("../data/help/circle.yaml")
	index := make(map[string]YAMLHelp)
	for _, keyword := range temp.Keywords {
		index[keyword] = temp
	}
	return index
}
