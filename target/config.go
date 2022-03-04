package target

import (
	"bytes"
	"text/template"
	"time"
)

type Config struct {
	Bucket      string `mapstructure:"bucket"`
	KeyTemplate string `mapstructure:"key_template"`
}

type StreamInfo struct {
	StreamName    string
	SyncStartTime time.Time
}

const DEFAULT_KEY_TEMPLATE = "{{.StreamName}}.jsonl"

func FillKeyTemplate(templateString string, info StreamInfo) (string, error) {
	keyTemplate := template.Must(template.New("streamKey").Parse(templateString))

	var keyBuffer bytes.Buffer
	err := keyTemplate.Execute(&keyBuffer, info)

	if err != nil {
		return "", err
	}

	return keyBuffer.String(), nil
}
