FROM alpine:3.16.2

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go

RUN wget -O go.tgz https://dl.google.com/go/go1.10.3.src.tar.gz

RUN tar -C /usr/local -xzf go.tgz

WORKDIR /usr/local/go/xgss

COPY . .

RUN ./make.bash

RUN export PATH="/usr/local/go/bin:$PATH"

RUN export GOPATH=/opt/go/

RUN export PATH=$PATH:$GOPATH/bin

RUN apk del .build-dep

RUN go build -o build/ ./...

WORKDIR build/

CMD [". xgss"]