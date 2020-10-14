package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tools "github.com/findy-network/findy-agent-api/tools/resolver"
)

type JSONError struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type JSON struct {
	Data   map[string]interface{} `json:"data"`
	Errors *[]JSONError           `json:"errors"`
}

func queryJSON(content string) string {
	content = strings.Replace(content, "\t", "", -1)
	content = strings.Replace(content, "\n", " ", -1)
	return `{
		"query": "` + content + `"
		}`
}

func doQuery(query string) (payload JSON) {
	request, _ := http.NewRequest(http.MethodPost, "/query", strings.NewReader(queryJSON(query)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	Server(&tools.Resolver{}).ServeHTTP(response, request)

	bytes := response.Body.Bytes()
	fmt.Println(string(bytes))
	_ = json.Unmarshal(bytes, &payload)
	return
}

func TestServerForError(t *testing.T) {
	got := doQuery("{}")
	if len(*got.Errors) == 0 {
		t.Errorf("Expected errors, none found")
	}
}

func TestServerForSuccess(t *testing.T) {
	got := doQuery("{\n  __schema {\n    queryType {\n      name\n    }\n  }\n}")
	if _, ok := got.Data["__schema"]; !ok {
		t.Errorf("Expected response, none found")
	}
}
