package grants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SubjectResponse struct {
	Address string          `json:"address"`
	Name    string          `json:"name"`
	Grants  []GrantedAccess `json:"grants"`
}

type GrantedAccess struct {
	Address      string `json:"address"`
	Name         string `json:"name"`
	Applications []struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	} `json:"applications"`
}

var (
	Client HTTPClient
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	Client = &http.Client{}
}

func GetGrants(grantsURL string, address string, apiKey string) (*SubjectResponse, error) {
	uri := strings.ReplaceAll(grantsURL, "{addr}", address)
	request, _ := http.NewRequest("GET", uri, nil)
	request.Header.Add("x-sender", address)
	// Add apikey if supplied.
	if apiKey != "" {
		request.Header.Add("apikey", apiKey)
	}
	resp, err := Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response SubjectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	return &response, nil
}
