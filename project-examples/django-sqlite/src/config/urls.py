from django.contrib import admin
from django.urls import include, path, re_path
from django.views.generic import TemplateView

urlpatterns = [
    path("", include("main.urls")),
    path("utils/", include("utils.urls")),
]