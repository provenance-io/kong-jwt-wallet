package signing

import (
	"encoding/base64"
	"fmt"
	"github.com/Kong/go-pdk"
	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

func init() {
	jwt.RegisterSigningMethod("ES256K", NewSecp256k1Signer)
}

func ParseKey(kong *pdk.PDK) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			if kong != nil {
				kong.Log.Warn("no claims")
			}
			return nil, fmt.Errorf("no claims")
		}
		sub := claims.Subject
		if sub == "" {
			if kong != nil {
				kong.Log.Warn("no subject")
			}
			return nil, fmt.Errorf("no subject")
		}
		keyB64 := strings.Split(sub, ",")[0]
		keyBytes, err := base64.RawURLEncoding.DecodeString(keyB64)
		if err != nil {
			return nil, err
		}
		pubk, err := secp256k1.ParsePubKey(keyBytes, secp256k1.S256())
		if err != nil {
			return nil, err
		}
		return pubk, nil
	}
}


type secp256k1Sig struct {
}

var _ jwt.SigningMethod = (*secp256k1Sig)(nil)

func (s secp256k1Sig) Verify(signingString, signature string, key interface{}) error {
	fmt.Printf("verify(" + signingString + "," + signature + ")")

	sigBytes, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	sig, err := secp256k1.ParseSignature(sigBytes, secp256k1.S256())
	if err != nil {
		return fmt.Errorf("sig parse failed: %w", err)
	}
	ok := sig.Verify([]byte(signingString), key.(*secp256k1.PublicKey))
	if !ok {
		return fmt.Errorf("sig verify failed")
	}
	return nil
}

func (s secp256k1Sig) Sign(signingString string, key interface{}) (string, error) {
	sig, err := key.(*secp256k1.PrivateKey).Sign([]byte(signingString))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sig.Serialize()), nil
}

func (s secp256k1Sig) Alg() string {
	return "ES256K"
}

func NewSecp256k1Signer() jwt.SigningMethod {
	return &secp256k1Sig{}
}
