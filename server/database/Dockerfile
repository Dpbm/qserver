FROM postgres:17.2-alpine3.21

COPY entrypoint.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/entrypoint.sh

ENV QASM_PATH=/qasm
RUN mkdir -p $QASM_PATH

HEALTHCHECK --interval=5s --timeout=30s --start-period=5s --retries=3 \
    CMD pg_isready -U postgres