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
make -f .docker/Makefile push TAG=aggregator
make -f .docker/Makefile push TAG=mini-collector
```

See `.docker/Makefile` for more detail.
