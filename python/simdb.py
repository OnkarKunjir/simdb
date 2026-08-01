from typing import Any

import requests


class SimDB:
    def __init__(self, url: str, embedder_url: str, embedder_model: str) -> None:
        self.url = url
        self.embedder_url = embedder_url
        self.embedder_model = embedder_model

    def create_store(
        self, name: str, m: int = 16, ef_construct: int = 64, ef_search: int = 40
    ):
        response = requests.post(
            url=f"{self.url}/create-store",
            json={
                "name": name,
                "m": m,
                "efConstruct": ef_construct,
                "efSearch": ef_search,
                "url": self.embedder_url,
                "model": self.embedder_model,
            },
        )
        return response.json()

    def list_stores(self) -> list[str]:
        response = requests.get(url=f"{self.url}/list-stores")
        return response.json()

    def index_document(
        self, name: str, content: str, metadata: dict[str, str] | None = None
    ) -> dict[str, Any]:
        response = requests.post(
            url=f"{self.url}/index-document",
            json={"name": name, "content": content, "metadata": metadata},
        )

        return response.json()

    def search(self, name: str, search_string: str, k: int = 10) -> dict[str, Any]:
        response = requests.post(
            url=f"{self.url}/search",
            json={
                "name": name,
                "searchString": search_string,
                "k": k,
            },
        )
        return response.json()
