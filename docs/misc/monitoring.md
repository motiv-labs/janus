### Monitoring

`Janus` monitoring is built on top of [`hellofresh/stats-go`](https://github.com/hellofresh/stats-go) library.
You can configure it with the following env variables:

* `STATS_DSN` (default `log://`) - DSN of stats backend
* `STATS_IDS` - second level ID list for URLs to generalise metric names, see details in [Generalise resources by type and stripping resource ID](https://github.com/hellofresh/stats-go#generalise-resources-by-type-and-stripping-resource-id)
* `STATS_AUTO_DISCOVER_THRESHOLD` - threshold for second level IDs autodiscovery, see details in [Generalise resources by type and stripping resource ID](https://github.com/hellofresh/stats-go#generalise-resources-by-type-and-stripping-resource-id)
* `STATS_AUTO_DISCOVER_WHITE_LIST` - white list for second level IDs autodiscovery, see details in [Generalise resources by type and stripping resource ID](https://github.com/hellofresh/stats-go#generalise-resources-by-type-and-stripping-resource-id)
* `STATS_ERRORS_SECTION` (default `error-log`) - section for error logs monitoring, see details in [Usage for error logs monitoring](https://github.com/hellofresh/stats-go#usage-for-error-logs-monitoring-using-githubcomsirupsenlogrus)
