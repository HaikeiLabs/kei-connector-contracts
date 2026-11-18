# Path rendering: every value is one escaped segment, so no value can
# change the shape of the path. Mirrors gen/go/kei/connection/path.go.
from typing import Collection, Mapping
from urllib.parse import quote

# Characters url.PathEscape keeps: A-Za-z0-9 (always safe in urllib.parse.quote)
# plus the sub-delims and :@ it does not escape.
_SAFE_SEGMENT = "-._~!$&'()*+,;=:@"


def _escape_segment(v: str) -> str:
    """Percent-encode v as one path segment."""
    if v.strip(".") == "":
        # A value made only of dots ("." or "..") is encoded as %2E so it is
        # never read as a relative segment.
        return "%2E" * len(v)
    return quote(v, safe=_SAFE_SEGMENT)


def render_path(
    template: str,
    values: Mapping[str, str],
    multi_segment: Collection[str] = (),
) -> str:
    """Fill a path template such as "/repos/{owner}/{repository}" from values.

    Every value is percent-encoded as a single path segment, so no value can
    change the shape of the path: "/" becomes %2F, "%" becomes %25, "?" and
    "#" are encoded, and a value made only of dots ("." or "..") is encoded
    as %2E so it is never read as a relative segment. A missing or empty
    value, an unterminated or empty variable, or a template that does not
    start with "/" raises ValueError.

    A variable named in multi_segment (a param with multi_segment, such as a
    file path or a Git ref) keeps "/" between its segments and escapes each
    segment; an empty, "." or ".." segment raises ValueError.
    """
    if not template.startswith("/"):
        raise ValueError(f"connection: path template {template!r} must start with /")
    out: list[str] = []
    rest = template
    while True:
        open_idx = rest.find("{")
        close_idx = rest.find("}")
        if open_idx < 0:
            if close_idx >= 0:
                raise ValueError(f"connection: path template {template!r} has an unmatched }}")
            out.append(rest)
            return "".join(out)
        if close_idx >= 0 and close_idx < open_idx:
            raise ValueError(f"connection: path template {template!r} has an unmatched }}")
        end = rest.find("}", open_idx)
        if end < 0:
            raise ValueError(f"connection: path template {template!r} has an unterminated variable")
        name = rest[open_idx + 1 : end]
        if name == "" or "{" in name or "/" in name:
            raise ValueError(f"connection: path template {template!r} has an invalid variable")
        if name not in values:
            raise ValueError(f"connection: path variable {name!r} has no value")
        v = values[name]
        if v == "":
            raise ValueError(f"connection: path variable {name!r} is empty")
        out.append(rest[:open_idx])
        if name in multi_segment:
            segs = v.split("/")
            for i, s in enumerate(segs):
                if s == "" or s == "." or s == "..":
                    raise ValueError(
                        f"connection: path variable {name!r} has an empty, . or .. segment"
                    )
                segs[i] = quote(s, safe=_SAFE_SEGMENT)
            out.append("/".join(segs))
        else:
            out.append(_escape_segment(v))
        rest = rest[end + 1 :]
