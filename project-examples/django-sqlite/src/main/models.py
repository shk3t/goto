from django.db import models


class Todo(models.Model):
    collection = models.ForeignKey(
        "main.Collection",
        on_delete=models.CASCADE,
        related_name="todos",
        blank=True,
        null=True,
    )
    text = models.TextField()
    completed = models.BooleanField(default=False)


class Collection(models.Model):
    name = models.CharField(max_length=64)
    description = models.TextField()