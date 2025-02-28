FROM golang:1.23.5-alpine3.21 AS build

RUN apk --no-cache add ca-certificates

COPY ./restAPI/ /go/src/api
COPY ./shared/ /go/src/shared

WORKDIR /go/src/shared
RUN go mod download && go mod verify

WORKDIR /go/src/api
RUN go mod download && go mod verify
RUN go build -o serverExec .

FROM scratch
# https://stackoverflow.com/questions/52969195/docker-container-running-golang-http-client-getting-error-certificate-signed-by
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/api/serverExec /server

EXPOSE 3000

CMD ["/server"]