package main

func main() {
	cfg := Config{
		ListenAddr: ":4000",
	}

	server := NewServer(cfg)
	server.Start()

}
