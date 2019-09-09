FROM golang:1.13 AS builder

WORKDIR /stash/src/github.com/acidlemon/guardmech
ENV GOPATH=/stash

ADD . /stash/src/github.com/acidlemon/guardmech
RUN go get && go test -v && go build -v -o guardmech cmd/guardmech/main.go && mv guardmech /stash/

FROM debian
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*
COPY --from=builder /stash/guardmech /usr/local/bin/guardmech

CMD ["/usr/local/bin/guardmech"]

