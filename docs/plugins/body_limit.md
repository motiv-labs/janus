# Body Limit

Block incoming requests whose body is greater than a specific size in megabytes.

## Configuration

The plain request transformer config:

```json
"body_limit": {
    "enabled": true,
    "config": {
        "limit": "40M"
    }
}
```

Here is a simple definition of the available configurations.

| Configuration                 | Description                                                         |
|-------------------------------|---------------------------------------------------------------------|
| name                          | Name of the plugin to use, in this case: body_limit        |
| config.limit      | Allowed request payload size. You can set the size in `B` for bytes,`K` for kilobytes, `M` for megabytes, `G` for gigabytes and `T` for terabytes  |
