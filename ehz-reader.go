package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	"go.bug.st/serial"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type measurement struct {
	name       string
	pattern    []byte
	startIndex int
	length     int
	divisor    float64
}

var measurements = []measurement{
	measurement{name: "power", pattern: []byte{'\x07', '\x01', '\x00', '\x10', '\x07', '\x00'}, startIndex: 8, length: 4, divisor: 1},
	measurement{name: "total", pattern: []byte{'\x07', '\x01', '\x00', '\x01', '\x08', '\x00'}, startIndex: 12, length: 8, divisor: 10000}}

func splitMsg(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\x1b', '\x1b', '\x1b', '\x1b', '\x01', '\x01', '\x01', '\x01'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func parseMsg(msg []byte) (map[string]interface{}){
	// log.Printf("%x\n", msg)
	fields := make(map[string]interface{})
	for _, m := range measurements {
		if i := bytes.Index(msg, m.pattern); i > 0 {
			l := len(m.pattern)
			start := i + l + m.startIndex
			slice := msg[start : start+m.length]
			var value float64
			if m.length == 8 {
				value = float64(binary.BigEndian.Uint64(slice))
			} else {
				value = float64(binary.BigEndian.Uint32(slice))
			}
			fields[m.name] = value / m.divisor
		}
	}
	return fields
}

func writePoints(influx client.Client, fields map[string]interface{}, ts time.Time) {
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: os.Getenv("INFLUXDB_DB")})

	// Create a point and add to batch
	tags := map[string]string{"meter": "household"}
	pt, _ := client.NewPoint("power_consumption", tags, fields, ts)
	bp.AddPoint(pt)

	if err := influx.Write(bp); err != nil {
		log.Fatal("Error writing influx data", err)
	}
}

func main() {
	serialPort, envPresent := os.LookupEnv("SERIAL_PORT_NAME")
	if envPresent != true {
		log.Fatal("Env var SERIAL_PORT_NAME not set.")
	}
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open(serialPort, mode)
	if err != nil {
		msg := fmt.Sprintf("Cannot open '%s' - ", serialPort)
		log.Fatal(msg, err)
	}
	defer port.Close()	
	
	var influx client.Client
	influxUrl, envPresent := os.LookupEnv("INFLUXDB_URL")
	if envPresent == true {
 		influx, err = client.NewHTTPClient(client.HTTPConfig{
			Addr: influxUrl,
		})
		if err != nil {
			msg := fmt.Sprintf("Cannot reach influxdb '%s'", influxUrl)
			log.Fatal(msg, err)
		}
		defer influx.Close()				
	}

	scanner := bufio.NewScanner(bufio.NewReader(port))
	scanner.Buffer(make([]byte, 2048), 4*1024)
	scanner.Split(splitMsg)

	for scanner.Scan() {
		fields := parseMsg(scanner.Bytes())
		log.Printf("fields: %v", fields)		
		if influx != nil && len(fields) > 0 {
			go writePoints(influx, fields, time.Now())
		}
	}
	if scanner.Err() != nil {
        log.Fatal("Read from serial port aborted", scanner.Err())
    }
	log.Print("Done")
}
