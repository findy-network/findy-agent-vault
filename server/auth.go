package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
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
	jwtSecret       = "mySuperSecretKeyLol"
	unauthenticated = "UNAUTHENTICATED"
	hoursInDay      = 24
	hoursForTest    = 2
)

// JWTChecker checks the token for all requests
// The authentication error is generated here instead of resolvers to make sure all resolvers use authentication.
// Error should be in compatible GQL format so that frontend frameworks succeed in parsing.
// TODO: move authentication to resolvers so that errors are generated at correct level?
func jwtChecker(next http.Handler) http.Handler {
	checker := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
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
	})
	return checker.Handler(next)
}

func CreateToken(id string) (string, error) {
	signer := createToken(id, time.Hour*hoursInDay)
	return signer.SignedString([]byte(jwtSecret))
}

func CreateTestToken(id string) *jwt.Token {
	token := createToken(id, time.Hour*hoursForTest)
	str, _ := token.SignedString([]byte(jwtSecret))
	token.Raw = str
	return token
}

func createTokenString(id string, duration time.Duration) (string, error) {
	signer := createToken(id, duration)
	return signer.SignedString([]byte(jwtSecret))
}

func createToken(id string, duration time.Duration) *jwt.Token {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["un"] = fake.FakeCloudDID
	claims["label"] = "minnie mouse"
	claims["exp"] = time.Now().Add(duration).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func createTestToken() string {
	token, _ := createTokenString("test", time.Hour*hoursForTest)
	return token
}
