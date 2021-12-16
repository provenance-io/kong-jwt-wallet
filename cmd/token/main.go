package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/golang-jwt/jwt/v4"
	"github.com/provenance-io/jwt-wallet/signing"
	"time"
)

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

	//testToken := "eyJhbGciOiJFUzI1NksiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJwcm92ZW5hbmNlLmlvIiwic3ViIjoiQTQ1SndKcEptX2J1UW5vd2Z6SE9RQk5oYTRpNEQ5cUYyOUZPVnQ3NGlqQ1UiLCJpYXQiOiIxNjM5NjAyNzIwMDAwIiwiZXhwIjoiMTYzOTc4MjcyMDAwMCIsImFkZHIiOiJ0cDFmeWVkZmVnemd3ODhxZHg0MHh6NXBxd2podHQ2Zmo2ZHc2eHBwYSJ9.JCThW-MlX_3nCrketHDEqoGFPf_59nWAOMFMW_38UmlUECc0fLq-bBjHP9z1yrEkkMHpG3Kh_psRpFp2k4mqQA"

	//signing.NewSecp256k1Signer().Verify(signingStringTest, {})

	var newClaims signing.Claims
	newToken, err := jwt.ParseWithClaims(sig, &newClaims, signing.ParseKey(nil))

	//roles, err := handleRoles(token, "http://localhost:8080/api/v1/rbac/account/")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("\nroles:%s\n", roles)
	fmt.Printf("sig:%s\n", newToken.Signature)
	fmt.Printf("valid:%+v\n", newToken)
}
