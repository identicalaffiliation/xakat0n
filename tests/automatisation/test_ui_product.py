import pytest
from playwright.sync_api import Page

SEED_LIMITED = "iPhone 16 Pro"
SEED_REGULAR = "Футбольный мяч"


def open_product(page: Page, title: str) -> str:
    page.get_by_text(title, exact=True).click()
    page.wait_for_url("**/product/*")
    return page.url.rstrip("/").split("/")[-1]


@pytest.fixture
def product_page(page: Page, logged_in_page: Page) -> Page:
    return logged_in_page


def test_limited_product_detail_shows_badge(product_page: Page) -> None:
    open_product(product_page, SEED_LIMITED)
    product_page.get_by_role("button", name="Купить").wait_for()
    product_page.get_by_text("Лимитированный").wait_for()
    product_page.get_by_text("Категория: Электроника").wait_for()


def test_regular_product_has_no_limited_badge(product_page: Page) -> None:
    open_product(product_page, SEED_REGULAR)
    product_page.get_by_role("button", name="Купить").wait_for()
    assert product_page.get_by_text("Лимитированный").count() == 0
    product_page.get_by_text("Категория: Спорт").wait_for()


def test_back_button_returns_to_catalog(product_page: Page) -> None:
    open_product(product_page, SEED_REGULAR)
    product_page.get_by_role("button", name="Назад").click()
    product_page.wait_for_url("**/products")
