package client

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"

	"github.com/golang/snappy"
)

const HANDSHAKE_REQ = '?'
const HANDSHAKE_RESP = '!'

func NewUdpClient(addr string) *UdpClient {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &UdpClient{
		server:     udpAddr,
		inputChan:  make(byteBufferChan, 5),
		outputChan: make(chan rune, 3),
	}
}

type UdpClient struct {
	server     *net.UDPAddr
	conn       *net.UDPConn
	input      []byte
	inputChan  byteBufferChan
	outputChan chan rune
}

func (c *UdpClient) Connect() error {
	var err error
	c.conn, err = net.DialUDP("udp", nil, c.server)

	if err != nil {
		return err
	}

	go handleUdpReading(c.conn, c.inputChan)
	go handleUdpWriting(c.conn, c.outputChan)

	for {
		c.Write(HANDSHAKE_REQ)

		if string(c.Read()) != "" {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

// Blocking or not, thats the question
// For now, we dont block
func (c *UdpClient) Read() []byte {
	select {
	case value := <-c.inputChan:
		c.input = value[0]
		return c.input
	default:
		return c.input
	}
}

func (c *UdpClient) Write(char rune) {
	select {
	// Try to write to the channel
	case c.outputChan <- char:
	// Otherwise clear channel
	default:
		<-c.outputChan
		c.outputChan <- char
	}
}

func handleUdpReading(conn net.Conn, inputChan byteBufferChan) {
	for {
		// Read Payload length
		var lengthBuffer [4]byte
		_, err := io.ReadFull(conn, lengthBuffer[:])
		if err != nil {
			conn.Close()
			return
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
			log.Println("Error decompressing data:", err)
			return
		}

		select {
		// Try to write to the channel
		case inputChan <- byteBuffer{decompressed}:
		// Otherwise clear channel
		default:
			<-inputChan
			inputChan <- byteBuffer{decompressed}
		}
	}
}

func handleUdpWriting(conn net.Conn, outputChan chan rune) {
	for {
		message := <-outputChan

		// Write back the same message to the client.
		conn.Write([]byte(string(message)))
	}
}
