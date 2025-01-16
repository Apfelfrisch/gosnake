package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"unicode/utf8"

	"github.com/golang/snappy"
)

const HANDSHAKE_REQ = '?'
const HANDSHAKE_RESP = '!'

func NewUdpSever(addr string, connCount int) *UdpServer {
	tcpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &UdpServer{
		addr:        tcpAddr,
		clientCount: connCount,
		inputChans:  make(map[string]chan rune),
		outputChans: make(map[string]byteBufferChan),
	}
}

type UdpServer struct {
	addr        *net.UDPAddr
	conn        *net.UDPConn
	clients     []*net.UDPAddr
	clientCount int
	inputChans  map[string]chan rune
	outputChans map[string]byteBufferChan
}

func (s *UdpServer) Ready() bool {
	return len(s.clients) == s.clientCount
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
	var err error
	s.conn, err = net.ListenUDP("udp", s.addr)
	if err != nil {
		log.Fatal(err)
	}

	for len(s.clients) < s.clientCount {
		clientAddr, err := s.addClient()

		if err != nil {
			log.Printf("UDP-SERVER:" + err.Error())
			continue
		}

		go s.handleServerWriting(clientAddr, s.outputChans[clientAddr.String()])
	}

	// Start Server Reading after s.addClient
	// otherwise we get Deadlock
	go s.handleSeverReading()
}

func (s *UdpServer) addClient() (*net.UDPAddr, error) {
	buffer := make([]byte, 64)
	n, clientAddr, err := s.conn.ReadFromUDP(buffer)

	if err != nil {
		return clientAddr, err
	}

	if string(buffer[:n]) != string(HANDSHAKE_REQ) {
		return clientAddr, fmt.Errorf("Invalid Handshake: [%s]", string(buffer[:n]))
	}

	if _, ok := s.inputChans[clientAddr.String()]; ok {
		return clientAddr, fmt.Errorf("Client already connected [%s]", clientAddr.String())
	}

	s.clients = append(s.clients, clientAddr)
	s.inputChans[clientAddr.String()] = make(chan rune, 3)
	s.outputChans[clientAddr.String()] = make(byteBufferChan, 1)

	return clientAddr, nil
}

func (s *UdpServer) handleSeverReading() {
	buffer := make([]byte, 16)
	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
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

		if rune == HANDSHAKE_REQ {
			s.WriteConn(remoteAddr, []byte(string(HANDSHAKE_RESP)))
		}

		select {
		case inputChan <- rune:
		default:
			<-inputChan
			inputChan <- rune
		}
	}
}

func (s *UdpServer) handleServerWriting(clientAddr *net.UDPAddr, outputChan byteBufferChan) {
	writemessage := func(message []byte) {
		compressed := snappy.Encode(nil, message)

		var lengthBuffer bytes.Buffer
		binary.Write(&lengthBuffer, binary.BigEndian, uint32(len(compressed)))

		s.conn.WriteToUDP(lengthBuffer.Bytes(), clientAddr)
		s.conn.WriteToUDP(compressed, clientAddr)
	}

	writemessage([]byte(string(HANDSHAKE_RESP)))

	for {
		message := <-outputChan
		writemessage(message[0])
	}
}
