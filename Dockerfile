FROM golang:alpine AS builder

RUN apk add --no-cache git

# Create appuser
RUN adduser -D -g '' appuser

# Create app directory
WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH GOARM=$GOARM go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/app .

FROM scratch

ENV INFLUX_URL "http://influxdb:8086"
ENV SERIAL_PORT_NAME "/dev/ttyUSB0"

# Import from builder.
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable.
COPY --from=builder /go/bin/app /go/bin/app

# Use an unprivileged user.
USER appuser

# Run the hello binary.
ENTRYPOINT ["/go/bin/app"]