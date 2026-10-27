import pytest


@pytest.fixture
def broken_fixture():
    raise RuntimeError("fixture setup failed")


def test_add_item_passes():
    cart = ["apple"]
    cart.append("pear")
    assert len(cart) == 2


def test_total_is_wrong_fails():
    assert sum([10, 20]) == 40, "total must include both items"


def test_fixture_error(broken_fixture):
    pass


def test_skipped():
    pytest.skip("not supported on this platform")


@pytest.mark.parametrize("percent", [10, 50])
def test_discount(percent):
    assert 100 - percent < 100


def test_prints_to_stdout():
    print("checkout started")
    assert True
