Mini Collector
==============

Metrics collection and aggregation. Relies very heavily on 
[cgroups](https://man7.org/linux/man-pages/man7/cgroups.7.html) to capture metrics
via `mini-collector` generated binary.

This comes in two flavors:

* `aggregator` - the thing that gets collected metrics to push to a destination in a 
golang-based channel/queue
* `mini-collector` - the thing that collects metrics on a given machine to push to an
aggregator

These two work in harmony via protobufs. 

Building
--------

This mainly runs on docker and inspects/collects data about Docker containers,
so you'll probably want that too.

### Build Dependencies

You'll need a few build dependencies:

- Protobuf Compiler
- Protobuf Golang support (requires Go as the bullet suggests)

To install this you can use brew on mac (or equivalent installation instructions
exist for apt and other linux-based package managers):

```sh
brew install protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Building and Testing

Use `make build` and `make test`. It is suggested to use the docker compose image for
local testing on MacOS to avoid headache.

### Building Docker images

To build Docker images, use:

```
make -f .docker/Makefile push APP=aggregator TAG=aggregator-vX.Y.Z
make -f .docker/Makefile push APP=mini-collector TAG=mini-collector-vX.Y.Z
```

### Running Locally

`docker compose` is suggested for local development

```sh
# build the image
docker compose build --no-cache

# need to run in a privileged context (to get access to docker itself, and
# the cgroups.) Currently this has been added to the compose spec, but not yet
# released: https://github.com/compose-spec/compose-spec/pull/292#event-8075878994
# which should allow us to avoid this in the future
# when released: docker compose run -it frontend /bin/bash

# for mini-collector
docker run \
  --privileged \
  --cgroupns=host \
  --user root \
  -it \
  -v .:/Users/$USER/work/mini-collector \
  -v /var/run/docker.sock:/var/run/docker.sock \
  mini-collector-frontend /bin/bash
```

To run the binaries locally, the following variables (each found in `cmd/*/main.go`) 
are clean dupes to run locally (and test against DataDog).  Note below are both on the same **host**:

#### Mini Collector

For `mini-collector` (Note `CONTAINER_ID` should be a valid docker container id, note this
will be the LONG name, can be extracted with: `docker inspect <SHORTNAME OR TAG> | jq '.[0].Id'`):

```sh
mkdir /test_data

export MINI_COLLECTOR_REMOTE_ADDRESS="host.docker.internal:8000"
export MINI_COLLECTOR_ENVIRONMENT_NAME="local"
export MINI_COLLECTOR_SERVICE_NAME="notreal"
export MINI_COLLECTOR_APP_NAME="aptible_test_app"
export MINI_COLLECTOR_MOUNT_PATH="/test_data"
export MINI_COLLECTOR_DEBUG=true
export MINI_COLLECTOR_CGROUP_PATH="/sys/fs/cgroup"

# USE A REAL ID
export MINI_COLLECTOR_CONTAINER_ID="$CONTAINER_ID"

/usr/local/go/bin/go run cmd/mini-collector/*.go
```

Also you can set `export MINI_COLLECTOR_POLL_INTERVAL=1s` if you want to develop against this with
a very short reporting time on collection.

#### Aggregator

For `aggregator`, this can be run on your machine!

```sh
export AGGREGATOR_DATADOG_CONFIGURATION='{"api_key": "$DATADOG_KEY"}'
export MINI_COLLECTOR_DEBUG=true

go run cmd/aggregator/*.go
```