import os
import uuid
from collections.abc import Iterator

import pytest
from playwright.sync_api import APIRequestContext, Page, Playwright

FRONTEND_URL = os.environ.get("FRONTEND_URL", "http://localhost:3000")
API_URL = os.environ.get("API_URL", "http://localhost:8080/api/v1")
DB_URL = os.environ.get("DB_URL")


def unique_username(prefix: str = "test") -> str:
    return f"{prefix}_{uuid.uuid4().hex[:12]}"


@pytest.fixture(scope="session")
def frontend_url() -> str:
    return FRONTEND_URL


@pytest.fixture(scope="session")
def api_url() -> str:
    return API_URL


@pytest.fixture
def username() -> str:
    return unique_username()


@pytest.fixture(scope="session")
def api(playwright: Playwright) -> Iterator[APIRequestContext]:
    context = playwright.request.new_context()
    yield context
    context.dispose()


@pytest.fixture(autouse=True)
def reset_db(api: APIRequestContext) -> None:
    if not DB_URL:
        return
    try:
        import psycopg2
    except ImportError:
        pytest.skip("psycopg2-binary не установлен — сброс БД пропущен")

    conn = psycopg2.connect(DB_URL)
    try:
        with conn.cursor() as cursor:
            cursor.execute("TRUNCATE queues;")
        conn.commit()
    finally:
        conn.close()


@pytest.fixture
def api_login(api: APIRequestContext, api_url: str) -> object:
    def _login(user: str) -> str:
        response = api.post(f"{api_url}/auth/login", data={"username": user})
        assert response.status == 200, f"login failed: {response.text()}"
        body = response.json()
        assert "token" in body, body
        return body["token"]

    return _login


def _inject_token(page: Page, token: str) -> None:
    import json

    page.add_init_script(
        f"sessionStorage.setItem('sessionToken', {json.dumps(token)})"
    )


@pytest.fixture
def logged_in_page(
    page: Page,
    frontend_url: str,
    username: str,
    api_login: object,
) -> Page:
    token = api_login(username)
    _inject_token(page, token)
    page.goto(f"{frontend_url}/products")
    page.get_by_placeholder("Поиск по объявлениям").wait_for()
    return page


@pytest.fixture
def expect_dialog(page: Page) -> object:
    import time
    from contextlib import contextmanager

    @contextmanager
    def _expect(timeout: float = 10000) -> Iterator[dict]:
        holder: dict = {}

        def handler(dialog) -> None:
            holder["message"] = dialog.message
            dialog.accept()

        page.on("dialog", handler)
        try:
            yield holder
            deadline = time.monotonic() + timeout / 1000
            while "message" not in holder and time.monotonic() < deadline:
                page.wait_for_timeout(100)
        finally:
            page.remove_listener("dialog", handler)

    return _expect


