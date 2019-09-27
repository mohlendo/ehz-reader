FROM golang:alpine AS builder

RUN apk add --no-cache git \
    && go get github.com/tarm/serial \
    && go get github.com/influxdata/influxdb1-client/v2 \
    && apk del git 
# Create app directory
WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM scratch

ENV INFLUX_URL "http://influxdb:8086"
ENV SERIAL_PORT_NAME "/dev/ttyUSB0"

# Copy our static executable.
COPY --from=builder /go/bin/app /go/bin/app


CMD ["/go/bin/app"]