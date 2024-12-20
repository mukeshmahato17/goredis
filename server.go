package main

import (
	"fmt"
	"log"
	"net"
)

const defaultListenAddr = ":5000"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln net.Listener
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	log.Println("server running", "listenAddr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	fmt.Println("new incomming connection", conn)
}
