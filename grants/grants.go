package grants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RoleResponse struct {
	Account struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Type    string `json:"type"`
	} `json:"account"`
	Grants []struct {
		Org struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Type    string `json:"type"`
		} `json:"org"`
		Roles       []string      `json:"roles"`
		AuthzGrants []string      `json:"authzGrants"`
		Apps        []interface{} `json:"apps"`
	} `json:"grants"`
}

type Grants struct {
	Orgs []Org `json:"orgs"`
}

type Org struct {
	Name        string   `json:"name"`
	Roles       []string `json:"roles"`
	AuthzGrants []string `json:"authzGrants"`
}

func GetGrants(grantsURL, address, apiKey string) (*Grants, error) {
	client := &http.Client{}

	uri := strings.ReplaceAll(grantsURL, "{addr}", address)
	roleReq, _ := http.NewRequest("GET", uri, nil)
	roleReq.Header.Add("x-sender", address)
	// Add apikey if supplied.
	if apiKey != "" {
		roleReq.Header.Add("apikey", apiKey)
	}
	resp, err := client.Do(roleReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var roleResponse RoleResponse
	if err := json.Unmarshal(body, &roleResponse); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	var grants Grants
	for _, grant := range roleResponse.Grants {
		org := Org{
			Name:        grant.Org.Name,
			Roles:       grant.Roles,
			AuthzGrants: grant.AuthzGrants,
		}

		grants.Orgs = append(grants.Orgs, org)
	}
	return &grants, nil
}
