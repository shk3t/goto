import json

from django.forms.models import model_to_dict
from django.http import JsonResponse
from django.http.response import Http404
from main.models import Todo


def todo_view(request, id):
    if request.method == "GET":
        todo = Todo.objects.get(id=id)
        return JsonResponse(model_to_dict(todo))

    elif request.method == "DELETE":
        Todo.objects.filter(id=id).delete()
        return JsonResponse({})

    else:
        raise Http404


def todos_view(request):
    if request.method == "GET":
        todos = list(Todo.objects.all())
        response = [model_to_dict(x) for x in todos]
        return JsonResponse(response, safe=False)

    elif request.method == "POST":
        payload = json.loads(request.body)
        todo = Todo.objects.create(**payload)
        return JsonResponse(model_to_dict(todo), status=201)

    else:
        raise Http404