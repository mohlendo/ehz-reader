FROM golang:alpine AS builder

ENV INFLUXDB_URL "http://influxdb:8086"
ENV INFLUXDB_USERNAME "admin"
ENV INFLUXDB_PASSWORD "admin"
ENV SERIAL_PORT_NAME "/dev/ttyUSB0"

RUN apk add --no-cache git

WORKDIR /build

COPY . .

RUN go get -d -v ./...
RUN go build ./ehz-reader.go

FROM scratch

COPY --from=builder /build/ehz-reader /go/bin/app

CMD ["/go/bin/app"]