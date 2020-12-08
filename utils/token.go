package utils

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/lainio/err2"
)

type UserToken struct {
	Label   string
	AgentID string
	Token   string
}

func ParseToken(ctx context.Context) *UserToken {
	defer err2.Catch(func(err error) {
		panic(err)
	})

	user := ctx.Value("user")
	if user == nil {
		err2.Check(fmt.Errorf("no authenticated user found"))
	}

	jwtToken, ok := user.(*jwt.Token)
	if !ok {
		err2.Check(fmt.Errorf("no jwt token found for user"))
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		err2.Check(fmt.Errorf("no claims found for token"))
	}

	caDID, ok := claims["un"].(string)
	if !ok || caDID == "" {
		err2.Check(fmt.Errorf("no cloud agent DID found for token"))
	}

	if jwtToken.Raw == "" {
		err2.Check(fmt.Errorf("no raw token found for user %s", caDID))
	}

	label := "n/a"
	if l, ok := claims["label"].(string); ok {
		label = l
	}

	return &UserToken{
		AgentID: caDID, Label: label, Token: jwtToken.Raw,
	}
}
