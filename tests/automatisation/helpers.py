from playwright.sync_api import APIRequestContext

def seed_item_id(api: APIRequestContext, api_url: str, title: str) -> str:
    response = api.get(f"{api_url}/items")
    assert response.status == 200, response.text()
    items = response.json()
    assert isinstance(items, list) and items, "каталог пуст"
    for item in items:
        if item["title"] == title:
            return item["itemId"]
    raise AssertionError(f"товар '{title}' не найден в каталоге")
