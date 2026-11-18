# Copyright 2026 Haikei Labs
# SPDX-License-Identifier: Apache-2.0
#
# Python runner of the cross-language parity test (protocol: see
# scripts/parity/run.sh). It calls each operation through the generated
# kei.connectors.<connector>.v1.ops module by its method name, so a new
# operation or connector needs no change here. Needs protobuf and httpx.
import asyncio
import importlib
import json
import re
import sys

import httpx
from google.protobuf import descriptor_pool, json_format, message_factory

from kei.connection import http_connection


def snake(name: str) -> str:
    """GetIssue -> get_issue, as the Python template names methods."""
    return re.sub(r"(?<!^)(?=[A-Z])", "_", name).lower()


def queries_class(module):
    for attr in dir(module):
        if attr.endswith("Queries") and isinstance(getattr(module, attr), type):
            return getattr(module, attr)
    raise RuntimeError(f"{module.__name__} has no Queries class")


async def run(cases, base, token):
    async with httpx.AsyncClient() as client:
        for c in cases:
            line = {"op": c["op"]}
            try:
                module = importlib.import_module(f"kei.connectors.{c['connector']}.v1.ops")
                pool = descriptor_pool.Default()
                req_cls = message_factory.GetMessageClass(pool.FindMessageTypeByName(c["request_type"]))
                req = json_format.ParseDict(c["request"], req_cls())
                conn = http_connection(f"{base}/{c['op']}", "bearer", token, client)
                method = getattr(queries_class(module)(conn), snake(c["method"]))
                resp = await method(req)
                line["response"] = json_format.MessageToDict(resp)
            except Exception as e:  # reported per case; the test fails on it
                line["error"] = str(e)
            sys.stdout.write(json.dumps(line, ensure_ascii=False) + "\n")


def main():
    if len(sys.argv) != 4:
        sys.stderr.write("usage: runner.py <cases.json> <base url> <token>\n")
        sys.exit(2)
    with open(sys.argv[1], encoding="utf-8") as f:
        cases = json.load(f)
    asyncio.run(run(cases, sys.argv[2], sys.argv[3]))


if __name__ == "__main__":
    main()
