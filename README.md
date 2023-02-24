# target-jsonl-blob

JSONL Singer target for local storage, S3, GCS and Azure Blob Storage.

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

1. [Download the appropriate asset](#1-download-the-appropriate-asset)
1. [Add a custom Meltano plugin to your project](#3-add-a-custom-meltano-plugin-to-your-project)
1. [Run a pipeline](#4-run-a-pipeline)

### 1. Download the appropriate asset

You can see the full list of assets in the release page: https://github.com/MeltanoLabs/target-jsonl-blob/releases/latest.

The [`gh`](https://cli.github.com/) tool makes downloading an asset easy:

```bash
gh release download v0.0.6 \
  -R MeltanoLabs/target-jsonl-blob \
  -p '*darwin_amd64*' \
  --clobber

tar -xvf target-jsonl-blob_0.0.6_darwin_amd64.tar.gz target-jsonl-blob
```

You can also install with Homebrew:

```bash
brew tap MeltanoLabs/target-jsonl-blob https://github.com/MeltanoLabs/target-jsonl-blob
brew install target-jsonl-blob
```

### 2. Add a custom Meltano plugin to your project

```yaml
# meltano.yml
plugins:
  loaders:
  - name: target-jsonl-blob
    namespace: target_jsonl_blob
    executable: ./target-jsonl-blob
    settings:
    - name: bucket
      label: Bucket
      description: Target directory (local, S3, Azure Blob)
    - name: key_template
      label: Key Template
      description: Template string for file keys
    config:
      bucket: file://./output/my-bucket
      key_template: $MELTANO_EXTRACTOR_NAMESPACE/{{.StreamName}}.jsonl
```

You also need to ensure the local "bucket" exists:

```bash
mkdir output/my-bucket
```

### 3. Run a pipeline

```bash
meltano run tap-github target-jsonl-blob
```
