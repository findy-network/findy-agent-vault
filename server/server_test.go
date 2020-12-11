package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/findy-network/findy-agent-vault/resolver"
)

var testToken string = createTestToken()

const testQuery = "{\n  __schema {\n    queryType {\n      name\n    }\n  }\n}"

func queryJSON(content string) string {
	content = strings.Replace(content, "\t", "", -1)
	content = strings.Replace(content, "\n", " ", -1)
	return `{
		"query": "` + content + `"
		}`
}

func doQuery(query string, auth bool) (payload JSONPayload) {
	request, _ := http.NewRequestWithContext(context.TODO(), http.MethodPost, "/query", strings.NewReader(queryJSON(query)))
	request.Header.Set("Content-Type", "application/json")
	if auth {
		request.Header.Set("Authorization", "Bearer "+testToken)
	}
	response := httptest.NewRecorder()

	Server(&resolver.Resolver{}).ServeHTTP(response, request)

	bytes := response.Body.Bytes()
	fmt.Println(string(bytes))
	_ = json.Unmarshal(bytes, &payload)
	return
}

func doAuthQuery(query string) (payload JSONPayload) {
	return doQuery(query, true)
}

func TestServerForError(t *testing.T) {
	got := doAuthQuery("{}")
	if len(*got.Errors) == 0 {
		t.Errorf("Expected errors, none found")
	}
}

func TestServerForAuth(t *testing.T) {
	got := doQuery(testQuery, false)
	if len(*got.Errors) == 0 || (*got.Errors)[0].Extensions.Code != unauthenticated {
		t.Errorf("Expected UNAUTHENTICATED error, none found")
	}
}

func TestServerForSuccess(t *testing.T) {
	got := doAuthQuery(testQuery)
	if _, ok := (*got.Data)["__schema"]; !ok {
		t.Errorf("Expected response, none found")
	}
}
