package main

import (
	"github.com/Kong/go-pdk/server"
	jwtWallet "github.com/provenance-io/jwt-wallet"
)

func main() {
	err := server.StartServer(jwtWallet.New, "", 1)
	if err != nil {
		panic(err)
	}
}
