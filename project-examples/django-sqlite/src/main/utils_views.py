from django.http import JsonResponse
from django.http.response import Http404
from main.models import Todo


def ready_view(request):
    return JsonResponse({})

def clear_view(request):
    if request.method == "DELETE":
        Todo.objects.all().delete()
        return JsonResponse({})
    else:
        raise Http404