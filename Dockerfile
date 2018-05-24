FROM golang:1.10 AS builder

ADD . /stash/src/github.com/acidlemon/guardmech
WORKDIR /stash/src/github.com/acidlemon/guardmech
ENV GOPATH=/stash

RUN go get && go test && go build -v -o guardmech cmd/guardmech/main.go && mv guardmech /stash/

FROM debian
COPY --from=builder /stash/guardmech /usr/local/bin/guardmech

ENTRYPOINT ["/usr/local/bin/guardmech"]

