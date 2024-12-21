package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func NewTcpSever(addr string, connCount int) *Tcp {
	var inputChans []chan rune

	for i := 0; i < connCount; i++ {
		inputChans = append(inputChans, make(chan rune, 3))
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Tcp{
		addr:          tcpAddr,
		inputChans:    inputChans,
		broadcastChan: make(chan [1]string, 1),
	}
}

type Tcp struct {
	addr          *net.TCPAddr
	inputChans    []chan rune
	broadcastChan chan [1]string
	conns         []net.Conn
}

func (s *Tcp) Ready() bool {
	return len(s.inputChans) == len(s.conns)
}

func (s *Tcp) ReadConn(connIndex int) *rune {
	select {
	case value := <-s.inputChans[connIndex]:
		return &value
	default:
		return nil
	}
}

func (s *Tcp) Broadcast(content string) {
	select {
	// Try to write to the channel
	case s.broadcastChan <- [1]string{content + "\n"}:
	// Otherwise clear channel
	default:
		<-s.broadcastChan
		s.broadcastChan <- [1]string{content + "\n"}
	}
}

func (s *Tcp) Listen() {
	listener, err := net.ListenTCP("tcp", s.addr)
	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, inputChan := range s.inputChans {
		// Accept new connections
		conn, err := listener.Accept()
		s.conns = append(s.conns, conn)

		if err != nil {
			log.Fatal("Could not connect: " + err.Error())
		}
		// Handle new connections in a Goroutine for concurrency
		go handleSeverReading(conn, inputChan)
	}

	go handleServerWriting(s.conns, s.broadcastChan)
}

func (s *Tcp) Shutdown() {
	for _, conn := range s.conns {
		conn.Close()
	}
}

func handleSeverReading(conn net.Conn, inputChan chan rune) {
	for {
		// Read from the connection untill a new line is send
		data, _, err := bufio.NewReader(conn).ReadRune()
		if err != nil {
			fmt.Println(err)
			return
		}

		select {
		case inputChan <- data:
		default:
			<-inputChan
			inputChan <- data
		}
	}
}

func handleServerWriting(conns []net.Conn, broadcastChan chan [1]string) {
	for {
		message := <-broadcastChan

		for _, conn := range conns {
			conn.Write([]byte(message[0]))
		}
	}
}
