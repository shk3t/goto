import os
import sys
from pathlib import Path

import django

HOST = "http://localhost:8223"


def change_import_dir(path):
    os.chdir(path)
    sys.path[0] = str(Path.cwd())


class DjangoImport:
    setuped = False

    def __enter__(self):
        self.base_dir = Path.cwd()
        change_import_dir("src")

        if not DjangoImport.setuped:
            os.environ.setdefault("DJANGO_SETTINGS_MODULE", "config.settings")
            django.setup()
            DjangoImport.setuped = True

    def __exit__(self, exc_type, exc_value, exc_tb):
        change_import_dir(self.base_dir)