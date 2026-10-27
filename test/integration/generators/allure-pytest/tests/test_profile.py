import base64

import allure
import pytest

PNG = base64.b64decode(
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVR4nGP4//8/AAX+Av4N70a4AAAAAElFTkSuQmCC"
)


@allure.suite("Profile")
@allure.severity(allure.severity_level.CRITICAL)
@allure.label("layer", "api")
@allure.link("PRJ-12", link_type="tms", name="PRJ-12")
def test_update_profile_passes():
    with allure.step("Open profile page"):
        with allure.step("Load user data"):
            assert True
        with allure.step("Capture page"):
            allure.attach(PNG, name="profile page", attachment_type=allure.attachment_type.PNG)
    allure.attach("profile updated", name="response body", attachment_type=allure.attachment_type.TEXT)


@allure.suite("Profile")
def test_display_name_fails():
    with allure.step("Check display name"):
        assert "Alice" == "Bob", "display name must match"


@allure.suite("Profile")
def test_avatar_upload_broken():
    with allure.step("Upload avatar"):
        raise RuntimeError("storage backend unavailable")


@allure.suite("Profile")
def test_two_factor_skipped():
    pytest.skip("two-factor is not supported yet")


@allure.suite("Profile")
@pytest.mark.parametrize("language", ["en", "de"])
def test_language_switch(language):
    with allure.step("Switch language"):
        assert language in ("en", "de")
