package main

import (
	_ "fmt"
)

type configuration struct {
	ConfigItems []botConfig `json:"config_items,omitempty"`
}

func (c *configuration) Validate() error {
	if c == nil {
		return nil
	}

	items := c.ConfigItems
	for i := range items {
		if err := items[i].validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *configuration) SetDefault() {
	if c == nil {
		return
	}

	Items := c.ConfigItems
	for i := range Items {
		Items[i].setDefault()
	}
}

type botConfig struct {

	// CacheFileSourceUrl is the url of code hosting platform repository that need to cache something
	SigProjectPath string `yaml:"sig-project-path"`
}

func (c *botConfig) setDefault() {
	// do noting
}

func (c *botConfig) validate() error {
	// do noting
	return nil
}
