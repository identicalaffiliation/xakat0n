import uuid

from playwright.sync_api import APIRequestContext, Browser, Page

from helpers import seed_item_id

SEED_LIMITED = "Сушёные комары, 1 кг"
SEED_REGULAR = "Сушёные тараканы, 1 кг"


def open_product(page: Page, title: str) -> str:
    page.get_by_text(title, exact=True).click()
    page.wait_for_url("**/product/*")
    return page.url.rstrip("/").split("/")[-1]


def login_via_ui(page: Page, frontend_url: str) -> None:
    page.goto(f"{frontend_url}/")
    page.get_by_label("Имя").fill(f"e2e_{uuid.uuid4().hex[:10]}")
    page.get_by_role("button", name="Войти").click()
    page.wait_for_url("**/products")


def test_buy_regular_item_opens_checkout(
    page: Page, logged_in_page: Page
) -> None:
    item_id = open_product(page, SEED_REGULAR)
    page.get_by_role("button", name="Купить").click()
    page.wait_for_url(f"**/product/{item_id}/checkout")
    page.get_by_text("Оформление заказа").wait_for()
    page.get_by_role("heading", name=SEED_REGULAR).wait_for()
    page.get_by_role("button", name="Оплатить").wait_for()


def test_checkout_cancel_returns_to_catalog(page: Page, logged_in_page: Page) -> None:
    item_id = open_product(page, SEED_REGULAR)
    page.get_by_role("button", name="Купить").click()
    page.wait_for_url(f"**/product/{item_id}/checkout")
    page.get_by_role("button", name="Отменить заказ").click()
    page.wait_for_url("**/products")


def test_buy_limited_item_gets_offer_and_checkout(
    page: Page, logged_in_page: Page
) -> None:
    item_id = open_product(page, SEED_LIMITED)
    page.get_by_role("button", name="Купить").click()
    page.wait_for_url(f"**/product/{item_id}/queue")
    page.get_by_role("heading", name="Товар освободился!").wait_for()

    page.get_by_role("button", name="Перейти к оформлению").click()
    page.wait_for_url(f"**/product/{item_id}/checkout")
    page.get_by_text("Оформление заказа").wait_for()
    page.get_by_role("heading", name=SEED_LIMITED).wait_for()

    page.get_by_role("button", name="Отменить заказ").click()
    page.wait_for_url("**/products")


def test_second_user_queues_while_first_holds_right(
    page: Page, browser: Browser, frontend_url: str, api_url: str, api: APIRequestContext
) -> None:
    item_id = seed_item_id(api, api_url, SEED_LIMITED)

    login_via_ui(page, frontend_url)
    page.goto(f"{frontend_url}/product/{item_id}")
    page.get_by_role("button", name="Купить").click()
    page.wait_for_url(f"**/product/{item_id}/queue")
    page.get_by_role("heading", name="Товар освободился!").wait_for()

    context = browser.new_context()
    page2 = context.new_page()
    try:
        login_via_ui(page2, frontend_url)
        page2.goto(f"{frontend_url}/product/{item_id}")
        page2.get_by_role("button", name="Купить").click()
        page2.wait_for_url(f"**/product/{item_id}/queue")
        page2.get_by_text("Очередь ожидания на лимитированный товар").wait_for()
        page2.get_by_text("Ваше место").wait_for()
    finally:
        context.close()

    page.get_by_role("button", name="Отказаться").click()
    page.wait_for_url("**/products")


def test_queue_page_without_ticket_shows_empty_state(
    page: Page, logged_in_page: Page, frontend_url: str, api_url: str, api: APIRequestContext
) -> None:
    item_id = seed_item_id(api, api_url, SEED_LIMITED)
    page.goto(f"{frontend_url}/product/{item_id}/queue")
    page.get_by_text("Нет активной заявки").wait_for()
    page.get_by_role("button", name="В каталог").click()
    page.wait_for_url("**/products")