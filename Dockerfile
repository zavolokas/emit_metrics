FROM docker.svc.ring.com/golang:1.19-bullseye

COPY . /app
WORKDIR /app

RUN go env -w GOPROXY=direct
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["emit_metrics"]