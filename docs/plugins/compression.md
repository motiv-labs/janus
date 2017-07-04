# Compression

Enables gzip compression if the client supports it. By default, responses are not gzipped. If enabled, the default settings will ensure that images, videos, and archives (already compressed) are not gzipped.

The plain compression config is good enough for most things, but you can gain more control if needed:

```json
"compression": {
    "enabled": true
}
```
