package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/provenance-io/jwt-wallet/grants"
	"github.com/provenance-io/jwt-wallet/signing"

	"github.com/Kong/go-pdk"
	secp256k1 "github.com/btcsuite/btcd/btcec"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	GRPCURL  string `json:"grpc_url"`
	RolesURL string `json:"roles_url"`
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

	roles, err := handleRoles(tok, conf.RolesURL)
	if err != nil {
		kong.Log.Warn("err: " + err.Error())
		return
	}

	kong.ServiceRequest.AddHeader("x-roles", strings.Join(roles, ","))

	//
	kong.Log.Warn(tok)
	kong.Response.Exit(200, "{}", x)
	return
}

var parser = jwt.NewParser(jwt.WithoutClaimsValidation())

func handleRoles(token *jwt.Token, url string) ([]string, error) {
	fmt.Println(token.Claims.(*signing.Claims))
	if claims, ok := token.Claims.(*signing.Claims); ok {
		addr := claims.Addr

		addrString := fmt.Sprintf("%v", addr)
		roles, err := grants.GetGrants(url+addrString+"/grants", addrString) // temporary interpolation until better configuration solutions
		if err != nil {
			return nil, err
		}
		return roles, nil
	}
	return nil, fmt.Errorf("malformed claims")
}

func handleToken(kong *pdk.PDK, tokenString string) (*jwt.Token, error) {
	var claims signing.Claims
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

	claims := &signing.Claims{
		Addr: "tp1uz5g72pvfrdnm9qnjpyvsnwc64d4wygyqanx2t",
		RegisteredClaims: *&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Date(2099, 1, 1, 0, 0, 0, 0, loc)),
			IssuedAt:  jwt.NewNumericDate(time.Date(2021, 1, 1, 0, 0, 0, 0, loc)),
			Issuer:    "provenance.io",
			Subject:   base64.RawURLEncoding.EncodeToString(pubk.SerializeCompressed()),
		},
	}

	token := jwt.NewWithClaims(signing.NewSecp256k1Signer(), claims)
	sig, err := token.SignedString(prvk)
	fmt.Printf("signed:%s\n", sig)

	var newClaims signing.Claims
	newToken, err := jwt.ParseWithClaims(sig, &newClaims, signing.ParseKey(nil))

	roles, err := handleRoles(token, "http://localhost:8080/api/v1/rbac/account/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\nroles:%s\n", roles)
	fmt.Printf("sig:%s\n", newToken.Signature)
	fmt.Printf("valid:%+v\n", newToken)
}
