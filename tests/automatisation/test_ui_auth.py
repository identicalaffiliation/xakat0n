import pytest
from playwright.sync_api import Page


@pytest.fixture
def auth_page(page: Page, frontend_url: str) -> Page:
    page.goto(f"{frontend_url}/")
    page.get_by_role("button", name="Войти").wait_for()
    return page


def test_login_via_ui_redirects_to_catalog(auth_page: Page, username: str) -> None:
    auth_page.get_by_label("Имя").fill(username)
    auth_page.get_by_role("button", name="Войти").click()
    auth_page.wait_for_url("**/products")
    auth_page.get_by_placeholder("Поиск по объявлениям").wait_for()


def test_login_with_empty_name_blocks_with_alert(
    auth_page: Page, expect_dialog: object
) -> None:
    with expect_dialog() as dialog:
        auth_page.get_by_role("button", name="Войти").click()
    assert dialog["message"] == "Введите имя"
    assert "products" not in auth_page.url
    auth_page.get_by_role("button", name="Войти").is_visible()


def test_login_with_too_short_name_shows_error(
    auth_page: Page, expect_dialog: object
) -> None:
    auth_page.get_by_label("Имя").fill("ab")
    with expect_dialog() as dialog:
        auth_page.get_by_role("button", name="Войти").click()
    assert dialog["message"] == "Ошибка входа"
