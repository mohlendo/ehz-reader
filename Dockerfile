FROM golang:alpine

ENV INFLUX_URL "http://influxdb:8086"
ENV SERIAL_PORT_NAME "/dev/ttyUSB0"

RUN apk add --no-cache git

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]