package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateJWT(idUser uint, rolesUser string) (string, error)
	ExtractToken(t *jwt.Token) (uint, string, error)
}

func GenerateJWT(idUser uint, rolesUser string) (string, error) {
	var claim = jwt.MapClaims{}
	claim["id"] = idUser
	claim["role"] = rolesUser
	claim["iat"] = time.Now().UnixMilli()
	claim["exp"] = time.Now().Add(time.Minute * 30).UnixMilli()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	strToken, err := token.SignedString([]byte("$!1gnK3yyy!!!"))
	if err != nil {
		return "", err
	}

	return strToken, nil
}

func ExtractToken(t *jwt.Token) (uint, string, error) {
	var userID uint
	var rolesUser string
	expiredTime, err := t.Claims.GetExpirationTime()
	if err != nil {
		return 0, "", err
	}

	var eTime = *expiredTime

	if t.Valid && eTime.Compare(time.Now()) > 0 {
		id, ok := t.Claims.(jwt.MapClaims)["id"].(float64)
		if !ok {
			return 0, "", errors.New("tidak dapat mendapatkan ID dari token")
		}
		userID = uint(id)

		roles, ok := t.Claims.(jwt.MapClaims)["role"].(string)
		if !ok {
			return 0, "", errors.New("tidak dapat mendapatkan peran dari token")
		}
		rolesUser = roles

		return userID, rolesUser, nil
	}

	return 0, "", errors.New("token tidak valid")
}
