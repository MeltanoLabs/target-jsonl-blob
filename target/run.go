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
	"io"
	"log"
	"time"

	"gocloud.dev/blob"
)

type Target struct {
	Config Config
}

func handleNewStream(
	streamName string,
	ctx context.Context,
	config Config,
	bucket *blob.Bucket,
	streams map[string]StreamInfo,
	writers map[string]*blob.Writer,
) {
	now := time.Now()
	streams[streamName] = StreamInfo{
		StreamName:       streamName,
		Date:             now.Format("2006-01-02"),
		Minute:           now.Minute(),
		Hour:             now.Hour(),
		Day:              now.Day(),
		Month:            int(now.Month()),
		Year:             now.Year(),
		TimestampSeconds: now.Unix(),
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
	var message Message
	if err := json.Unmarshal(line, &message); err != nil {
		panic(err)
	}

	switch message.Type {
	case "RECORD":
		var recordMessage RecordMessage
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
		var schemaMessage SchemaMessage
		if err := json.Unmarshal(line, &schemaMessage); err != nil {
			panic(err)
		}
		log.Println("SCHEMA not implemented")
	case "STATE":
		var stateMessage StateMessage
		if err := json.Unmarshal(line, &stateMessage); err != nil {
			panic(err)
		}
		fmt.Println(string(stateMessage.Value))
	default:
		log.Printf("Unknown message type %s", message.Type)
	}
}

func (t Target) ProcessLines(r io.Reader) {
	scanner := bufio.NewScanner(r)
	streams := make(map[string]StreamInfo)
	writers := make(map[string]*blob.Writer)

	ctx := context.Background()
	bucket, err := blob.OpenBucket(ctx, t.Config.Bucket)
	if err != nil {
		log.Fatalf("Unable to open bucket, %v", err)
	}
	defer bucket.Close()

	for scanner.Scan() {
		processLine(scanner.Bytes(), ctx, t.Config, bucket, streams, writers)
	}

	for _, writer := range writers {
		closeErr := writer.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}
}
