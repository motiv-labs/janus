# CORS

Easily add Cross-origin resource sharing (CORS) to your API by enabling this plugin.

## Configuration

| Configuration   | Description                                                                                                                                                                  |
|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| domains         | A comma-separated list of allowed domains for the Access-Control-Allow-Origin header. If you wish to allow all origins, add * as a single value to this configuration field. |
| methods         | Value for the Access-Control-Allow-Methods header, expects a comma delimited string (e.g. GET,POST).                                                                         |
| request_headers | Value for the Access-Control-Allow-Headers header, expects a comma delimited string (e.g. Origin, Authorization).                                                            |
| exposed_headers | Value for the Access-Control-Expose-Headers header, expects a comma delimited string (e.g. Origin, Authorization). If not specified, no custom headers are exposed.          |
