package jwtwallet

import (
	"fmt"

	"github.com/Kong/go-pdk"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	NetworkURI     string `json:"network_uri"`
	RBACServiceURI string `json:"rbac_service_uri"`
}

func New() interface{} {
	return &Config{}
}

func (conf *Config) Access(kong *pdk.PDK) {
	header, ok := kong.Request.GetHeader("Authorization")

	x := make(map[string][]string)
	x["Content-Type"] = append(x["Content-Type"], "application/json")
	if ok != nil {
		kong.Response.Exit(401, "Missing access token", x)
	}
	kong.Log.Warn(header)
	key, err := kong.Request.GetQueryArg("key")
	apiKey := conf.NetworkURI
	tmp := conf.RBACServiceURI
	if tmp == "" {
		kong.Log.Err("oopsies")
	}

	if err != nil {
		kong.Log.Err(err.Error())
	}

	if apiKey != key {
		kong.Response.Exit(403, "You have no correct key", x)
	}
}

func handleToken(tokenString string) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["exp"])
	} else {
		fmt.Println(err)
	}
}
