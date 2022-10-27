FROM golang:1.19.2 AS builder

COPY . /go/src/xgss

WORKDIR /go/src/xgss

RUN go get && CGO_ENABLED=0 GOOS=linux go build -o xgss .

FROM scratch

COPY --from=builder /go/src/xgss/xgss /opt/xgss/xgss

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 42100

WORKDIR /opt/xgss

ENTRYPOINT ["./xgss"]