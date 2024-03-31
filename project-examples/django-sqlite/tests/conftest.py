import subprocess
import time

import pathlib
import requests
import env

server_process: subprocess.Popen


def pytest_sessionstart(session):
    global server_process

    pathlib.Path("./db.sqlite3").unlink(missing_ok=True)
    subprocess.run(["python", "src/manage.py", "makemigrations", "main"])
    subprocess.run(["python", "src/manage.py", "migrate"])

    server_process = subprocess.Popen(["python", "src/manage.py", "runserver"])

    while True:
        time.sleep(1)
        response = requests.get(f"{env.HOST}/ready")
        if response.status_code == 200:
            break


def pytest_sessionfinish(session, exitstatus):
    server_process.terminate()