package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
	"unicode/utf8"

	"github.com/golang/snappy"
)

const HANDSHAKE_RESP = '!'

func NewUdpSever(addr string, connCount int) *UdpServer {
	tcpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &UdpServer{
		addr:        tcpAddr,
		connCount:   connCount,
		inputChans:  make(map[string]chan rune),
		outputChans: make(map[string]byteBufferChan),
	}
}

type UdpServer struct {
	addr        *net.UDPAddr
	conns       []*net.UDPAddr
	connCount   int
	inputChans  map[string]chan rune
	outputChans map[string]byteBufferChan
}

func (s *UdpServer) Ready() bool {
	return len(s.conns) == s.connCount
}

func (s *UdpServer) ReadConn(addr *net.UDPAddr) *rune {
	select {
	case value := <-s.inputChans[addr.String()]:
		return &value
	default:
		return nil
	}
}

func (s *UdpServer) WriteConn(addr *net.UDPAddr, content []byte) {

	select {
	// Try to write to the channel
	case s.outputChans[addr.String()] <- byteBuffer{content}:
	// Otherwise clear channel
	default:
		<-s.outputChans[addr.String()]
		s.outputChans[addr.String()] <- byteBuffer{content}
	}
}

func (s *UdpServer) Listen() {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		log.Fatal(err)
	}

	for len(s.conns) < s.connCount {
		var buf [64]byte
		_, clientAddr, err := conn.ReadFromUDP(buf[0:])

		if err != nil {
			log.Print(err)
			continue
		}

		if _, ok := s.inputChans[clientAddr.String()]; ok {
			continue
		}

		s.conns = append(s.conns, clientAddr)
		s.inputChans[clientAddr.String()] = make(chan rune, 3)
		s.outputChans[clientAddr.String()] = make(byteBufferChan, 1)

		go s.handleServerWriting(conn, clientAddr, s.outputChans[clientAddr.String()])
	}

	go s.handleSeverReading(conn)
}

func (s *UdpServer) handleSeverReading(conn *net.UDPConn) {
	buffer := make([]byte, 16)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		data := buffer[:n]
		for len(data) == 0 {
			continue
		}

		inputChan, ok := s.inputChans[remoteAddr.String()]
		if !ok {
			continue
		}

		rune, size := utf8.DecodeRune(data)

		if rune == utf8.RuneError && size == 1 {
			continue
		}

		select {
		case inputChan <- rune:
		default:
			<-inputChan
			inputChan <- rune
		}
	}
}

func (s *UdpServer) handleServerWriting(conn *net.UDPConn, clientAddr *net.UDPAddr, outputChan byteBufferChan) {
	writemessage := func(message []byte) {
		compressed := snappy.Encode(nil, message)

		var lengthBuffer bytes.Buffer
		binary.Write(&lengthBuffer, binary.BigEndian, uint32(len(compressed)))

		conn.WriteToUDP(lengthBuffer.Bytes(), clientAddr)
		conn.WriteToUDP(compressed, clientAddr)
	}

	writemessage([]byte(string(HANDSHAKE_RESP)))

	time.Sleep(time.Second)

	for {
		message := <-outputChan
		writemessage(message[0])
	}
}
