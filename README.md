# PoolCounter Prometheus exporter
A Prometheus exporter for the [poolcounter](https://www.mediawiki.org/wiki/PoolCounter) daemon.

## Building 🛠
After [setting up](https://golang.org/doc/install) your Go development environment, you can run `go build`
to create an executable. Alternatively you may use the provided Dockerfile to build a container image.

## Configuration 📋
Configuration is done via environment variables. All settings are optional and have sane defaults:
* `EXPORTER_POOL_COUNTER_ADDRESS` - host:port of the poolcounter instance to gather metrics from. Default: `localhost:7531`.
* `EXPORTER_LISTEN_ADDRESS` - host:port the collector should listen on. Default: `localhost:8000`.
* `EXPORTER_LOGS_AS_JSON` - whether to format stdout logs as JSON or as human readable output. Default: `false`.
* `EXPORTER_COLLECTOR_TIMEOUT_SECONDS` - TCP timeout value used by the metrics collector, in seconds. Default: `5.
* `EXPORTER_SERVER_TIMEOUT_SECONDS` - HTTP timeout values used by the server, in seconds. Default: `3`.

## Available metrics 📡
Metrics are made available at the `/prometheus` HTTP endpoint. They correspond to the [internal metrics](https://www.mediawiki.org/wiki/PoolCounter#Testing) tracked by poolcounter.
### Counters 💯
* `poolcounter_total_processing_time_seconds` - total seconds spent by workers on performing poolcounter-protected tasks
* `poolcounter_total_gained_time_seconds` - total seconds of processing time saved by the use of PoolCounter in seconds
* `poolcounter_total_excl_wait_time_seconds` - total seconds spent by workers on waiting on exclusive locks
* `poolcounter_total_shared_wait_time_seconds` - total seconds spent by workers on waiting on shared locks
* `poolcounter_total_acquired` - total acquired locks count
* `poolcounter_total_releases` - total released locks count
* `poolcounter_connect_errors` - total number of client connection errors
* `poolcounter_lock_mismatch` - total number of mismatched locks
* `poolcounter_release_mismatch` - total number of received `RELEASE` commands for which no lock was found
* `poolcounter_processed_count` - total number of tasks processed
### Gauges ⏲
* `poolcounter_hashtable_entries` - number of entries in poolcounter hash table
* `poolcounter_processing_workers` - number of workers busy performing poolcounter-protected tasks
* `poolcounter_waiting_workers` - number of workers waiting in the queue
* `poolcounter_full_queues` - number of queues that are full of waiting workers

## Logging
Logs are sent to standard output, either in human readable form or as JSON depending on the value of
the `EXPORTER_LOGS_AS_JSON` environment variable.
