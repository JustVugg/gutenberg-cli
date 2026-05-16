# Go Setup

Gutenberg uses Go as the primary generated tool target.

For this workspace, Go is installed outside the Windows-mounted tree for speed:

```text
/tmp/black-forge-go/go/bin/go
```

The project wrapper is:

```bash
scripts/use-go.sh version
```

Use it for generated tools:

```bash
scripts/use-go.sh test ./...
scripts/use-go.sh build ./cmd/petstore
```

The downloaded archive is kept under:

```text
.tools/downloads/go1.26.3.linux-amd64.tar.gz
```

Official checksum used:

```text
2b2cfc7148493da5e73981bffbf3353af381d5f93e789c82c79aff64962eb556
```
