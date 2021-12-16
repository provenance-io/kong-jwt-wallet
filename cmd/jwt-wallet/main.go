package main

import (
	"github.com/Kong/go-pdk/server"
	jwtwallet "github.com/provenance-io/jwt-wallet"
)

var (
	Priority = 1
	Version  = "1.0.0"
)

func main() {
	err := server.StartServer(jwtwallet.New, Version, Priority)
	if err != nil {
		panic(err)
	}
}
