package main

import (
	"log"
	"log/slog"
	"net"
)

const defaultListenAddr = ":5000"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln        net.Listener
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan []byte
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	go s.loop()

	slog.Info("server running", "listenAddr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(string(rawMsg))
	if err != nil {
		return net.ErrClosed
	}
	switch v := cmd.(type) {
	case SetCommand:
		slog.Info("somebody want to set a key in the hash table", "key", v.key, "value", v.val)
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Info("handle raw message", "error", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
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
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}
