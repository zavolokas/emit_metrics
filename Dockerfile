FROM golang:1.20-bullseye as builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./emit_metrics ./cmd

FROM debian:bullseye

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /build/emit_metrics /emit_metrics

RUN groupadd --gid 900 bot && useradd -M --shell /usr/sbin/nologin --uid 900 --gid 900 bot && \
    openssl req -newkey rsa:2048 -nodes -keyout /server.key -x509 -days 3650 -out /server.crt -subj '/CN=server' && \
    chown bot /server.key /server.crt && chmod 0400 /server.key /server.crt
    
USER bot

ENTRYPOINT ["/emit_metrics"]