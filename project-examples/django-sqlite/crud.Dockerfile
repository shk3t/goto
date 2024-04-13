FROM python:3.11-alpine
WORKDIR /django-sqlite
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY src src
CMD ["pytest", "tests/test_crud.py"]
