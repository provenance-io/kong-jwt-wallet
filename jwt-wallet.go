package jwtwallet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"strings"

	"github.com/cosmos/btcutil/bech32"
	"github.com/provenance-io/kong-jwt-wallet/grants"
	"github.com/provenance-io/kong-jwt-wallet/signing"
	"golang.org/x/crypto/ripemd160"

	"github.com/Kong/go-pdk"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	RBAC         string `json:"rbac"`
	APIKey       string `json:"apikey"`
	AuthHeader   string `json:"authHeader"`
	AccessHeader string `json:"accessHeader"`
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

	if conf.AuthHeader == "" {
		conf.AuthHeader = "Authorization"
	}
	header, err := kong.Request.GetHeader(conf.AuthHeader)
	if err != nil {
		kong.Log.Warn("missing auth header")
		kong.Response.Exit(401, "{}", x)
		return
	}

	authToken := strings.Split(header, "Bearer")
	if len(authToken) < 2 {
		kong.Log.Warn("malformed auth header")
		kong.Response.Exit(401, "{}", x)
		return
	}

	tok, err := handleToken(kong, strings.TrimSpace(authToken[1]))
	if err != nil {
		kong.Log.Warn("err: " + err.Error())
		kong.Response.Exit(401, "{}", x)
		return
	}

	access, err := handleGrantedAccess(tok, conf.RBAC, conf.APIKey)
	if err != nil {
		kong.Log.Warn("err: " + err.Error())
		kong.Response.Exit(400, err.Error(), x)
		return
	}

	accessJson, err := json.Marshal(access)
	if err != nil {
		kong.Response.Exit(500, "something went wrong", x)
		return
	}
	if conf.AccessHeader == "" {
		conf.AccessHeader = "x-wallet-access"
	}
	kong.ServiceRequest.AddHeader(conf.AccessHeader, string(accessJson))

	kong.Log.Warn(tok)

}

var parser = jwt.NewParser()

func handleGrantedAccess(token *jwt.Token, url string, apiKey string) (*[]grants.GrantedAccess, error) {
	if claims, ok := token.Claims.(*signing.Claims); ok {
		if claims.Addr == "" {
			return nil, fmt.Errorf("missing addr claim")
		}
		if claims.Hrp == "" {
			return nil, fmt.Errorf("missing hrp claim")
		}

		if !verifyAddress(claims.Addr, claims.Subject, claims.Hrp) {
			return nil, fmt.Errorf("address does not match public key")
		}

		grantedAccess, err := grants.GetGrants(url, claims.Addr, apiKey)
		if err != nil {
			return nil, err
		}
		return grantedAccess, nil
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

func verifyAddress(addr string, pubKey string, hrp string) bool {
	keyB64 := strings.Split(pubKey, ",")[0]
	keyBytes, err := base64.RawURLEncoding.DecodeString(keyB64)

	if err != nil {
		fmt.Printf("Could not decode public key")
		return false
	}

	hash160Bytes := Hash160(keyBytes)

	dataBits, err := bech32.ConvertBits(hash160Bytes, 8, 5, true)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false
	}

	pubKeyAddr, err := bech32.Encode(hrp, dataBits)

	if err != nil {
		fmt.Printf("error: %v", err)
		return false
	}

	return strings.EqualFold(addr, pubKeyAddr)
}

// Calculate the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

// Hash160 calculates the hash ripemd160(sha256(b)).
func Hash160(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), ripemd160.New())
}
