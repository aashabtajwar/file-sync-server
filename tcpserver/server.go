package tcpserver

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type TcpFServer struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func NewServerTwo() *TcpFServer {
	s := &TcpFServer{
		quit: make(chan interface{}),
	}
	fmt.Println("turning on tcp server...")
	l, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tcp server is on!")
	s.listener = l
	s.wg.Add(1)
	return s
}

func NewServer(addr string) *TcpFServer {
	s := &TcpFServer{
		quit: make(chan interface{}),
	}
	fmt.Println("turning on tcp server...")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tcp server is on!")
	s.listener = l
	s.wg.Add(1)
	s.serve()
	return s
}

func (s *TcpFServer) Stop() {
	fmt.Println("stopping tcp server...")
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *TcpFServer) serve() {
	fmt.Println("accepting connections!")
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				fmt.Println("Connection found", conn)
				s.handleConnection(conn)
				fmt.Println("done with this...	")
				s.wg.Done()
			}()
		}
	}
}

func (s *TcpFServer) handleConnection(conn net.Conn) {
	activeConnections = append(activeConnections, conn)
	conns = append(conns, conn)
	for _, e := range activeConnections {
		fmt.Println("Connection: ", e)

	}
	fmt.Println("New Connection: ", conn)
	s.CheckReceivedData(conn, activeConnections)
}

var server []net.Listener
var activeConnections []net.Conn

var conns []net.Conn

func updatedConnections() []net.Conn {
	return conns
}

func SetupConn() net.Listener {
	if len(server) == 1 {
		return server[0]
	}
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Turning on TCP Server\n", ln)
	server = append(server, ln)
	return ln
}

func Start() {
	// ln, err := net.Listen("tcp", ":3030")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// msg := make(chan string)

	// s := NewServer(":3030")

	// s.Stop()

	// old code here
	// s := NewServerTwo()

	s := &TcpFServer{
		quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", "192.168.0.103:3030")
	if err != nil {
		log.Fatal(err)
	}
	s.listener = l
	s.wg.Add(1)
	// ln := SetupConn()

	for {
		conn, err := s.listener.Accept()
		activeConnections = append(activeConnections, conn)
		conns = append(conns, conn)
		for _, e := range activeConnections {
			fmt.Println("connection: ", e)
		}
		fmt.Println("New Connection: ", conn)
		fmt.Println("CC = ", connections)
		if err != nil {
			log.Fatal(err)
		}
		go s.CheckReceivedData(conn, activeConnections)
	}

	// old code ends here
}
