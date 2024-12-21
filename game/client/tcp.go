package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
		// Read from the connection untill a new line is send
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		select {
		// Try to write to the channel
		case inputChan <- [1]string{data}:
		// Otherwise clear channel
		default:
			<-inputChan
			inputChan <- [1]string{data}
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
