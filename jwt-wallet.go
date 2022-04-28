package jwtwallet

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/provenance-io/kong-jwt-wallet/grants"
	"github.com/provenance-io/kong-jwt-wallet/signing"

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
		kong.Log.Warn("err:" + err.Error())
		kong.Response.Exit(401, "{}", x)
		return
	}

	access, err := handleGrantedAccess(tok, conf.RBAC, conf.APIKey)
	if err != nil {
		kong.Log.Warn("err: " + err.Error())
		kong.Response.Exit(400, "account does not exist", x)
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
