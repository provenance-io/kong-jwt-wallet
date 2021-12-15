package main

import "github.com/Kong/go-pdk/server"

var (
	Priority = 1
	Version  = "1.0.0"
)

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		panic(err)
	}
}
