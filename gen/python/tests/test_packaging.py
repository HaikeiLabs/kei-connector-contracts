# Packaging tests: pyproject.toml must stay valid for PyPI (the 0.20.2
# publish failed with "project.urls.dependencies must be string" because
# the [project.urls] table swallowed the dependencies key).
import shutil
import subprocess
import sys
import tomllib
from pathlib import Path

import pytest

PYPROJECT = Path(__file__).resolve().parent.parent / "pyproject.toml"


def test_dependencies_is_a_list_with_protobuf():
    with PYPROJECT.open("rb") as f:
        data = tomllib.load(f)
    deps = data["project"]["dependencies"]
    assert isinstance(deps, list) and deps
    assert any(req.startswith("protobuf") for req in deps)


def test_urls_values_are_strings():
    with PYPROJECT.open("rb") as f:
        data = tomllib.load(f)
    urls = data["project"]["urls"]
    assert urls
    for key, value in urls.items():
        assert isinstance(value, str), f"project.urls.{key} must be a string"
    assert urls["Repository"] == (
        "https://github.com/HaikeiLabs/kei-connector-contracts"
    )


def test_build_sdist_and_wheel(tmp_path):
    if shutil.which("uv"):
        cmd = ["uv", "build", "--sdist", "--wheel", "--out-dir", str(tmp_path)]
    elif _build_module_available():
        cmd = [sys.executable, "-m", "build", "--sdist", "--wheel", "--outdir", str(tmp_path)]
    else:
        pytest.skip("neither uv nor python -m build available")
    result = subprocess.run(
        cmd, cwd=PYPROJECT.parent, capture_output=True, text=True
    )
    assert result.returncode == 0, (
        f"build failed:\n{result.stdout}\n{result.stderr}"
    )
    assert any(tmp_path.glob("*.tar.gz")), "no sdist produced"
    assert any(tmp_path.glob("*.whl")), "no wheel produced"


def _build_module_available() -> bool:
    try:
        import importlib.util

        return importlib.util.find_spec("build") is not None
    except Exception:
        return False
