package server

import (
	"encoding/json"
	"net/http"

	"github.com/findy-network/findy-agent-vault/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
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

type jwtChecker struct {
	checker *jwtmiddleware.JWTMiddleware
	secret  []byte
}

func newJWTChecker(jwtSecretKey string) *jwtChecker {
	secret := []byte(jwtSecretKey)
	return &jwtChecker{
		secret: secret,
		checker: jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			},
			SigningMethod:       jwt.SigningMethodHS256,
			EnableAuthOnOptions: true,
			Extractor: jwtmiddleware.FromFirst(
				jwtmiddleware.FromAuthHeader,
				jwtmiddleware.FromParameter("access_token"), // TODO: unsafe but needed for browser websocket auth
			),
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
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
			},
		}),
	}
}

func (j *jwtChecker) handler(next http.Handler) http.Handler {
	return j.checker.Handler(next)
}
