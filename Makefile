all: bin

# lib: jwt-wallet.go
# 	go build -o jwt-wallet.so -buildmode plugin ./jwt-wallet.go

.PHONY: bin
bin:
	go build -o jwt-wallet ./cmd/jwt-wallet

.PHONY: lint
lint:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -d -s

.PHONY: http
http:
	python3 -m http.server 8888 --directory http/

.PHONY: docker
docker:
	docker build -t kong-test .

.PHONY: docker-run
docker-run:
	docker run --net host -it --name kong-test --rm \
		-e "KONG_DATABASE=off" \
		-e "KONG_LOG_LEVEL=debug" \
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-e "KONG_DECLARATIVE_CONFIG=/opt/config.yml" \
		-p 9000:8000 \
		kong-test:latest

.PHONY: docker-bash
docker-bash:
	docker run --net host -it --name kong-test --rm \
		-e "KONG_DATABASE=off" \
		-e "KONG_LOG_LEVEL=debug" \
		-e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
		-e "KONG_DECLARATIVE_CONFIG=/opt/config.yml" \
		-p 9000:8000 \
		kong-test:latest

.PHONY: curl
curl:
	# We only care about headers getting created
	curl -q -sSL -v -Hauthorization:$(shell cat ./token) localhost:8000/ >/dev/null

.PHONY: token
token:
	go run ./cmd/token/main.go

.PHONY: clean
clean:
	rm -f *.so jwt_wallet main main.so
