package loader

import (
	"fmt"
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
