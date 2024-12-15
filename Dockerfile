ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

FROM debian:bookworm

RUN apt update -y
RUN apt-get install ca-certificates -y
RUN update-ca-certificates

COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/static /usr/local/bin/
CMD ["run-app", "server"]
