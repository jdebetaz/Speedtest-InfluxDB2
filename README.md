# ⚠️ WARNING ⚠️

This container is only compatible with InfluxDB v1.8 and v2.

# Speedtest.net Collector For InfluxDB2 and Grafana
[![CI to Docker Hub](https://github.com/jdebetaz/Speedtest-InfluxDB2/actions/workflows/main.yml/badge.svg)](https://github.com/jdebetaz/Speedtest-InfluxDB2/actions/workflows/main.yml)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/redkill2108/speedtest-influxdb2)](https://hub.docker.com/repository/docker/redkill2108/speedtest-influxdb2)

Runs Ookla's [Speedtest CLI](https://www.speedtest.net/apps/cli) program in Docker, sends the results to InfluxDB
  - Source code: [GitHub](https://github.com/jdebetaz/Speedtest-InfluxDB2)
  - Docker container: [Docker Hub](https://hub.docker.com/repository/docker/redkill2108/speedtest-influxdb2)
  - Image base: [Node](https://hub.docker.com/_/node)
  - Init system: N/A
  - Application: [Speedtest CLI](https://www.speedtest.net/apps/cli)
  - Architecture: `linux/amd64`

## Explanation

  - This runs Ooka's Speedtest CLI program on an interval, then writes the data to an InfluxDB database (you can later graph this data with Grafana or Chronograf)
  - This does **NOT** use the open-source [speedtest-cli](https://github.com/sivel/speedtest-cli). That program uses the Speedtest.net HTTP API. This program uses Ookla's official CLI application.
  - ⚠️ Ookla's speedtest application is closed-source (the binary applications are [here](https://www.speedtest.net/apps/cli)) and Ookla's reasoning for this decision is [here](https://www.reddit.com/r/HomeNetworking/comments/dpalqu/speedtestnet_just_launched_an_official_c_cli/f5tm9up/) ⚠️
  - ⚠️ Ookla's speedtest application reports all data back to Ookla ⚠️
  - ⚠️ This application uses Ookla's recommendation to install by piping curl to bash  ⚠️
  - The default output unit of measurement is **bytes-per-second**. You will most likely want to convert to megabits-per-second by dividing your output by 125000.

```
For example, if your download speed is 11702913 bytes-per-second:
11702913 / 125000 = 93.623304 megabits-per-second
```

## Requirements

  - This only works with InfluxDB v1.8 and v2, because I'm using [this](https://github.com/influxdata/influxdb-client-js) client library.
  - You must already have an InfluxDB database created, along with a user that has `WRITE` and `READ` permissions on that database.
  - This Docker container needs to be able to reach that InfluxDB instance by hostname, IP address, or Docker service name (I run this container on the same Docker network as my InfluxDB instance).
  - ⚠️ Depending on how often you run this, you may need to monitor your internet connection's usage. If you have a data cap, you could exceed it. The standard speedtest uses about 750MB of data per run. See below for an example. ⚠️

## Docker image information

### Docker image tags
  - `latest`: Latest version
  - `X.X.X`: [Semantic version](https://semver.org/) (use if you want to stick on a specific version)

### Environment variables
| Variable | Required? | Definition | Example | Comments |
|----------|-----------|------------|---------|----------|
| APP_INTERVAL | Yes | Minutes to sleep between runs | 60 |   
| INFLUX_HOST | Yes | Server hosting the InfluxDB | 'http://localhost:8086'  | |
| INFLUX_TOKEN | Yes | Token to connect to bucket | asdfghjkl | Needs to have the correct permissions. Setting this assumes we're talking to an InfluxDBv2 instance |
| INFLUX_ORG | Yes | Organization | my_test_org | |
| INFLUX_BUCKET | Yes | Database name | SpeedtestStats | Must already be created. In InfluxDBv2, this is the "bucket". |

