package loader

import (
	"flux/internal/models"
	"fmt"
	"maps"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadData() *models.Data {
	file, err := os.ReadFile("data.yaml")
	if err != nil {
		fmt.Println("Error opening yaml file:", err)
		os.Exit(1)
	}

	var data models.Data

	err = yaml.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error processing yaml data:", err)
		os.Exit(1)
	}

	return &data
}

func WriteData(cfg *models.Data, values map[string]models.Entry) {
	if cfg.Storage == nil {
		cfg.Storage = make(map[string]models.Entry)
	}

	maps.Copy(cfg.Storage, values)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Println("Error marshaling yaml data:", err)
		os.Exit(1)
	}

	err = os.WriteFile("data.yaml", data, 0644)
	if err != nil {
		fmt.Println("Error writing yaml file:", err)
		os.Exit(1)
	}
}
