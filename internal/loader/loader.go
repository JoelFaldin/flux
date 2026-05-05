package loader

import (
	"fmt"
	"maps"
	"os"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Storage map[string]string `yaml:"storage"`
}

func LoadData() *Data {
	file, err := os.ReadFile("data.yaml")
	if err != nil {
		fmt.Println("Error opening yaml file:", err)
		os.Exit(1)
	}

	var data Data

	err = yaml.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error processing yaml data:", err)
		os.Exit(1)
	}

	return &data
}

func WriteData(entry Data, values map[string]string) {
	maps.Copy(entry.Storage, values)

	yamlData, err := yaml.Marshal(&entry)
	if err != nil {
		fmt.Println("Error marshaling yaml data:", err)
		os.Exit(1)
	}

	err = os.WriteFile("data.yaml", yamlData, 0644)
	if err != nil {
		fmt.Println("Error writing yaml file:", err)
		os.Exit(1)
	}
}
