FROM golang:1.11-alpine

RUN  apk add --update musl-dev gcc git

ENV UNPROCESSED_IMAGES=/tagallery/images/unprocessed
ENV GO111MODULE=on

# Override these env variables according to your settings
ENV PORT=3333
ENV DATABASE=tagallery
ENV DATABASE_HOST=tg-mongodb

# Create the image directory if it doesn't exist already
RUN mkdir -p ${UNPROCESSED_IMAGES}

WORKDIR /go/src/app
COPY . .

RUN go get -d ./...
RUN go build

CMD ["./api"]
