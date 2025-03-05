package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"reflect"
)

const defaultListenAddr = ":4000"

type Config struct {
	ListenAddr string
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	ln        net.Listener
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	delPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message

	kv *KV
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	go s.loop()

	slog.Info("goredis server is running", "listenAddr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) handleMessage(msg Message) error {
	slog.Info("go message from client", "type", reflect.TypeOf(msg.cmd))

	switch v := msg.cmd.(type) {
	case SetCommand:
		if err := s.kv.Set(v.key, v.val); err != nil {
			return err
		}
		_, err := msg.peer.Send([]byte("OK\r\n"))
		if err != nil {
			return fmt.Errorf("peer send error %s", err)
		}
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		_, err := msg.peer.Send(val)
		if err != nil {
			return fmt.Errorf("peer send error %s", err)
		}
	case HelloCommand:
		spec := map[string]string{
			"server": "redis",
			"role":   "master",
		}
		_, err := msg.peer.Send(respWriteMap(spec))
		if err != nil {
			return fmt.Errorf("peer send error %s", err)
		}

		fmt.Println("this is the hello command from the client")
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Info("handle raw message error", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true

		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
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
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}
