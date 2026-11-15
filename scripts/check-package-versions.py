#!/usr/bin/env python3
"""Check that generated package manifests share a release version."""
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def toml_version(path: str, section: str) -> str:
    text = (ROOT / path).read_text(encoding="utf-8")
    match = re.search(rf"(?ms)^\[{section}\]\s*$(.*?)(?=^\[|\Z)", text)
    if match is None:
        raise ValueError(f"missing [{section}] in {path}")
    version = re.search(r'^version\s*=\s*"([^"\n]+)"\s*$', match.group(1), re.MULTILINE)
    if version is None:
        raise ValueError(f"missing version in {path}")
    return version.group(1)


try:
    versions = {
        "python": toml_version("gen/python/pyproject.toml", "project"),
        "typescript": json.loads((ROOT / "gen/ts/package.json").read_text())["version"],
        "rust": toml_version("gen/rust/Cargo.toml", "package"),
    }
    if len(set(versions.values())) != 1:
        raise ValueError(f"package versions disagree: {versions}")
except (ValueError, KeyError, OSError) as error:
    print(error, file=sys.stderr)
    sys.exit(1)
print(f"Package versions agree: {versions['python']}")
