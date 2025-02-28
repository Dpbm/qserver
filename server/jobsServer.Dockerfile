FROM golang:1.23.5-alpine3.21 AS build

COPY ./jobsServer/ /go/src/server
COPY ./shared/ /go/src/shared

WORKDIR /go/src/shared
RUN go mod download && go mod verify

WORKDIR /go/src/server
RUN go mod download && go mod verify
RUN go build -o serverExec .

RUN mkdir -p /qasm

FROM scratch
COPY --from=build /go/src/server/serverExec /server
COPY --from=build /qasm /qasm

EXPOSE 50051

CMD ["/server"]