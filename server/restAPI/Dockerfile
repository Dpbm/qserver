FROM golang:1.23.5-alpine3.21 AS build

WORKDIR /go/src

COPY . .
RUN go mod download && go mod verify
RUN go build -o server .

FROM scratch
COPY --from=build /go/src/server /server

CMD ["/server"]