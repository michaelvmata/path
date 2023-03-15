package help

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Content struct {
}

type YAMLHelp struct {
	UUID     string   `yaml:"UUID"`
	Keywords []string `yaml:"Keywords""`
	Content  string   `yaml:"Content"`
}

func (y *YAMLHelp) Describe() string {
	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("Keywords: %s", strings.Join(y.Keywords, ", ")))
	parts = append(parts, "")
	parts = append(parts, y.Content)
	return strings.Join(parts, "\n")
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

func Build(root string) map[string]YAMLHelp {
	index := make(map[string]YAMLHelp)

	nodes, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("Error reading help directory")
	}

	for _, f := range nodes {
		if f.IsDir() {
			continue
		}
		temp := buildFromPath(root + "/" + f.Name())
		for _, keyword := range temp.Keywords {
			index[keyword] = temp
		}
	}
	return index
}
