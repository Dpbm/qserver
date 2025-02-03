FROM golang:1.23.5-alpine3.21 AS build

WORKDIR /go/src

COPY . .
RUN go mod download && go mod verify
RUN go build -C ./server/

RUN mkdir -p /qasm

FROM scratch
COPY --from=build /go/src/server/server /bin/server
COPY --from=build /qasm /qasm

EXPOSE 50051

CMD ["/bin/server"]

