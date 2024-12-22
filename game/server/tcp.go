package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"

	"github.com/golang/snappy"
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
	conns         tcpConns
}

func (s *Tcp) Ready() bool {
	return len(s.inputChans) == s.conns.count()
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
		s.conns.add(conn)

		if err != nil {
			log.Fatal("Could not connect: " + err.Error())
		}

		// Handle new connections in a Goroutine for concurrency
		go s.handleSeverReading(conn, inputChan)
	}

	go s.handleServerWriting(s.broadcastChan)
}

func (s *Tcp) Shutdown() {
	s.conns.reset()
}

func (s *Tcp) handleSeverReading(conn net.Conn, inputChan chan rune) {
	for {
		// Read from the connection untill a new line is send
		data, _, err := bufio.NewReader(conn).ReadRune()
		if err != nil {
			s.Shutdown()
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

func (s *Tcp) handleServerWriting(broadcastChan chan [1]string) {
	for {
		message := <-broadcastChan
		compressed := snappy.Encode(nil, []byte(message[0]))

		var lengthBuffer bytes.Buffer
		binary.Write(&lengthBuffer, binary.BigEndian, uint32(len(compressed)))

		for _, conn := range s.conns.get() {
			conn.Write(lengthBuffer.Bytes())
			conn.Write(compressed)
		}
	}
}

type tcpConns struct {
	conns []net.Conn
	mu    sync.Mutex
}

func (tc *tcpConns) add(conn net.Conn) {
	tc.mu.Lock()
	tc.conns = append(tc.conns, conn)
	tc.mu.Unlock()
}

func (tc *tcpConns) count() int {
	return len(tc.conns)
}

func (tc *tcpConns) get() []net.Conn {
	return tc.conns
}

func (tc *tcpConns) reset() {
	tc.mu.Lock()
	for _, conn := range tc.conns {
		conn.Close()
	}
	tc.conns = []net.Conn{}
	tc.mu.Unlock()
}
