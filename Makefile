all: lib bin

lib: jwt-wallet.go
	go build -o jwt-wallet.so -buildmode plugin ./jwt-wallet.go

bin: jwt-wallet.go main.go
	go build -o jwt-wallet ./jwt-wallet.go

lint:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -d -s

.PHONY: http
http:
	python3 -m http.server 8888 --directory http/

docker:
	docker build -t kong-test .

docker-run:
	docker run --net host -it --name kong-test --rm \
		-e "KONG_DATABASE=off" \
		-e "KONG_LOG_LEVEL=debug" \
		-e "KONG_PLUGINSERVER_NAMES=jwt-wallet" \
		-e "KONG_PLUGINSERVER_JWT_WALLET_START_CMD=/usr/local/bin/jwt-wallet" \
		-e "KONG_PLUGINSERVER_JWT_WALLET_QUERY_CMD=/usr/local/bin/jwt-wallet -dump" \
		-e "KONG_PLUGINS=bundled,jwt-wallet" \
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-p 9000:8000 \
		kong:latest

docker-bash:
	docker run --net host -it --name kong-test --rm \
		-e "KONG_DATABASE=off" \
		-e "KONG_LOG_LEVEL=debug" \
		-e "KONG_PLUGINSERVER_NAMES=jwt-wallet" \
		-e "KONG_PLUGINSERVER_JWT_WALLET_START_CMD=/usr/local/bin/jwt-wallet" \
		-e "KONG_PLUGINSERVER_JWT_WALLET_QUERY_CMD=/usr/local/bin/jwt-wallet -dump" \
		-e "KONG_PLUGINS=bundled,jwt-wallet" \
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-p 9000:8000 \
		kong:latest


clean:
	rm -f *.so jwt_wallet main main.so
