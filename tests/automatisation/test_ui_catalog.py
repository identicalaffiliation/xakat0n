import pytest
from playwright.sync_api import Page

SEED_LIMITED = "Сушёные комары, 1 кг"
SEED_REGULAR = "Сушёные тараканы, 1 кг"
SEED_OTHER = "Пакет пакетов обычный, свежий"


@pytest.fixture
def catalog(page: Page, logged_in_page: Page) -> Page:
    return logged_in_page


def test_seed_items_and_limited_badge_visible(catalog: Page) -> None:
    catalog.get_by_text(SEED_LIMITED, exact=True).wait_for()
    catalog.get_by_text(SEED_REGULAR, exact=True).wait_for()
    catalog.get_by_text("Участок на Луне, Море Спокойствия, 1 га", exact=True).wait_for()
    assert catalog.get_by_text("Лимитированный").count() >= 1


def test_search_filters_cards(catalog: Page) -> None:
    search = catalog.get_by_placeholder("Поиск по объявлениям")
    search.fill("комары")
    catalog.get_by_text(SEED_LIMITED, exact=True).wait_for()
    assert catalog.get_by_text(SEED_REGULAR, exact=True).count() == 0


def test_search_cleared_restores_all_cards(catalog: Page) -> None:
    search = catalog.get_by_placeholder("Поиск по объявлениям")
    search.fill("комары")
    catalog.get_by_text(SEED_LIMITED, exact=True).wait_for()
    search.fill("")
    catalog.get_by_text(SEED_REGULAR, exact=True).wait_for()


def test_search_with_no_results_shows_empty_message(catalog: Page) -> None:
    catalog.get_by_placeholder("Поиск по объявлениям").fill("такого товара нет")
    catalog.get_by_text("Ничего не найдено").wait_for()


def test_category_filter_shows_only_matching(catalog: Page) -> None:
    catalog.get_by_role("button", name="Сушёные насекомые").click()
    catalog.get_by_text(SEED_LIMITED, exact=True).wait_for()
    catalog.get_by_text(SEED_REGULAR, exact=True).wait_for()
    assert catalog.get_by_text(SEED_OTHER, exact=True).count() == 0


def test_product_card_opens_detail_page(catalog: Page) -> None:
    catalog.get_by_text(SEED_REGULAR, exact=True).click()
    catalog.wait_for_url("**/product/*")
    catalog.get_by_role("button", name="Купить").wait_for()