from django.urls import path
from main.views import todo_view, todos_view

urlpatterns = [
    path("todo/<int:id>", todo_view),
    path("todos", todos_view),
]