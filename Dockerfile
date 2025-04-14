FROM golang:1.24

RUN apt-get update && apt-get install -y make

WORKDIR ${GOPATH}/pvz/
COPY . ${GOPATH}/pvz/

RUN go build -o /build ./cmd/ \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]