import os
import shutil
import subprocess
import time
from pathlib import Path

import requests
from tests.utils import HOST

server_process: subprocess.Popen


def pytest_sessionstart(session):
    global server_process

    if os.path.isdir("src/main/migrations"):
        shutil.rmtree("src/main/migrations")
    Path("db.sqlite3").unlink(missing_ok=True)
    subprocess.run(["python", "src/manage.py", "makemigrations", "main"])
    subprocess.run(["python", "src/manage.py", "migrate"])

    server_process = subprocess.Popen(["python", "src/manage.py", "runserver"])

    while True:
        time.sleep(1)
        response = requests.get(f"{HOST}/utils/ready")
        if response.status_code == 200:
            break


def pytest_sessionfinish(session, exitstatus):
    server_process.terminate()