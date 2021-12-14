all: lib bin

lib: jwt-wallet.go plugin/main.go
	go build -o jwt-wallet.so -buildmode plugin ./plugin/main.go

bin: jwt-wallet.go cmd/jwt_wallet/main.go
	go build -o jwt-wallet ./cmd/jwt_wallet/main.go

docker:
	docker build -t kong-test .

docker-run:
	docker run -it --rm \
		-e "KONG_DATABASE=off" \
		-e "KONG_GO_PLUGINS_DIR=/opt/go-plugins" \
		-e "KONG_DECLARATIVE_CONFIG=/opt/config.yml" \
		-e "KONG_PLUGINS=jwt-wallet" \
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-p 8000:8000 \
		kong-test


clean:
	rm -f *.so jwt_wallet main main.so
