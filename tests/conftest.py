import os
import uuid
from collections.abc import Iterator
from pathlib import Path

import pytest
from playwright.sync_api import APIRequestContext, Page, Playwright

FRONTEND_URL = os.environ.get("FRONTEND_URL", "http://localhost:3000")
API_URL = os.environ.get("API_URL", "http://localhost:8080/api/v1")


def _db_url_from_env_file() -> str | None:
    env_file = Path(__file__).resolve().parent.parent / ".env"
    if not env_file.exists():
        return None

    values: dict[str, str] = {}
    for line in env_file.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, _, value = line.partition("=")
        values[key.strip()] = value.strip().strip('"').strip("'")

    user = values.get("PG_USER", "postgres")
    password = values.get("PG_PASSWORD", "change_me")
    port = values.get("PG_PORT", "5433")
    name = values.get("PG_NAME", "xakat0n")
    return f"postgresql://{user}:{password}@localhost:{port}/{name}?sslmode=disable"


DB_URL = os.environ.get("DB_URL") or _db_url_from_env_file()


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


