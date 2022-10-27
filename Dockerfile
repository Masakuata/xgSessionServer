FROM golang:1.19-alpine3.16

WORKDIR /user/local/xgss

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/xgss/bin/ ./...

WORKDIR /usr/local/xgss/bin

CMD ["xgss"]