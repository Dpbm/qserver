FROM python:3.12.9-slim-bullseye

RUN apt update && apt install -y git gcc make cmake libopenblas-dev python3-dev build-essential curl
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

ENV PATH="/root/.cargo/bin:${PATH}"

WORKDIR /server
COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

HEALTHCHECK --interval=1m --timeout=10s --start-period=5s --retries=3 \
    CMD pidof python || exit 1

CMD ["python", "worker.py"]