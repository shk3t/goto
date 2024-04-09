import pytest
import requests
from tests.utils import HOST


def teardown_module(module):
    requests.delete(f"{HOST}/utils/clear")


class TestCRUD:
    texts = ("buy products", "cook okroshka on kefir", "sleep")
    data = []

    @pytest.mark.parametrize("text", texts)
    def test_create(self, text):
        response = requests.post(f"{HOST}/todos", json={"text": text})
        assert response.status_code == 201
        assert response.json()["text"] == text
        self.data.append(response.json())

    @pytest.mark.parametrize("i", range(len(texts)))
    def test_get_one(self, i):
        id = self.data[i]["id"]
        response = requests.get(f"{HOST}/todo/{id}")
        assert response.status_code == 200
        assert self.data[i] == response.json()

    def test_delete(self):
        i = 1
        id = self.data.pop(i)["id"]
        response = requests.delete(f"{HOST}/todo/{id}")
        assert response.status_code == 200

    def test_get_all(self):
        response = requests.get(f"{HOST}/todos")
        assert response.status_code == 200
        assert self.data == response.json()