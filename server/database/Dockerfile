FROM postgres:17.2-alpine3.21

COPY postgresql.conf /docker-entrypoint-initdb.d/
COPY entrypoint.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/entrypoint.sh

ENV QASM_PATH=/qasm
RUN mkdir -p $QASM_PATH
RUN mkdir -p /var/log/postgres
RUN chmod -R 777 /var/log/postgres && \
    chown postgres /var/log/postgres

HEALTHCHECK --interval=5s --timeout=30s --start-period=5s --retries=3 \
    CMD pg_isready -U postgres

CMD ["postgres", "-c", "config_file=/docker-entrypoint-initdb.d/postgresql.conf"]
