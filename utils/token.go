package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type UserToken struct {
	Label   string
	AgentID string
	Token   string
}

func ParseToken(ctx context.Context) (*UserToken, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("no authenticated user found")
	}

	jwtToken, ok := user.(*jwt.Token)
	if !ok {
		return nil, errors.New("no authenticated user found")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("no claims found for token")
	}

	caDID, ok := claims["un"].(string)
	if !ok || caDID == "" {
		return nil, errors.New("no cloud agent DID found for token")
	}

	if jwtToken.Raw == "" {
		return nil, errors.New(fmt.Sprintf("no raw token found for user %s", caDID))
	}

	label := "n/a"
	if l, ok := claims["label"].(string); ok {
		label = l
	}

	return &UserToken{
		AgentID: caDID, Label: label, Token: jwtToken.Raw,
	}, nil
}
