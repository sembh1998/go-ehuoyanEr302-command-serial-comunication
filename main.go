package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// List of know status values. Status value is returned when executing device
// commands in raw mode, with RawCommand() method.
const (
	StatusOK    byte = 0x00
	StatusNoTag byte = 0x01
)

// Command is an enum type for encoding different commands that can be sent to
// device.
type Command uint16

// List of supported commands.
const (
	CommandInitializePort        Command = 0x0101
	CommandSetDeviceNodeNumber   Command = 0x0201
	CommandReadDeviceNodeNumber  Command = 0x0301
	CommandReadDeviceMode        Command = 0x0401
	CommandSetBuzzerBeep         Command = 0x0601
	CommandSetLedColor           Command = 0x0701
	CommandRFU                   Command = 0x0801
	CommandSetAntennaStatus      Command = 0x0c01
	CommandMifareRequest         Command = 0x0102
	CommandMifareAnticollision   Command = 0x0202
	CommandMifareSelect          Command = 0x0302
	CommandMifareHlta            Command = 0x0402
	CommandMifareAuthentication2 Command = 0x0702
	CommandMifareRead            Command = 0x0802
	CommandMifareWrite           Command = 0x0902
	CommandMifareInitval         Command = 0x0A02
	CommandMifareReadBalance     Command = 0x0B02
	CommandMifareDecrement       Command = 0x0C02
	CommandMifareIncrement       Command = 0x0D02
	CommandRF_UL_SELECT          Command = 0x1202
	CommandRF_UL_WRITE           Command = 0x1302
)

const (
	HeadSize         int = 2
	LenghtSize       int = 2
	NodeIDSize       int = 2
	FunctionCodeSize int = 2
	XORSize          int = 1
	StatusSize       int = 1
)

// LedMode encodes device led status - off, red or green.
type LedMode byte

// List of LED modes accepted by LED change command.
const (
	LedOff  LedMode = 0x00
	LedBlue LedMode = 0x01
	LedRed  LedMode = 0x02
)

var NodeId uint16

// beepUnit is a minimum beep duration supported by this device. When sending
// beep command to the device beep duration is sent as a number of beep units.

// msgPrefix is a static header added for all messages (sent and received) when
// communicating with device.
const msgPrefix uint16 = 0xAABB

func xorChecksum(buf []byte) byte {
	var csum byte = 0x00
	for _, b := range buf {
		csum ^= b
	}
	return csum
}

func copyUint16(buf []byte, n uint16) int {
	binary.BigEndian.PutUint16(buf, n)
	return 2
}

func newRequest(nodeId []byte, cmd Command, data []byte) []byte {
	// header+length+command+data+checksum
	tamanio := HeadSize + LenghtSize + NodeIDSize + FunctionCodeSize + len(data) + XORSize
	fmt.Print(tamanio)
	buf := make([]byte, tamanio)
	pos := 0
	pos += copyUint16(buf[pos:], msgPrefix)                                             // header
	pos += copyUint16(buf[pos:], uint16(NodeIDSize+FunctionCodeSize+len(data)+XORSize)) // length
	//invertimos posicion
	x := buf[2]
	y := buf[3]
	buf[2] = y
	buf[3] = x
	pos += copy(buf[pos:], nodeId)            // node ID
	pos += copyUint16(buf[pos:], uint16(cmd)) // command
	pos += copy(buf[pos:], data)              // data
	buf[pos] = xorChecksum(buf[4:pos])        // checksum

	return buf
}

func rx(r io.Reader) ([]byte, error) {
	fmt.Println("rx: 1")
	header := make([]byte, 4)
	fmt.Println("rx: 2")
	pos, err := io.ReadAtLeast(r, header, len(header))
	fmt.Println("rx: 3")
	if err != nil {
		return nil, fmt.Errorf("header read error: %s", err)
	}
	if binary.BigEndian.Uint16(header) != msgPrefix {
		return nil, errors.New("bad response header")
	}
	length := binary.BigEndian.Uint16(header[2:])
	if length < 4 {
		return nil, errors.New("invalid payload length value")
	}
	buf := make([]byte, len(header)+int(length))
	copy(buf, header)
	_, err = io.ReadAtLeast(r, buf[pos:], int(length))
	if err != nil {
		return nil, fmt.Errorf("payload read error: %s", err)
	}
	if buf[len(buf)-1] != xorChecksum(buf[4:len(buf)-1]) {
		return nil, errors.New("response have invalid checksum")
	}
	return buf, nil
}

type Device struct {
	port *serial.Port
}

func OpenDevice(dev string) (*Device, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: dev,
		Baud: 19200,
	})
	if err != nil {
		return nil, fmt.Errorf("rfid: error opening serial port: %s", err)
	}
	return &Device{
		port: port,
	}, nil
}

/*
func main() {
	//var cmd Command
	var data []byte
	d, err := OpenDevice("COM3")

	req := newRequest(0x0106, data)
	fmt.Println("TX:", req)
	if _, err := d.port.Write(req); err != nil {
		fmt.Errorf("rfid: error sending command: %s", err)
	}
	resp, err := rx(d.port)
	if err != nil {
		fmt.Errorf("rfid: error reading response: %s", err)
	}
	fmt.Println("RX:", resp)

}*/

func main() {

	// Open the serial port
	c := &serial.Config{Name: "COM3", Baud: 115200, ReadTimeout: 0}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	// load the reader
	scanner := bufio.NewScanner(s)
	scanner.Buffer(make([]byte, 16), 1)

	NodeId := make([]byte, 2)
	NodeId[0] = 0x00
	NodeId[1] = 0x00

	// tun on light blue
	request := newRequest(NodeId, CommandMifareRead, []byte{0x02})
	writerOfContent(s, request)

	wg.Add(1)
	go readerOfContent(scanner, &wg)
	time.Sleep(3 * time.Second)
	s.Flush()

	s.Close()
	wg.Wait()
}

func writerOfContent(s *serial.Port, request []byte) {
	fmt.Println("===========================================")
	fmt.Println("request:", request)
	fmt.Println("requesthex:", hex.EncodeToString(request))
	n, err := s.Write(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("n: ", n)
	fmt.Println("===========================================")
}

func readerOfContent(scanner *bufio.Scanner, wg *sync.WaitGroup) {
	i := 1
	condition := true
	for condition {
		if scanner.Scan() {
			fmt.Println("::::::::::::::::::::::::::::::::::::::::::")
			data := scanner.Bytes()
			fmt.Println(i, "responsehex: ", hex.EncodeToString(data))
			i++
			condition = scanner.Scan()
			fmt.Println("::::::::::::::::::::::::::::::::::::::::::")
		}
	}
	wg.Done()
}
