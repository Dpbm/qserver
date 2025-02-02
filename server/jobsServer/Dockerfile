FROM golang:1.23.5-alpine3.21 

WORKDIR /go/src

COPY . .
RUN go mod download && go mod verify
RUN go build -C ./server/

RUN mkdir -p /qasm

CMD ["/go/src/server/server"]

