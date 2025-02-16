FROM golang:1.23.5-alpine3.21 AS build

RUN apk --no-cache add ca-certificates

WORKDIR /go/src

COPY . .
RUN go mod download && go mod verify
RUN go build -o server .

FROM scratch
# https://stackoverflow.com/questions/52969195/docker-container-running-golang-http-client-getting-error-certificate-signed-by
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/server /server

CMD ["/server"]