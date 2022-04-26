# kong-wallet-jwt

Adds an extra layer of security and functions as a role authority. This plugin takes in an `Authorization` header with a user signed JWT token.  With a verified JWT this plugin can also function as a role authority and provide `x-wallet-access` that belong to the associated account that signed the JWT. 

## Getting started

When using this plugin you can use `go install github.com/provenance-io/kong-jwt-wallet/cmd/jwt-wallet@v0.3.0` directly or download a release version (soon to come)

### Configuration

When using the plugin, add it to your kong service definition and include an rbac url of choice. Currently the rbac url should contain an `{addr}` string target. 
```
  plugins:
  - name: jwt-wallet
    config:
      rbac: http://localhost:8888/{addr}/
```

### Running locally

```
make docker && make docker-run
```

## Creating a JWT

This example uses the standard jwt format but sings with an `secp256k1` elliptic curve key. When generating your jwt you must set the public key as the `sub` field on the payload and it must be compressed public key bytes (base64). If wanting grants to return then also include the wallet address as the `addr` field. 

### Header: 

```
{
  "alg": "ES256K",
  "typ": "JWT"
}
```

### Payload: 

```
{
  "addr": wallet_address,
  "sub": wallet_public_key,
  "iss": your_org,
  "iat": 1609459200,
  "exp": 4070908800
}
```

### Signature: 

```
ecdsa.Sign(
  SHA256(base64UrlEncode(header) + "." +
  base64UrlEncode(payload)))
```

Full token representation: 

```
base64UrlEncode(header) + "." +
  base64UrlEncode(payload) + "." +
    base64UrlEncode(signature)
```

