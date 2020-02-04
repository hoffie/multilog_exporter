# multilog_exporter
Watches one or more configured log file paths, analyzes them based on configured patterns and exposes the result in the Prometheus metrics format.

The files are watched using `inotify` and can handle log file rotations (based on [grok_exporter](https://github.com/fstab/grok_exporter) code).

The patterns are regular expressions.
The pattern may contain groups which can later be used as label values or metric value.

Currently, the following actions are supported:

  * **add** an arbitrary number (Counter, Gauge)
  * **set** to a new value (Gauge)
  * **dec** to decrease by an arbitrary number (Gauge)

Supported values (for metric results and label values):

  * Static numbers, such as `1`, `10.3e5`
  * Regex group references, such as `$pool` from the regexp `pool=(?P<pool>[^ ]+) msg`
  * `now()` which returns the current Unix time

Metrics are registered at startup so that the relevant time series are generated even if no relevant log lines have been processed yet (only works for those without labels).

## Building
Tested on Linux 4.x using go 1.10: `go get -u github.com/hoffie/multilog_exporter`.
Run unit and integration tests using: `make`

## Configuration
`multilog_exporter` is mainly configured using a configuration file in yaml syntax.

The following options are command line flags, though:

  * `--metrics.listen-addr`, which specifies the ip and port for listening.
  * `--config.file`, the path to the mentioned config file
  * `--debug`, which helps diagnosing potential issues

The configuration file syntax is best described by the [annotated example configuration](doc/example.yaml).
In general, the config consists of a list of log file paths.
For each log file path, a list of patterns is configured.
The pattern describes how to map log lines to Prometheus metrics.

### Reloading configuration

Configuration can be reloaded without stopping the application, by means of a `SIGHUP` signal. This can be useful to
change the configuration (for instance, adding new log files or patterns) without losing current exposed metrics.

```bash
$ kill -SIGHUP $(pidof multilog_exporter)
```

When doing so, the application log should reflect the reload

```text
INFO[0137] SIGHUP received, reloading config file        configFile=doc/example.yaml
INFO[0137] Loading config file                           configFile=doc/example.yaml
```

## Limitations
### Multi-cardinality
There is no support for working with multi-cardinality in log files.
In other words, the log line `errors in the following pools: prod, test, dev` **cannot** generate multiple time series with the same name, such as: `errors{pool="prod"}`, `errors{pool="test"}`, `errors{pool="dev"}`.
If it turns out there is a use case for that, support might be added.

### Arithmetic
There is no support for any calculations or variables (except for references to regex groups).
If it turns out this is required and cannot be solved on the scraping side, support might be added.

## Similar tools
You may have noticed that there are two very popular tools which seem to have similar goals to this project: Google's [mtail](https://github.com/google/mtail) and @fstab's [grok_exporter](https://github.com/fstab/grok_exporter).
There is basically two reasons, why this tool exists:

  * multilog_exporter provides an easy way to monitor multiple log files while still handling them as independent files.
    For instance, you could use it to monitor your system log and count the frequency of certain events while concurrently monitoring an application log with totally different data.
  * multilog_exporter aims to provide a simple, integrated way to configure it, while the other tools implement own domain-specific languages or require knowledge of the grok format.

Both named tools seem to have a large user base and if their feature set matches your requirements, you are encouraged to use them.
In fact, multilog_exporter even uses grok_exporter's library for handling file system watches.

## License
This software is released under the [Apache 2.0 license](LICENSE).

## Author
multilog_exporter has been created by [Christian Hoffmann](https://hoffmann-christian.info/).
If you find this project useful, please star it or drop me a short [mail](mailto:mail@hoffmann-christian.info).
