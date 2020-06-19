# Kiosk

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![build](https://travis-ci.org/jibitters/kiosk.svg?branch=master)](https://travis-ci.org/jibitters/kiosk)

A typical ticketing system that is designed for micro services environments with highly scalable and concurrent needs.

## How it works
Kiosk designed by scalability and ease of use in mind, so we chose nats as message bus because of its modern distributed
patterns. Every kiosk node listens to all subjects but in queue grouped manner, so the requests will distribute between
different nodes. The message protocol is typical JSON format, so it can be used by all nats clients.

For more information about subject names and request/response models see Wiki pages.

## How to test and build

The requirements to test and build the project are as follows:

|Requirement                                                                   |Version|
|---                                                                           |---    |
|go                                                                            |1.14   |
|postgres                                                                      |11     |
|nats                                                                          |2.1    |

To prepare your environment for next steps first run:

`./scripts/setup.sh`

To build an executable instance of the project, use:

`./scripts/build.sh`

To run tests you can use the test.sh script (Ensure Docker is up and running):

`./scripts/test.sh`

To build a docker image (Images also available on [Docker Hub](https://hub.docker.com/r/jibitters/kiosk))

`docker build -t image:tag .`

## How to run

`./kiosk-linux-[version] --config path/to/kiosk.json` starts the project, easily!

See `configs/kiosk.json` for an example configuration.
