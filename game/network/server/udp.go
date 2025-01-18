package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"

	"github.com/golang/snappy"
)

const HANDSHAKE_REQ = '?'
const HANDSHAKE_RESP = '!'

func NewUdpSever(addr string, connCount int) *UdpServer {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &UdpServer{
		addr:        udpAddr,
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
	stopChan    chan struct{}
	isOpen      bool
}

func (s *UdpServer) Disconnect() {
	if s.isOpen {
		close(s.stopChan)
		s.isOpen = false
	}
	if s.conn != nil {
		log.Println("Connection closed")
		s.conn.Close()
		s.clients = []*net.UDPAddr{}
		s.inputChans = make(map[string]chan rune)
		s.outputChans = make(map[string]byteBufferChan)
	}
}

func (s *UdpServer) IsListining() bool {
	return s.isOpen
}

func (s *UdpServer) IsReady() bool {
	return len(s.clients) == s.clientCount
}

func (s *UdpServer) Listen(ctx context.Context) {
	var err error
	s.conn, err = net.ListenUDP("udp", s.addr)
	if err != nil {
		log.Fatal(err)
	}

	s.stopChan = make(chan struct{})
	s.isOpen = true

	buffer := make([]byte, 64)

	for len(s.clients) < s.clientCount {
		select {
		case <-ctx.Done():
			s.Disconnect()
			return
		default:
			if err := s.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
				log.Fatal("UDP-SERVER:" + err.Error())
			}
			n, clientAddr, err := s.conn.ReadFromUDP(buffer)

			if err != nil {
				if netErr, ok := err.(*net.OpError); ok && !netErr.Timeout() {
					log.Printf("UDP-SERVER:" + err.Error())
				}
				continue
			}

			if string(buffer[:n]) != string(HANDSHAKE_REQ) {
				log.Printf("UDP-SERVER:" + fmt.Sprintf("Invalid Handshake: [%s]", string(buffer[:n])))
				continue
			}

			if _, ok := s.inputChans[clientAddr.String()]; ok {
				continue
			}

			s.clients = append(s.clients, clientAddr)
			s.inputChans[clientAddr.String()] = make(chan rune, 3)
			s.outputChans[clientAddr.String()] = make(byteBufferChan, 1)

			go s.handleServerWriting(clientAddr, s.outputChans[clientAddr.String()])
		}
	}

	// Start Server Reading after s.addClient
	// otherwise we get Deadlock
	go s.handleSeverReading()
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

func (s *UdpServer) addClient() (*net.UDPAddr, error) {
	buffer := make([]byte, 64)
	for {
		n, clientAddr, err := s.conn.ReadFromUDP(buffer)

		if err != nil {
			return clientAddr, err
		}

		if string(buffer[:n]) != string(HANDSHAKE_REQ) {
			return clientAddr, fmt.Errorf("Invalid Handshake: [%s]", string(buffer[:n]))
		}

		if _, ok := s.inputChans[clientAddr.String()]; ok {
			continue
		}

		s.clients = append(s.clients, clientAddr)
		s.inputChans[clientAddr.String()] = make(chan rune, 3)
		s.outputChans[clientAddr.String()] = make(byteBufferChan, 1)

		return clientAddr, nil
	}
}

func (s *UdpServer) handleSeverReading() {
	if err := s.conn.SetReadDeadline(time.Time{}); err != nil {
		log.Fatal("UDP-SERVER:" + err.Error())
	}

	buffer := make([]byte, 16)
	readConnection := func() {
		n, remoteAddr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := buffer[:n]
		for len(data) == 0 {
			return
		}

		inputChan, ok := s.inputChans[remoteAddr.String()]
		if !ok {
			return
		}

		rune, size := utf8.DecodeRune(data)

		if rune == utf8.RuneError && size == 1 {
			return
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

	for {
		select {
		case <-s.stopChan:
			return
		default:
			readConnection()
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

	for {
		select {
		case <-s.stopChan:
			return
		case message := <-outputChan:
			writemessage(message[0])
		}
	}
}
