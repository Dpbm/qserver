FROM nginx:1.27-alpine3.21

RUN mkdir -p /logs/nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

HEALTHCHECK --interval=1m --timeout=10s --start-period=5s --retries=3 \
    CMD wget --spider -q http://172.18.0.30:8080/healthcheck/ || exit 1