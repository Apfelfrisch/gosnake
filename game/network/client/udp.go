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
	server        *net.UDPAddr
	conn          *net.UDPConn
	input         []byte
	lastHandshake time.Time
	inputChan     byteBufferChan
	outputChan    chan rune
}

func (c *UdpClient) Connect() error {
	var err error
	c.conn, err = net.DialUDP("udp", nil, c.server)

	if err != nil {
		return err
	}

	go c.handleUdpReading()
	go c.handleUdpWriting()

	beforeHandshare := time.Now()
	for {
		c.Write(HANDSHAKE_REQ)

		time.Sleep(time.Second / 5)

		if c.lastHandshake.After(beforeHandshare) {
			break
		}
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

func (c *UdpClient) handleUdpReading() {
	for {
		var lengthBuffer [4]byte
		n, err := io.ReadFull(c.conn, lengthBuffer[:])
		if len(lengthBuffer[:n]) == 0 {
			continue
		}
		if err != nil {
			log.Println("UDP-CLIENT: Could not read response: ", string(lengthBuffer[:]))
			continue
		}

		// Read Payload
		compressed := make([]byte, binary.BigEndian.Uint32(lengthBuffer[:]))
		_, err = io.ReadFull(c.conn, compressed)
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

		if len(decompressed) == 1 && string(decompressed) == string(HANDSHAKE_RESP) {
			c.lastHandshake = time.Now()
			continue
		}

		select {
		// Try to write to the channel
		case c.inputChan <- byteBuffer{decompressed}:
		// Otherwise clear channel
		default:
			<-c.inputChan
			c.inputChan <- byteBuffer{decompressed}
		}
	}
}

func (c *UdpClient) handleUdpWriting() {
	for {
		message := <-c.outputChan

		c.conn.Write([]byte(string(message)))
	}
}
