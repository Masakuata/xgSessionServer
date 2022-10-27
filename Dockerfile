FROM golang:1.19-alpine3.16

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app/ ./...

WORKDIR /usr/local/bin/app

CMD ["xgss"]