package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/provenance-io/jwt-wallet/signing"
	"time"

	"github.com/Kong/go-pdk"
	secp256k1 "github.com/btcsuite/btcd/btcec"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	NetworkURI     string `json:"network_uri"`
	RBACServiceURI string `json:"rbac_service_uri"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	defer func() {
		err := recover()
		if err != nil {
			e, ok := err.(error)
			if ok {
				kong.Log.Err(e.Error())
				kong.Response.Exit(500, "{}", map[string][]string{})
			} else {
				kong.Log.Err(fmt.Sprintf("%v", err))
				kong.Response.Exit(500, "{}", map[string][]string{})
			}
		}
	}()

	x := make(map[string][]string)
	x["Content-Type"] = []string{"application/json"}

	header, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		kong.Log.Warn("missing auth header")
		kong.Response.Exit(401, "{}", x)
		return
	}

	tok, err := handleToken(kong, header)
	if err != nil {
		kong.Log.Warn("err:" + err.Error())
		kong.Response.Exit(401, "{}", x)
		return
	}
	//
	kong.Log.Warn(tok)
	kong.Response.Exit(200, "{}", x)
	return
}




var parser = jwt.NewParser(jwt.WithoutClaimsValidation())

func handleToken(kong *pdk.PDK, tokenString string) (*jwt.Token, error) {
	var claims jwt.RegisteredClaims
	token, err := parser.ParseWithClaims(tokenString, &claims, signing.ParseKey(kong))
	if err != nil {
		if kong != nil {
			kong.Log.Warn("parse error:" + err.Error())
		}
		return nil, err
	}
	return token, nil
}

func main() {
	pkBytes, err := hex.DecodeString("8C037EFC21AB3F0F8D32CF209D90FDBF41D10071FF600BA66A30EFA994F268A3")
	if err != nil {
		panic(err)
	}
	prvk, pubk := secp256k1.PrivKeyFromBytes(secp256k1.S256(), pkBytes)
	fmt.Printf("prvKey:%X\n", prvk.Serialize())
	fmt.Printf("pubKey:%X\n", pubk.SerializeCompressed())

	loc, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}

	preClaims := jwt.RegisteredClaims{}
	preClaims.Subject = base64.RawURLEncoding.EncodeToString(pubk.SerializeCompressed())
	preClaims.Issuer = "provenance.io"
	preClaims.IssuedAt = jwt.NewNumericDate(time.Date(2021, 1, 1, 0, 0, 0, 0, loc))
	preClaims.ExpiresAt = jwt.NewNumericDate(time.Date(2099, 1, 1, 0, 0, 0, 0, loc))

	token := jwt.NewWithClaims(signing.NewSecp256k1Signer(), preClaims)
	sig, err := token.SignedString(prvk)
	fmt.Printf("signed:%s\n", sig)

	newClaims := jwt.RegisteredClaims{}
	newToken, err := jwt.ParseWithClaims(sig, &newClaims, signing.ParseKey(nil))
	fmt.Printf("sig:%s\n", newToken.Signature)
	fmt.Printf("valid:%+v\n", newToken)
}

