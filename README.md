# `kuberes`: Kubernetes Resources Analysis Tool

![Latest GitHub release](https://img.shields.io/github/release/matteogazzadi/kuberes.svg)
![Build](https://github.com/matteogazzadi/kuberes/actions/workflows/build.yml/badge.svg)

This repository provides `kuberes` tool.

## What is `kuberes`?

`kuberes` is a tool to summarize the configured  resources `requests` and `limits` for Kubereetes POD in a given cluster. 

Here is a demo of `kuberes`:

## Arguments

`kuberes` accept multiple arguments allowing to produce report directly in console or to a `.csv` file.

Arguments can be passed usin `-` or `--`.

|    Argument   |  Type  |  Default |                Description                |
| ------------- | ------ | -------- | ----------------------------------------- |
| `output`      | String | `table`  |  Output type. Valid values are: table,csv |
| `group-by-ns` |  Bool  | `true`   |  Should group statistics by namespace ?   |
| `csv-path`    | String | ``       |  Full Path to the .CSV File to produce    |
 