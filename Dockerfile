FROM alpine:3.16.2

WORKDIR /usr/local/xgss

COPY . .

RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go \
    wget -O go.tgz https://dl.google.com/go/go1.10.3.src.tar.gz \
    tar -C /usr/local -xzf go.tgz \
    cd /usr/local/go/src/ \
    ./make.bash \
    export PATH="/usr/local/go/bin:$PATH" \
    export GOPATH=/opt/go/ \
    export PATH=$PATH:$GOPATH/bin \
    apk del .build-deps \
    go build -o build/ ./...

WORKDIR build/

CMD [". xgss"]