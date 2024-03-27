package auth

import (
	"context"
	"encoding/json"
	"findigitalservice/internal/model/response"
	mUser "findigitalservice/internal/model/user"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Hash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func GenerateToken(user mUser.User) (response.TokenDto, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["username"] = user.Username
	claims["firstName"] = user.FirstName
	claims["lastName"] = user.LastName
	claims["email"] = user.Email
	claims["role"] = user.Role
	expire, err := strconv.ParseInt(os.Getenv("JWT_EXP"), 10, 64)
	if err != nil {
		expire = 15
	}
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(expire)).Unix()

	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return response.TokenDto{}, err
	}
	return response.TokenDto{Token: t}, nil
}

func UserLogin(ctx context.Context, logger *logrus.Logger) (mUser.User, error) {
	var u mUser.User
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		logger.Error("err FromContext: ", err)
		return u, err
	}
	uString, _ := json.Marshal(claims) // Convert to a json string
	err = json.Unmarshal(uString, &u)
	if err != nil {
		logger.Error("err Unmarshal: ", err)
		return u, err
	}
	return u, nil
}
