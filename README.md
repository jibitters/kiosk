# Kiosk

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![build](https://travis-ci.org/jibitters/kiosk.svg?branch=master)](https://travis-ci.org/jibitters/kiosk)
[![codecov](https://codecov.io/gh/jibitters/kiosk/branch/master/graph/badge.svg)](https://codecov.io/gh/jibitters/kiosk)
[![release](https://img.shields.io/badge/release-v0.0.5-1eb0fc.svg)](https://github.com/jibitters/kiosk/releases/tag/v0.0.5)

## About

A typical ticketing system that provides both gRPC and REST interfaces. This project is intended to be used by internal microservices so we recommend to not expose interfaces directly to the public network.

## How to run

If you want to setup kiosk in your environment please read this section; otherwise you can skip forward to the next topic.

You can download the latest stable release from [releases](https://github.com/jibitters/kiosk/releases) page and run it with --help to see the configuration requirements.

`./kiosk-linux-v1 --config path/to/kiosk.json` starts the project, easily! See configs/kiosk.json for an example configuration.

## How to build

The requirements to test and build the project are as follow:

|Requirement                           |Version|
|---                                   |---    |
|protoc                                |3.9.2  |
|golang/protobuf                       |1.3.2  |
|postgres                              |11.6   |
|nats                                  |2.1.2  |

Optional requirements are as follow:

|Requirement                           |Version|
|---                                   |---    |
|prometheus                            |latest |
|grafana                               |latest |

To build an executable instance of the project, use:

`./scripts/build.sh`

Also to run tests you can use the test.sh script (Note: Ensure docker is up and running):

`./scripts/test.sh`

And to build a docker image:

`docker build -t image:tag .`

## Interface Documentation

Kiosk provides two different types of interfaces, one is based on gRPC as protobuf definitions and the other is REST API.

|Type                                                      |
|---                                                       |
|[Protobuf Definitions](api/protobuf-spec)                 |
|[REST API Specification](api/rest-spec)                   |
