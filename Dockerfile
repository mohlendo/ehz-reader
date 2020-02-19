FROM golang:alpine AS builder

ENV INFLUXDB_URL "http://influxdb:8086"
ENV INFLUXDB_USERNAME "admin"
ENV INFLUXDB_PASSWORD "admin"
ENV SERIAL_PORT_NAME "/dev/ttyUSB0"

RUN apk add --no-cache git

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM scratch

COPY --from=builder /go/src/app/ehz-reader /go/bin/app

CMD ["/go/bin/app"]