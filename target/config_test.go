package target

import (
	"testing"
)

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
