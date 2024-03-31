from django.urls import path
from main.utils_views import clear_view, ready_view

urlpatterns = [
    path("ready", ready_view),
    path("clear", clear_view),
]