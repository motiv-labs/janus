# Request Transformer

Transform the request sent by a client on the fly on Janus, before hitting the upstream server.

## Configuration

The plain request transformer config:

```json
"request_transformer": {
    "enabled": true,
    "config": {
        "add": {
            "headers": {
                "X-Something": "Value"
            },
            "querystring": {
                "test": "value"
            }
        },
        "append": {
            "headers": {
                "X-Something-More": "Value"
            },
            "querystring": {
                "extra": "stuff"
            }
        },
        "replace": {
            "headers": {
                "X-Something": "New Value"
            },
            "querystring": {
                "test": "new value"
            }
        },
        "remove": {
            "headers": {
                "X-Something": ""
            },
            "querystring": {
                "test": ""
            }
        }
    }
}
```


Here is a simple definition of the available configurations.

| Configuration                 | Description                                                         |
|-------------------------------|---------------------------------------------------------------------|
| name                          | Name of the plugin to use, in this case: request-transformer        |
| config.remove.headers         | List of header names. Unset the headers with the given name.        |
| config.remove.querystring     | List of querystring names. Remove the querystring if it is present. |
| config.replace.headers        | List of headername:value pairs. If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set.        |
| config.replace.querystring    | List of queryname:value pairs. If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set. |
| config.add.headers            | List of headername:value pairs. If and only if the header is not already set, set a new header with the given value. Ignored if the header is already set.        |
| config.add.querystring        | List of queryname:value pairs. If and only if the querystring is not already set, set a new querystring with the given value. Ignored if the header is already set. |
| config.append.headers         | List of headername:value pairs. If the header is not set, set it with the given value. If it is already set, a new header with the same name and the new value will be set.        |
| config.append.querystring     | 	List of queryname:value pairs. If the querystring is not set, set it with the given value. If it is already set, a new querystring with the same name and the new value will be set. |

## Order of execution

Plugin performs the response transformation in following order

`remove --> replace --> add --> append`
