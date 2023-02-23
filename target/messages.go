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
	"time"
)

type Message struct {
	Type string `json:"type"`
}

type SchemaMessage struct {
	Message
	Stream             string          `json:"stream"`
	Schema             json.RawMessage `json:"schema"`
	KeyProperties      []string        `json:"key_properties"`
	BookmarkProperties []string        `json:"bookmark_properties"`
}

type RecordMessage struct {
	Message
	Stream        string          `json:"stream"`
	Record        json.RawMessage `json:"record"`
	TimeExtracted time.Time       `json:"time_extracted"`
}

type StateMessage struct {
	Message
	Value json.RawMessage `json:"value"`
}
