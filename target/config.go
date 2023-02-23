/*
Copyright © 2022 Edgar Ramírez-Mondragón

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package target

import (
	"bytes"
	"errors"
	"log"
	"text/template"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Bucket      string `mapstructure:"bucket"`
	KeyTemplate string `mapstructure:"key_template"`
}

type StreamInfo struct {
	StreamName       string
	Date             string
	TimestampSeconds int64
	Minute           int
	Hour             int
	Day              int
	Month            time.Month
	Year             int
}

const DEFAULT_KEY_TEMPLATE = "{{.StreamName}}.jsonl"

func ReadConfig(file string, c *Config) error {
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file %s", viper.ConfigFileUsed())
	} else {
		return err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return err
	}

	if c.Bucket == "" {
		return errors.New("bucket is required")
	}

	if c.KeyTemplate == "" {
		c.KeyTemplate = DEFAULT_KEY_TEMPLATE
	}

	return nil
}

func FillKeyTemplate(templateString string, info StreamInfo) (string, error) {
	keyTemplate := template.Must(template.New("streamKey").Parse(templateString))

	var keyBuffer bytes.Buffer
	err := keyTemplate.Execute(&keyBuffer, info)

	if err != nil {
		return "", err
	}

	return keyBuffer.String(), nil
}
