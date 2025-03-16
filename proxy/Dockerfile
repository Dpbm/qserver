FROM nginx:1.27-alpine3.21

COPY setup-logs.sh /docker-entrypoint.d/

# once after mounting the volume the directory will be replaced
# we need to add an entrypoint to ensure the directory will be setup
# after mounting the logs volume
RUN chmod +x /docker-entrypoint.d/setup-logs.sh

COPY nginx.conf /etc/nginx/

HEALTHCHECK --interval=1m --timeout=10s --start-period=5s --retries=3 \
    CMD wget --spider -q http://172.18.0.30:8080/healthcheck/ || exit 1