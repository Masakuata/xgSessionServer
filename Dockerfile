FROM golang:1.19.2-bullseye

WORKDIR /usr/local/xgss

COPY . .

RUN go build -o build/

WORKDIR build/

CMD [". xgss"]