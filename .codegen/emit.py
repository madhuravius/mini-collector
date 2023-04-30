#!/usr/bin/env python3
import os
import itertools
import json
import sys
import logging

from google.protobuf.compiler import plugin_pb2 as plugin

from jinja2 import Template


LOGGER = logging.getLogger(__name__)


def to_camel_case(s):
    bits = s.split('_')
    return "".join(c.title() for c in bits)


# REWRITE WHATEVER THIS IS
def generate_code(request, response):
    for proto_file in request.proto_file:
        output = []

        # Parse request
        for item in proto_file.message_type:
            if item.name != "PublishRequest":
                continue

            LOGGER.info("visit %s", item.name)

            for package, name in [
                ["writer/datadog", "datadog"],
                ["writer/influxdb", "influxdb"],
                ["publisher", "publisher"],
            ]:
                src = os.path.join(
                    os.path.dirname(__file__),
                    "{0}_formatter.go.jinja2".format(name)
                )

                f = response.file.add()
                f.name = os.path.join(
                    package,
                    "{0}.{1}_formatter.go".format(proto_file.name, name)
                )

                with open(src) as fh:
                    f.content = Template(fh.read()).render(
                        item=item,
                        to_camel_case=to_camel_case
                    )

def main():
    # See: https://www.expobrain.net/2015/09/13/create-a-plugin-for-google-protocol-buffer/

    # Read request message from stdin
    data = sys.stdin.buffer.read()

    # Parse request
    request = plugin.CodeGeneratorRequest()
    request.ParseFromString(data)

    # Create response
    response = plugin.CodeGeneratorResponse()

    # Generate code
    generate_code(request, response)

    # Serialise response message
    output = response.SerializeToString()

    # Write to stdout
    sys.stdout.buffer.write(output)


if __name__ == '__main__':
    main()
