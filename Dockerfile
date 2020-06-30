FROM golang:latest AS builder

COPY . /go/src/github.com/dapixio/fio-store/
WORKDIR /go/src/github.com/dapixio/fio-store/

ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -o /fiostore cmd/fiostore/main.go

FROM scratch

COPY --chown=65535:65535 --from=builder /fiostore /fiostore

USER 65535
CMD ["/fiostore"]

