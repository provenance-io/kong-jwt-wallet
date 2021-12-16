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
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-e "KONG_DECLARATIVE_CONFIG=/opt/config.yml" \
		-p 9000:8000 \
		kong-test:latest

docker-bash:
	docker run --net host -it --name kong-test --rm \
		-e "KONG_DATABASE=off" \
		-e "KONG_LOG_LEVEL=debug" \
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-e "KONG_DECLARATIVE_CONFIG=/opt/config.yml" \
		-p 9000:8000 \
		kong-test:latest


clean:
	rm -f *.so jwt_wallet main main.so
