import pytest
import requests
from django.forms.models import model_to_dict
from tests.utils import HOST, DjangoImport

with DjangoImport():
    from main.models import Collection, Todo


def teardown_module(module):
    requests.delete(f"{HOST}/utils/clear")


class TestModels:
    collections_data = [
        {"id": 1, "name": "utrom", "description": "nice"},
        {"id": 2, "name": "vecherom", "description": "very nice"},
    ]
    todos_data = [
        {"id": 1, "collection_id": 1, "text": "press kachat", "completed": True},
        {"id": 2, "collection_id": 1, "text": "begit", "completed": True},
        {"id": 3, "collection_id": 1, "text": "turnik", "completed": True},
        {"id": 4, "collection_id": 1, "text": "anjumanya", "completed": True},
        {"id": 5, "collection_id": 2, "text": "press kachat", "completed": False},
        {"id": 6, "collection_id": 2, "text": "begit", "completed": False},
        {"id": 7, "collection_id": 2, "text": "turnik", "completed": False},
        {"id": 8, "collection_id": 2, "text": "anjumanya", "completed": False},
    ]

    def test_collection(self):
        Collection.objects.bulk_create(
            [Collection(**x) for x in TestModels.collections_data]
        )
        created_collections = list(Collection.objects.all().values())
        assert TestModels.collections_data == created_collections

    def test_todo(self):
        Todo.objects.bulk_create([Todo(**x) for x in TestModels.todos_data])
        created_todos = list(Todo.objects.all().values())
        assert TestModels.todos_data == created_todos

    def test_collection_and_todo(self):
        utrom_collection = Collection.objects.get(pk=1)
        related_todos = list(utrom_collection.todos.all().values())  # type: ignore
        assert TestModels.todos_data[: len(related_todos)] == related_todos