FROM golang:1-alpine

RUN  apk add --update musl-dev gcc git

WORKDIR /go/src/app
COPY . .

RUN go get -d ./...
RUN go build

CMD ["./api"]