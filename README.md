# target-jsonl-blob

JSONL Singer target for local storage, S3 and Azure Blob Storage.

## Configuration

| Setting | Required | Default | Description |
|----------------|----------|-------------------------|-------------------------------|
| `bucket` | Yes | - | Blob storage [bucket URL](#bucket-urls) |
| `key_template` | No | `{{.StreamName}}.jsonl` | Template string for file keys |

### Bucket URLs

| Storage | Example URL                 |
|---------|-----------------------------|
| local   | `file:///path/to/directory` |
| S3      | `s3://my-bucket`            |
| Azure   | `azblob://my-container`     |

### Available fields for `key_template`

- `StreamName`
- `Date` (YYYY-MM-DD)
- `TimestampSeconds`
- `Minute`
- `Hour`
- `Day`
- `Month`
- `Year`

Example: `{{.StreamName}}/{{.Year}}/{{.Month}}/{{.Day}}/{{.Hour}}/{{.Minute}}/{{.StreamName}}.jsonl`

## Build from source

```shell
go build -o target-jsonl-blob
```

## Usage with Meltano

_TODO_

## Roadmap

- Support GCS

  Currently blocked by

  ```
  cloud.google.com/go/storage@v1.16.1/storage.go:1416:53: o.GetCustomerEncryption().GetKeySha256 undefined (type *"google.golang.org/genproto/googleapis/storage/v2".Object_CustomerEncryption has no field or method GetKeySha256)
  ```

- Build a lighter binary
