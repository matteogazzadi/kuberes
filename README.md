# `kuberes`: Kubernetes Resources Analysis Tool

![Latest GitHub release](https://img.shields.io/github/release/matteogazzadi/kuberes.svg)
![Build](https://github.com/matteogazzadi/kuberes/actions/workflows/build.yml/badge.svg)

This repository provides `kuberes` tool.

## What is `kuberes`?

`kuberes` is a tool to summarize the configured  resources `requests` and `limits` for Kubernetes PODs in a given cluster. 

It allows to produce reports in the following formats:
 - `table`: directly on screen formated as a table
 - `csv`: in a .CSV file, comma separated
 - `xlsx`: in a .XLSX file

Here is a demo of `kuberes`:
![kuberes demo](docs/kuberes-demo.gif)

## Arguments

`kuberes` accepts multiple arguments to define output format or filtering.

Arguments can be passed using `-` or `--`.

|      Argument    |  Type  |  Default |                    Description                     |
| ---------------- | ------ | -------- | -------------------------------------------------- |
| `output`         | String | `table`  |  Output type. Valid values are: table,csv          |
| `group-by-ns`    |  Bool  | `true`   |  Should group statistics by namespace ?            |
| `csv-path`       | String | `""`     |  Full Path to the .CSV File to produce             |
| `exclude-ns`     | String | `""`     |  Namespaces names to be ignored, comma separated   |
| `match-ns-regex` | String | `""`     |  Namespaces Names to be matched on the given RegEx |

## Examples

```sh
# List Resources grouped by Namespace in table format
.\kuberes

# List Resources by POD in table format
.\kuberes --group-by-ns=false

# List Resource grouped by Namespace and save in a .CSV file
.\kuberes --output=csv --csv-path="output.csv"

# List Resources grouped by Namespace in table format, exclude namespaces: default and kube-system
.\kuberes --exclude-ns="default,kube-system"

# List Resources grouped by Namespace in table format, include namespace containing word "test" only
.\kuberes --math-ns-regex="^.*test.*$"

```