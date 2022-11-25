FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/app

COPY . .

RUN go test
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v

FROM scratch

ENV INFLUXDB_URL "http://influxdb:8086"
ENV INFLUXDB_USERNAME "admin"
ENV INFLUXDB_PASSWORD "admin"
ENV SERIAL_PORT_NAME "/dev/ttyUSB0"

COPY --from=builder /go/src/app/ehz-reader /ehz-reader

ENTRYPOINT ["/ehz-reader"]