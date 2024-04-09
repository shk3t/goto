import pytest
import requests
from tests.utils import HOST, DjangoImport

with DjangoImport():
    from main.models import Todo


def teardown_module(module):
    requests.delete(f"{HOST}/utils/clear")


class TestModels:
    def test_collection(self):
        todos = Todo.objects.all()
        print(todos)

    def test_todo(self):
        pass