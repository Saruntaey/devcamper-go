package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Payload struct {
	Id string
	jwt.StandardClaims
}

func GetJwt(id string) (string, error) {
	tokenDuration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	payload := Payload{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(tokenDuration)).Unix(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	ss, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("cannot signed token: %w", err)
	}
	return ss, nil
}

func checkJwt(ss string) bool {
	t, err := jwt.ParseWithClaims(ss, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return []byte{}, errors.New("signed algo not match")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return err == nil && t.Valid
}
