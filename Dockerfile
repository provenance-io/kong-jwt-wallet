FROM golang:1.17-alpine as build
RUN apk --no-cache --update add build-base && \
	go install github.com/Kong/go-pluginserver@v0.6.1
ADD go.mod go.sum jwt-wallet.go /app/
WORKDIR /app
RUN go build -buildmode plugin ./...

FROM kong:2.0.4-alpine
COPY --from=build /app/jwt-wallet.so /opt/go-plugins/
COPY --from=build /go/bin/go-pluginserver /usr/local/bin/
ADD config.yml /opt/
