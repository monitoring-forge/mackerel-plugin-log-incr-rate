# mackerel-plugin-log-incr-rate

Counts logs and measures the rate of log growth against base log

## Usage

```
$ mackerel-plugin-log-incr-rate -h
Usage:
  mackerel-plugin-log-incr-rate [OPTIONS]

Application Options:
      --log-file=      path to log file calcurate lines increased
      --base-log-file= path to base log file count lines
      --key-prefix=    Metric key prefix
  -v, --version        Show version

Help Options:
  -h, --help           Show this help message

```

Sample

```
$ mackerel-plugin-log-incr-rate --key-prefix err_per_access --log-file error_log --base-log-file access_log
log-incr-rate.err_per_access_count.log  438.986301      1571629417
log-incr-rate.err_per_access_count.base 454.438356      1571629417
log-incr-rate.err_per_access_rate.log   0.965997        1571629417
```

## Install

Please download release page or `mkr plugin install monitoring-forge/mackerel-plugin-log-incr-rate`.
