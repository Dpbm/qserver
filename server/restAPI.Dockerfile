FROM golang:1.23.5-alpine3.21 AS build

RUN apk --no-cache add ca-certificates

COPY ./restAPI/ /go/src/api
COPY ./shared/ /go/src/shared

WORKDIR /go/src/shared
RUN go mod download && go mod verify

WORKDIR /go/src/api
RUN go mod download && go mod verify
RUN go build -o serverExec .

FROM busybox:1.37.0
# https://stackoverflow.com/questions/52969195/docker-container-running-golang-http-client-getting-error-certificate-signed-by
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/api/serverExec /server

HEALTHCHECK --interval=1m --timeout=10s --start-period=5s --retries=3 \
    CMD wget --spider -q http://172.18.0.28:3000/api/v1/health/ || exit 1

EXPOSE 3000

CMD ["/server"]