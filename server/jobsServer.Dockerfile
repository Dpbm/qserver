FROM golang:1.23.5-alpine3.21 AS build

COPY ./jobsServer/ /go/src/server
COPY ./shared/ /go/src/shared

WORKDIR /go/src/shared
RUN go mod download && go mod verify

WORKDIR /go/src/server
RUN go mod download && go mod verify
RUN go build -o serverExec .

RUN mkdir -p /qasm
RUN mkdir -p /bin-tmp

WORKDIR /bin-tmp
RUN ARCH=$(uname -m) && \
    if [ $ARCH = 'aarch64' ]; then ARCH='arm64'; fi && \
    wget https://github.com/fullstorydev/grpcurl/releases/download/v1.9.2/grpcurl_1.9.2_linux_${ARCH}.tar.gz && \
    tar -xvf grpcurl_1.9.2_linux_${ARCH}.tar.gz

FROM scratch
COPY --from=build /go/src/server/serverExec /server
COPY --from=build /qasm /qasm
COPY --from=build /bin-tmp/grpcurl /grpcurl

HEALTHCHECK --interval=1m --timeout=10s --start-period=5s --retries=3 \
    CMD ["/grpcurl", "-plaintext", "172.18.0.29:50051", "Jobs/HealthCheck"]

EXPOSE 50051

CMD ["/server"]