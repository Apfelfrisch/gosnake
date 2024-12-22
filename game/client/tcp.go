package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/snappy"
)

func NewTcpClient(addr string) *Tcp {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Tcp{
		server:     tcpAddr,
		inputChan:  make(chan [1]string, 1),
		outputChan: make(chan rune, 3),
	}
}

type Tcp struct {
	server     *net.TCPAddr
	conn       *net.TCPConn
	input      string
	inputChan  chan [1]string
	outputChan chan rune
}

func (c *Tcp) Connect() error {
	var err error
	c.conn, err = net.DialTCP("tcp4", nil, c.server)

	if err != nil {
		return err
	}

	go handleClientReading(c.conn, c.inputChan)
	go handleClientWriting(c.conn, c.outputChan)

	return nil
}

// Blocking or not, thats the question
// For now, we dont block
func (c *Tcp) Read() string {
	select {
	case value := <-c.inputChan:
		c.input = value[0]
		return c.input
	default:
		return c.input
	}
}

func (c *Tcp) Write(char rune) {
	select {
	// Try to write to the channel
	case c.outputChan <- char:
	// Otherwise clear channel
	default:
		<-c.outputChan
		c.outputChan <- char
	}
}

func handleClientReading(conn net.Conn, inputChan chan [1]string) {
	for {
		// Read Payload length
		var lengthBuffer [4]byte
		_, err := io.ReadFull(conn, lengthBuffer[:])
		if err != nil {
			log.Println("Error reading length:", err)
			continue
		}

		// Read Payload
		compressed := make([]byte, binary.BigEndian.Uint32(lengthBuffer[:]))
		_, err = io.ReadFull(conn, compressed)
		if err != nil {
			log.Println(err)
			return
		}

		// Decompress Payload
		decompressed, err := snappy.Decode(nil, compressed)
		if err != nil {
			fmt.Println("Error decompressing data:", err)
			return
		}

		select {
		// Try to write to the channel
		case inputChan <- [1]string{string(decompressed)}:
		// Otherwise clear channel
		default:
			<-inputChan
			inputChan <- [1]string{string(decompressed)}
		}
	}
}

func handleClientWriting(conn net.Conn, outputChan chan rune) {
	for {
		message := <-outputChan

		// Write back the same message to the client.
		conn.Write([]byte(string(message)))
	}
}
