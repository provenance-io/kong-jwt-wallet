all: lib

lib: jwt-wallet.go
	go build -buildmode plugin jwt-wallet.go
