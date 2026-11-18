# Runtime tests for the hand-written kei.connection package: path rendering
# and error hygiene (no credential or response body ever reaches an error).
import pytest

from kei.connection import OperationError, map_error, render_path


def test_render_path_escapes_segments():
    assert (
        render_path("/repos/{owner}/{repository}", {"owner": "a/b", "repository": "x y"})
        == "/repos/a%2Fb/x%20y"
    )
    with pytest.raises(ValueError):
        render_path("/{missing}", {})


def test_render_path_matches_go_path_escape():
    # Go's url.PathEscape is the reference: keep $ & + : = @, escape ! ' ( ) * , ;
    assert (
        render_path("/x/{v}", {"v": "a!b'c(d)e*f,g;h$i&j+k:l=m@n"})
        == "/x/a%21b%27c%28d%29e%2Af%2Cg%3Bh$i&j+k:l=m@n"
    )
    assert render_path("/x/{v}", {"v": "d/a!b"}, multi_segment={"v"}) == "/x/d/a%21b"


def test_render_path_multi_segment():
    assert (
        render_path("/contents/{path}", {"path": "docs/notes.txt"}, multi_segment={"path"})
        == "/contents/docs/notes.txt"
    )
    with pytest.raises(ValueError):
        render_path("/contents/{path}", {"path": "docs/../etc"}, multi_segment={"path"})


def test_no_credential_in_errors():
    # The provider echoes the credential in its 401 body; the error must
    # carry only the op id, kind and status.
    body = '{"message": "denied for tok_SECRET"}'
    err = map_error("github.issues.get", 401, [], body)
    assert isinstance(err, OperationError)
    assert "tok_SECRET" not in str(err)
    assert "tok_SECRET" not in repr(err)
