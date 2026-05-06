package models

import "time"

type Entry struct {
	Value     any        `yaml:"value"`
	ExpiresAt *time.Time `yaml:"expires_at"`
}

type Data struct {
	Storage map[string]Entry `yaml:"storage"`
}
