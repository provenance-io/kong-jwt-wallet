FROM kong/go-plugin-tool:2.0.4-alpine-latest AS builder

RUN mkdir -p /tmp/jwt-wallet/

COPY . /tmp/jwt-wallet/

RUN go version

RUN cd /tmp/jwt-wallet/ && \
    go get github.com/Kong/go-pdk && \
    go mod init kong-go-plugin && \
    go get -d -v github.com/Kong/go-pluginserver && \
    go build github.com/Kong/go-pluginserver && \
    go build -buildmode plugin jwt-wallet.go

FROM kong:2.0.4-alpine

RUN mkdir /tmp/go-plugins
COPY --from=builder  /tmp/jwt-wallet/go-pluginserver /usr/local/bin/go-pluginserver
COPY --from=builder  /tmp/jwt-wallet/jwt-wallet.so /tmp/go-plugins
COPY config.yml /tmp/config.yml

USER root
RUN chmod -R 777 /tmp
RUN /usr/local/bin/go-pluginserver -version && \
    cd /tmp/go-plugins && \
    /usr/local/bin/go-pluginserver -dump-plugin-info jwt-wallet
USER kong