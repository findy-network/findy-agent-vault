package server

import (
	"encoding/json"
	"net/http"

	"github.com/findy-network/findy-agent-vault/utils"
)

type JSONErrorExtension struct {
	Code string `json:"code"`
}

type JSONError struct {
	Message    string              `json:"message"`
	Path       []string            `json:"path"`
	Extensions *JSONErrorExtension `json:"extensions"`
}

type JSONPayload struct {
	Data   *map[string]interface{} `json:"data"`
	Errors *[]JSONError            `json:"errors"`
}

const (
	unauthenticated = "UNAUTHENTICATED"
)

func onAuthError(w http.ResponseWriter, r *http.Request, err string) {
	utils.LogLow().Infof("auth failed: %s", err)
	if r.Method == http.MethodPost {
		js, e := json.Marshal(
			&JSONPayload{
				Errors: &[]JSONError{{
					Extensions: &JSONErrorExtension{Code: unauthenticated},
				}},
			})

		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
		return
	}

	http.Error(w, err, http.StatusUnauthorized)
}
