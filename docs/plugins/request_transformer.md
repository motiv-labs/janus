# Request Transformer

Transform the request sent by a client on the fly on Janus, before hitting the upstream server.

## Configuration

> To enable this plugin to your API you should configure an [OAuth Server](auth/oauth.md) first!

Here is a simple definition of the available configurations.

| Configuration                 | Description                                                         |
|-------------------------------|---------------------------------------------------------------------|
| name                          | Name of the plugin to use, in this case: request-transformer        |
| config.remove.headers         | List of header names. Unset the headers with the given name.        |
| config.remove.querystring     | List of querystring names. Remove the querystring if it is present. |
| config.remove.body            | List of parameter names. Remove the parameter if and only if content-type is one the following  |
| config.replace.headers        | List of headername:value pairs. If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set.        |
| config.replace.querystring    | List of queryname:value pairs. If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set. |
| config.replace.body           | List of paramname:value pairs. If and only if content-type is one the following [application/json, multipart/form-data, application/x-www-form-urlencoded] and the parameter is already present, replace its old value with the new one. Ignored if the parameter is not already present.  |
| config.add.headers            | List of headername:value pairs. If and only if the header is not already set, set a new header with the given value. Ignored if the header is already set.        |
| config.add.querystring        | List of queryname:value pairs. If and only if the querystring is not already set, set a new querystring with the given value. Ignored if the header is already set. |
| config.add.body               | List of pramname:value pairs. If and only if content-type is one the following [application/json, multipart/form-data, application/x-www-form-urlencoded] and the parameter is not present, add a new parameter with the given value to form-encoded body. Ignored if the parameter is already present.  |
| config.append.headers         | List of headername:value pairs. If the header is not set, set it with the given value. If it is already set, a new header with the same name and the new value will be set.        |
| config.append.querystring     | 	List of queryname:value pairs. If the querystring is not set, set it with the given value. If it is already set, a new querystring with the same name and the new value will be set. |
| config.append.body            | List of paramname:value pairs. If the content-type is one the following [application/json, application/x-www-form-urlencoded], add a new parameter with the given value if the parameter is not present, otherwise if it is already present, the two values (old and new) will be aggregated in an array.  |

## Order of execution

Plugin performs the response transformation in following order

`remove --> replace --> add --> append`
