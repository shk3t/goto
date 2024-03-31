#!/usr/bin/env python

import os
import pathlib
import sys

import django
from asgiref.sync import async_to_sync, sync_to_async

os.chdir("..")
sys.path[0] = str(pathlib.Path(sys.path[0]).parent)
os.environ.setdefault("DJANGO_SETTINGS_MODULE", "config.settings")
django.setup()

from stats.api import *
from stats.models import *
from users.models import *

if __name__ == "__main__":
    WbAccount.objects.first().load_all_data()