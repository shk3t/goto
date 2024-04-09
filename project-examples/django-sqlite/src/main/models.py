from django.db import models


class Todo(models.Model):
    collection = models.ForeignKey(
        "main.Collection", on_delete=models.CASCADE, blank=True, null=True
    )
    text = models.TextField()
    completed = models.BooleanField(default=False)


class Collection(models.Model):
    text = models.TextField()
    completed = models.BooleanField(default=False)