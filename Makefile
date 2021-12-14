all: lib

lib: jwt-wallet.go plugin/main.go
	go build -o jwt-wallet.so -buildmode plugin ./plugin/main.go

bin: jwt-wallet.go cmd/jwt_wallet/main.go
	go build -o jwt-wallet ./cmd/jwt_wallet/main.go

clean:
	rm -f *.so jwt_wallet main main.so
