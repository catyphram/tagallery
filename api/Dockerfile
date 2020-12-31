FROM golang:1-alpine AS builder
RUN apk add --update --no-cache musl-dev gcc git

WORKDIR /go/src/app
COPY . .

RUN go get -d ./...
RUN go build

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /go/src/app/api .
CMD ["./api"]  
