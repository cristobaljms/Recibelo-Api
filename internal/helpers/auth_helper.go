package helpers

import (
	"fmt"
	"net/http"
	"recibe_me/configs"
	"recibe_me/internal/models"
	"recibe_me/pkg/crypto"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

// AuthenticatedUser contains the user authenticated
var AuthenticatedUser models.User

// GetUserByToken returns a User by token
func GetUserByToken(tokenString string) (models.User, error) {

	user := models.User{}

	err := crypto.ValidateToken([]byte(tokenString))

	if err != nil {
		return user, fmt.Errorf("Invalid Token")
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.SecurityCfg.TokenSecret), nil
	})

	if err != nil {
		return user, err
	}

	if claims["jti"] == nil || !bson.IsObjectIdHex(claims["jti"].(string)) {
		return user, fmt.Errorf("Invalid Id")
	}

	oid := bson.ObjectIdHex(claims["jti"].(string))

	err = UsersCollection.FindId(oid).One(&user)

	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserFromRequest returns a User from request by token
func GetUserFromRequest(request *http.Request) (models.User, error) {

	user := models.User{}

	token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")

	user, err := GetUserByToken(token)

	if err != nil {
		return user, err
	}

	return user, nil
}
