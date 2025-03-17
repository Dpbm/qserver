FROM python:3.12.9-alpine3.21

ENV PYTHONUNBUFFERED=1

WORKDIR /server
COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

HEALTHCHECK --interval=1m --timeout=10s --start-period=5s --retries=3 \
    CMD pidof python || exit 1

CMD ["python", "-u", "worker.py"]