package main

type ServerConfig struct {
	Addr          string
	StoreInterval int64
	Path          string
	Restore       bool
	LogLevel      string
}
