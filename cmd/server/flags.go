package main

import "flag"

var addr string

func parseFlag() {
	flag.StringVar(&addr, "a", ":8080", "server address")
	flag.Parse()
}
