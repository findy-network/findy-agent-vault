package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang/glog"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

const jwtSecret = "supersecret"

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
			glog.V(3).Infof("auth failed: %s", err)
			if r.Method == http.MethodPost {
				code := map[string]interface{}{"code": "UNAUTHENTICATED"}
				extensions := map[string]interface{}{"extensions": &code}
				errs := map[string]interface{}{"errors": []interface{}{&extensions}}
				js, e := json.Marshal(&errs)

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
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return signer.SignedString([]byte(jwtSecret))
}
