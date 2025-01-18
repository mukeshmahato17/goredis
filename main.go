package main

import "flag"

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listen address of goredis server")
	flag.Parse()
	server := NewServer(Config{
		ListenAddr: *listenAddr,
	})
	server.Start()
}
