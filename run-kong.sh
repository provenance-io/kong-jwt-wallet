#!/usr/bin/env bash

docker run -it --rm \
  -e "KONG_DATABASE=off" \
  -e "KONG_GO_PLUGINS_DIR=/opt/go-plugins" \
  -e "KONG_DECLARATIVE_CONFIG=/opt/config.yml" \
  -e "KONG_PLUGINS=jwt-wallet" \
  -e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
  -p 8000:8000 \
  kong-test
