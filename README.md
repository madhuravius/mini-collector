Mini Collector
==============

Metrics collection and aggregation.


Building
--------

### Build Dependencies

You'll need a few build dependencies:

- Protobuf Compiler
- Protobuf Python support
- Protobuf Golang support

To find out how to install those, see `.travis.yml`.


### Building and Testing

Use `make build` and `make test`.

See `Makefile` for more detail.


### Building Docker images

To build Docker images, use:

```
make -f .docker/Makefile push APP=aggregator TAG=aggregator-vX.Y.Z
make -f .docker/Makefile push APP=mini-collector TAG=mini-collector-vX.Y.Z
```

See `.docker/Makefile` for more detail.
