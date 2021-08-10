FROM golang:1.16 AS backend-builder

WORKDIR /stash/src/github.com/acidlemon/guardmech/backend
ENV GOPATH=/stash

ADD backend/ /stash/src/github.com/acidlemon/guardmech/backend
RUN go get && go test -v ./... && go build -v -o guardmech cmd/guardmech/main.go && mv guardmech /stash/ && cp run_guardmech.sh /stash/

FROM node:14-slim AS frontend-builder

WORKDIR /stash/src/github.com/acidlemon/guardmech/frontend

ADD frontend/ /stash/src/github.com/acidlemon/guardmech/frontend

RUN npm install && npm run build && mkdir -p /stash/guardmech && cp -a /stash/src/github.com/acidlemon/guardmech/frontend/dist /stash/guardmech/dist


FROM debian
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*
COPY --from=backend-builder /stash/guardmech /opt/guardmech/guardmech
COPY --from=backend-builder /stash/run_guardmech.sh /opt/guardmech/run_guardmech.sh
COPY --from=frontend-builder /stash/guardmech/dist /opt/guardmech/dist

WORKDIR /opt/guardmech

CMD ["sh", "-e", "/opt/guardmech/run_guardmech.sh"]

