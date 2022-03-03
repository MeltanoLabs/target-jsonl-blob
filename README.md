# target-jsonl-blob

JSONL Singer target for local storage, S3, GCS and Azure Blob Storage.

## Immediate Roadmap

- Build a working singer target that can write to the destinations mentioned above
- Make the destination easily configurable with a "bucket" URL (e.g. `"s3://my-bucket?region=us-west-1"`)
- Use a configurable template for the object keys with a rich set of available inputs (e.g. `prefix/{stream_name}/{day}/{hour}.jsonl`)

## Build from source

_TODO_

## Usage with Meltano

_TODO_
