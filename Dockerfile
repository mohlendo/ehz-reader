FROM arm64v8/golang:alpine

RUN apk add --no-cache git \
    && go get github.com/tarm/serial \
    && go get github.com/influxdata/influxdb1-client/v2 \
    && apk del git 
# Create app directory
WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]