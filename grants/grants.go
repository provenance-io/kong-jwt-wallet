package grants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		AuthzGrants []interface{} `json:"authzGrants"`
		Apps        []interface{} `json:"apps"`
	} `json:"grants"`
}

func GetGrants(grantsUrl string, address string) ([]string, error) {
	client := &http.Client{}

	roleReq, _ := http.NewRequest("GET", grantsUrl, nil)
	roleReq.Header.Add("x-sender", address)

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

	roles := []string{}

	for _, grant := range roleResponse.Grants {
		roles = append(roles, grant.Roles...)
	}
	return roles, nil
}
