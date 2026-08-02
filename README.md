# mackerel-plugin-log-incr-rate

2つのログファイルの増加率（増加割合）を Mackerel のメトリクスとして出力するプラグインです。

## 概要

このプラグインは、対象のログファイル（`--log-file`）と基準となるログファイル（`--base-log-file`）の行数増加量を同時にカウントし、対象ログの増加量が基準ログの増加量に対してどの程度の比率になっているかを計測します。

例えば、以下のような用途で使うことができます。

- アクセスログに対するエラーログの比率を監視する
- 基準ログに対して特定のログがどれだけ発生しているかをトレンドとして把握する

出力されるメトリクスは次の3つです。

| メトリクス名 | 意味 |
| --- | --- |
| `log-incr-rate.<prefix>_count.log` | 対象ログの1秒あたり平均増加行数 |
| `log-incr-rate.<prefix>_count.base` | 基準ログの1秒あたり平均増加行数 |
| `log-incr-rate.<prefix>_rate.log` | 対象ログの増加量を基準ログの増加量で割った比率 |

## インストール

リリースページからバイナリをダウンロードするか、以下の `mkr` コマンドでインストールしてください。

```
$ mkr plugin install monitoring-forge/mackerel-plugin-log-incr-rate
```

## 使い方

```
$ mackerel-plugin-log-incr-rate -h
Usage:
  mackerel-plugin-log-incr-rate [OPTIONS]

Application Options:
      --log-file=      path to log file calculate lines increased
      --base-log-file= path to base log file count lines
      --key-prefix=    Metric key prefix
  -v, --version        Show version

Help Options:
  -h, --help           Show this help message

```

### オプション

| オプション | 必須 | 説明 |
| --- | --- | --- |
| `--log-file` | 必須 | 比率の分子となる対象ログファイルのパス |
| `--base-log-file` | 必須 | 比率の分母となる基準ログファイルのパス |
| `--key-prefix` | 必須 | メトリクス名のプレフィックス |
| `--verbose` | 任意 | 詳細なログを出力する |

### 実行例

エラーログ（`error_log`）をアクセスログ（`access_log`）に対する比率で監視する例です。

```
$ mackerel-plugin-log-incr-rate --key-prefix err_per_access --log-file error_log --base-log-file access_log
log-incr-rate.err_per_access_count.log  438.986301      1571629417
log-incr-rate.err_per_access_count.base 454.438356      1571629417
log-incr-rate.err_per_access_rate.log   0.965997        1571629417
```

出力結果の意味は次の通りです。

- `log-incr-rate.err_per_access_count.log`: 対象ログ（`error_log`）の1秒あたり平均行数
- `log-incr-rate.err_per_access_count.base`: 基準ログ（`access_log`）の1秒あたり平均行数
- `log-incr-rate.err_per_access_rate.log`: `error_log / access_log` の増加比率

## 計算方法

各ログファイルは前回実行時からの追加分のみをカウントします。対象ログと基準ログの1秒あたり平均増加行数をそれぞれ求め、以下の式で比率を計算します。

```
rate = (対象ログの1秒あたり平均増加行数) / (基準ログの1秒あたり平均増加行数)
```

基準ログの増加行数が 0 の場合、`rate` は出力されません。
