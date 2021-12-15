all: lib bin

lib: jwt-wallet.go
	go build -o jwt-wallet.so -buildmode plugin ./jwt-wallet.go

bin: jwt-wallet.go main.go
	go build -o jwt-wallet ./main.go ./jwt-wallet.go

lint:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -d -s

docker:
	docker build -t kong-test .

docker-run:
	docker run -it --name kong-test --rm \
		-v $(CURDIR):/opt/go-plugins \
		-e "KONG_DATABASE=off" \
		-e "KONG_GO_PLUGINS_DIR=/opt/go-plugins" \
		-e "KONG_DECLARATIVE_CONFIG=/opt/config.yml" \
		-e "KONG_PLUGINS=jwt-wallet" \
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-p 9000:8000 \
		kong-test


clean:
	rm -f *.so jwt_wallet main main.so
