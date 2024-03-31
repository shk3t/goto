import pytest
import requests
import env


class TestCRUD:
    texts = ("buy products", "cook okroshka on kefir", "sleep")
    data = []

    @pytest.mark.parametrize("text", texts)
    def test_create(self, text):
        response = requests.post(f"{env.HOST}/todos", json={"text": text})
        assert response.status_code == 201
        assert response.json()["text"] == text
        self.data.append(response.json())

    @pytest.mark.parametrize("i", range(len(texts)))
    def test_get_one(self, i):
        print(i)
        print("CALLED")
        id = self.data[i]["id"]
        response = requests.get(f"{env.HOST}/todo/{id}")
        assert response.status_code == 200
        assert self.data[i] == response.json()

    def test_delete(self):
        i = 1
        id = self.data.pop(i)["id"]
        response = requests.delete(f"{env.HOST}/todo/{id}")
        assert response.status_code == 200

    def test_get_all(self):
        response = requests.get(f"{env.HOST}/todos")
        assert response.status_code == 200
        assert self.data == response.json()