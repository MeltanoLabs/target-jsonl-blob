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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gocloud.dev/blob"
)

type singerMessage struct {
	Type string `json:"type"`
}

type singerSchema struct {
	Stream             string          `json:"stream"`
	Schema             json.RawMessage `json:"schema"`
	KeyProperties      []string        `json:"key_properties"`
	BookmarkProperties []string        `json:"bookmark_properties"`
	singerMessage
}

type singerRecord struct {
	Stream        string          `json:"stream"`
	Record        json.RawMessage `json:"record"`
	TimeExtracted *time.Time      `json:"time_extracted"`
	singerMessage
}

type singerState struct {
	singerMessage
	Value json.RawMessage `json:"value"`
}

func handleNewStream(
	streamName string,
	ctx context.Context,
	config Config,
	bucket *blob.Bucket,
	streams map[string]StreamInfo,
	writers map[string]*blob.Writer,
) {
	streams[streamName] = StreamInfo{
		StreamName:    streamName,
		SyncStartTime: time.Now(),
	}

	objectKey, err := FillKeyTemplate(config.KeyTemplate, streams[streamName])
	if err != nil {
		log.Fatal(err)
	}

	w, err := bucket.NewWriter(ctx, objectKey, nil)
	if err != nil {
		log.Fatalf("Unable to write to bucket, %v", err)
	}
	writers[streamName] = w
}

func processLine(
	line []byte,
	ctx context.Context,
	config Config,
	bucket *blob.Bucket,
	streams map[string]StreamInfo,
	writers map[string]*blob.Writer,
) {
	var message singerMessage
	if err := json.Unmarshal(line, &message); err != nil {
		panic(err)
	}

	switch message.Type {
	case "RECORD":
		var recordMessage singerRecord
		if err := json.Unmarshal(line, &recordMessage); err != nil {
			panic(err)
		}

		_, ok := streams[recordMessage.Stream]
		if !ok {
			handleNewStream(recordMessage.Stream, ctx, config, bucket, streams, writers)
		}

		w := writers[recordMessage.Stream]
		_, writeErr := fmt.Fprintln(w, string(recordMessage.Record))
		if writeErr != nil {
			log.Fatal(writeErr)
		}
	case "SCHEMA":
		var schemaMessage singerSchema
		if err := json.Unmarshal(line, &schemaMessage); err != nil {
			panic(err)
		}
		log.Println("SCHEMA not implemented")
	case "STATE":
		var stateMessage singerState
		if err := json.Unmarshal(line, &stateMessage); err != nil {
			panic(err)
		}
		fmt.Println(string(stateMessage.Value))
	default:
		log.Printf("Unknown message type %s", message.Type)
	}
}

func ProcessLines(
	config Config,
	ctx context.Context,
	bucket *blob.Bucket,
	writers map[string]*blob.Writer,
) {
	scanner := bufio.NewScanner(os.Stdin)
	streams := make(map[string]StreamInfo)

	for scanner.Scan() {
		processLine(scanner.Bytes(), ctx, config, bucket, streams, writers)
	}
}
