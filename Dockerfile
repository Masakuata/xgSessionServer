FROM golang:1.19-alpine3.16 AS builder

RUN apk update \
    && apk add --no-cache git \
    && apk add --no-cache ca-certificates \
    && apk add --update gcc musl-dev alpine-sdk \
    && update-ca-certificates

WORKDIR $GOPATH

# Fetch dependencies.
RUN go get -u -d -v

COPY . .

# Go build the binary, specifying the final OS and architecture we're looking for
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app/

FROM scratch

# Import the files from the builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable.
COPY --from=builder /go/bin/app/ /go/bin/app/

WORKDIR /go/bin/app

CMD [". xgss"]