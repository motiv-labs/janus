# Response Transformer

Transform the response sent by a client on the fly on Janus, before giving it back to the client.

## Configuration

The plain response transformer config:

```json
"response_transformer": {
    "enabled": true,
    "config": {
        "add": {
            "headers": {
                "X-Something": "Value"
            }
        },
        "append": {
            "headers": {
                "X-Something-More": "Value"
            }
        },
        "replace": {
            "headers": {
                "X-Something": "New Value"
            }
        },
        "remove": {
            "headers": {
                "X-Something": ""
            }
        }
    }
}
```


Here is a simple definition of the available configurations.

| Configuration                 | Description                                                         |
|-------------------------------|---------------------------------------------------------------------|
| name                          | Name of the plugin to use, in this case: response_transformer        |
| config.remove.headers         | List of header names. Unset the headers with the given name.        |
| config.replace.headers        | List of headername:value pairs. If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set.        |
| config.add.headers            | List of headername:value pairs. If and only if the header is not already set, set a new header with the given value. Ignored if the header is already set.        |
| config.append.headers         | List of headername:value pairs. If the header is not set, set it with the given value. If it is already set, a new header with the same name and the new value will be set.        |

## Order of execution

Plugin performs the response transformation in following order

`remove --> replace --> add --> append`
