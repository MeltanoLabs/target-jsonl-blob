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
	"encoding/json"
	"testing"
)

func TestUnmarshalRecord(t *testing.T) {
	var message RecordMessage

	raw := `{
		"type": "RECORD",
		"stream": "test",
		"record": {"id": 3},
		"time_extracted": "2022-01-01T00:00:00Z"
	}`

	if err := json.Unmarshal([]byte(raw), &message); err != nil {
		panic(err)
	}

	if message.Type != "RECORD" {
		t.Errorf("Expected %s, got %s", "RECORD", message.Type)
	}

	if message.Stream != "test" {
		t.Errorf("Expected %s, got %s", "test", message.Stream)
	}

	if tsStr := message.TimeExtracted.String(); tsStr != "2022-01-01 00:00:00 +0000 UTC" {
		t.Errorf("Expected %s, got %s", "2022-01-01 00:00:00 +0000 UTC", tsStr)
	}
}
