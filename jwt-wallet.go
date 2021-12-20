package jwtwallet

import (
	"encoding/json"
	"fmt"

	"github.com/provenance-io/kong-jwt-wallet/grants"
	"github.com/provenance-io/kong-jwt-wallet/signing"

	"github.com/Kong/go-pdk"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	RBAC string `json:"rbac"`
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

	grants, err := handleRoles(tok, conf.RBAC)
	if err != nil {
		kong.Log.Warn("err: " + err.Error())
		kong.Response.Exit(400, "account does not exist", x)
		return
	}

	grantsJson, err := json.Marshal(grants)
	if err != nil {
		kong.Response.Exit(500, "someting went wrong", x)
		return
	}
	kong.ServiceRequest.AddHeader("x-roles", string(grantsJson))
	//
	kong.Log.Warn(tok)

}

var parser = jwt.NewParser(jwt.WithoutClaimsValidation())

func handleRoles(token *jwt.Token, url string) (*grants.Grants, error) {
	fmt.Println(token.Claims.(*signing.Claims))
	if claims, ok := token.Claims.(*signing.Claims); ok {
		grants, err := grants.GetGrants(url, claims.Addr) // temporary interpolation until better configuration solutions
		if err != nil {
			return nil, err
		}
		return grants, nil
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
