FROM golang:latest AS builder

COPY . /go/src/github.com/fioprotocol/fiostore/
WORKDIR /go/src/github.com/fioprotocol/fiostore/

ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -o /fiostore cmd/fiostore/main.go

FROM scratch

COPY --from=builder /fiostore /fiostore

USER 65535
CMD ["/fiostore"]

