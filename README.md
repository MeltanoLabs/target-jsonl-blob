# target-jsonl-blob

JSONL Singer target for local storage, S3 and Azure Blob Storage.

## Installation

To install this Singer tap, you can download a [prebuilt binary](https://github.com/MeltanoLabs/target-jsonl-blob/releases), or you can [build it from source](#build-from-source).

## Configuration

| Setting | Required | Default | Description |
|----------------|----------|-------------------------|-------------------------------|
| `bucket` | Yes | - | Blob storage [bucket URL](#bucket-urls) |
| `key_template` | No | `{{.StreamName}}.jsonl` | Template string for file keys |

### Bucket URLs

| Storage | Example URL                 | Supported URL parameters                                                            |
|---------|-----------------------------|-------------------------------------------------------------------------------------|
| local   | `file:///path/to/directory` | See [supported parameters](https://pkg.go.dev/gocloud.dev/blob/fileblob#URLOpener)  |
| S3      | `s3://my-bucket`            | See [supported parameters](https://pkg.go.dev/gocloud.dev/blob/s3blob#URLOpener)    |
| Azure   | `azblob://my-container`     | See [supported parameters](https://pkg.go.dev/gocloud.dev/blob/azureblob#URLOpener) |
| GCS     | `gs://my-bucket`            | See [supported parameters](https://pkg.go.dev/gocloud.dev/blob/gcsblob#URLOpener)   |

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
