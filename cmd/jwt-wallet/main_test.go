package main

import (
	"net/http"
	"testing"

	"github.com/Kong/go-pdk/test"
	jwtwallet "github.com/provenance-io/kong-jwt-wallet"
	"github.com/stretchr/testify/assert"
)

func TestInvalidJwt(t *testing.T) {
	env, err := test.New(t, test.Request{
		Method:  "GET",
		Url:     "http://example.com",
		Headers: map[string][]string{"Authorization": {"not-a-valid-jwt"}},
	})
	assert.NoError(t, err)

	config := &jwtwallet.Config{}
	env.DoHttp(config)

	assert.Equal(t, 401, env.ClientRes.Status)
}

func TestMissingAddrClaim(t *testing.T) {
	env, err := test.New(t, test.Request{
		Method:  "GET",
		Url:     "http://example.com",
		Headers: map[string][]string{"Authorization": {"eyJhbGciOiJFUzI1NksiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJwcm92ZW5hbmNlLmlvIiwic3ViIjoiQTQ1SndKcEptX2J1UW5vd2Z6SE9RQk5oYTRpNEQ5cUYyOUZPVnQ3NGlqQ1UiLCJpYXQiOiIxNjQwMDMzOTU2MDAwIiwiZXhwIjoiMTY0MDIxMzk1NjAwMCJ9.jKKH5C7dl_fv7MHNCYxq_CUb7ZZAMIHgkKcasfNNLIxFq1xQ-8g2FyYUPdJZbXset-0I7TCb-VrBcH8DvJZDaQ"}},
	})
	assert.NoError(t, err)

	config := &jwtwallet.Config{}
	env.DoHttp(config)

	assert.Equal(t, 400, env.ClientRes.Status)
}

func TestMissingSubClaim(t *testing.T) {
	env, err := test.New(t, test.Request{
		Method:  "GET",
		Url:     "http://example.com",
		Headers: map[string][]string{"Authorization": {"eyJhbGciOiJFUzI1NksiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJwcm92ZW5hbmNlLmlvIiwiaWF0IjoiMTY0MDAzNDQzOTAwMCIsImV4cCI6IjE2NDAyMTQ0MzkwMDAiLCJhZGRyIjoidHAxZnllZGZlZ3pndzg4cWR4NDB4ejVwcXdqaHR0NmZqNmR3NnhwcGEifQ.JEw7aDqyX5IL3IwUo7SxaNZ6syh0bao8aBXAmo4UAJoSo5bB2BldLhhO4LatfdxDV_6juyaiwjQGv8e6FjwnIA"}},
	})
	assert.NoError(t, err)

	config := &jwtwallet.Config{}
	env.DoHttp(config)

	assert.Equal(t, 401, env.ClientRes.Status)
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	// coming soon!
}
func TestValidJwt(t *testing.T) {
	env, err := test.New(t, test.Request{
		Method:  "GET",
		Url:     "http://example.com",
		Headers: map[string][]string{"Authorization": {"eyJhbGciOiJFUzI1NksiLCJ0eXAiOiJKV1QifQ.eyJhZGRyIjoidHAxdXo1ZzcycHZmcmRubTlxbmpweXZzbndjNjRkNHd5Z3lxYW54MnQiLCJpc3MiOiJwcm92ZW5hbmNlLmlvIiwic3ViIjoiQWtqUlJMaGtzdU5rWF9Ba2pjVlpBX01ZOVpCNEd1cEtva2RlbU9LYnFRUFEiLCJleHAiOjQwNzA5MDg4MDAsImlhdCI6MTYwOTQ1OTIwMH0.zdGle-_d5qg0iVp_2gJ7zwBkqPCiO0YXDzCF37rviu0c7eP32qcCv5NTeKttKpXqPzaIWDkqxdrwYNlaSy26xQ"}},
	})
	assert.NoError(t, err)

	config := &jwtwallet.Config{
		RBAC: "http://localhost",
	}
	env.DoHttp(config)

	assert.Equal(t, 200, env.ClientRes.Status)
}
