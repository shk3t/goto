FROM python:3.11-alpine
WORKDIR /django-sqlite-example
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
CMD ["pytest ./tests/test_crud.py"]