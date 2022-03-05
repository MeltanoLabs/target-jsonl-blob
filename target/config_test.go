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
	"testing"
)

func TestReadConfig(t *testing.T) {
	var config Config

	err := ReadConfig("testdata/config.good.example.json", &config)
	if err != nil {
		t.Error(err)
	}

	if config.Bucket != "s3://my-bucket" {
		t.Errorf("Expected %s, got %s", "s3://my-bucket", config.Bucket)
	}

	if config.KeyTemplate != DEFAULT_KEY_TEMPLATE {
		t.Errorf("Expected %s, got %s", DEFAULT_KEY_TEMPLATE, config.KeyTemplate)
	}
}

func TestReadConfigMissingBucket(t *testing.T) {
	var config Config

	err := ReadConfig("testdata/config.no-bucket.example.json", &config)
	if err == nil || err.Error() != "bucket is required" {
		t.Errorf("Expected %s, got %s", "bucket is required", err.Error())
	}
}

func TestFillKeyTemplate(t *testing.T) {
	config := Config{KeyTemplate: "prefix/{{.StreamName}}-123.jsonl", Bucket: "my-bucket"}
	streamInfo := StreamInfo{StreamName: "my-stream"}

	key, err := FillKeyTemplate(config.KeyTemplate, streamInfo)
	if err != nil {
		t.Error(err)
	}

	expected := "prefix/my-stream-123.jsonl"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestFillKeyBadTemplate(t *testing.T) {
	config := Config{KeyTemplate: "{{.StreamName}}-{{.Missing}}-123.jsonl", Bucket: "my-bucket"}
	streamInfo := StreamInfo{StreamName: "my-stream"}

	_, err := FillKeyTemplate(config.KeyTemplate, streamInfo)
	if err == nil {
		t.Error("Expected a failure")
	}

	expectedMessage := `template: streamKey:1:18: executing "streamKey" at <.Missing>: can't evaluate field Missing in type target.StreamInfo`
	if err.Error() != expectedMessage {
		t.Errorf("Expected %s, got %s", expectedMessage, err.Error())
	}
}
