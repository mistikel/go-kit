package main

import (
	"github.com/mistikel/go-kit/server"
)

func main() {
	s := server.New()
	s.Serve()
}
