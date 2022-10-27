FROM golang:1.19-alpine3.16

WORKDIR /user/local/xgss

COPY build/xgss ./xgss

CMD ["xgss"]