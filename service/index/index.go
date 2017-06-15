package index

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	Redis "github.com/pwcong/simver/db/redis"
	Init "github.com/pwcong/simver/init"
)

type VertifyKeyClaims struct {
	VisitCounts float64 `json:"visitCounts"`
	CheckCounts float64 `json:"checkCounts"`
	jwt.StandardClaims
}

func GenerateAndSetNewVertifyKey(ip string, visitCounts float64, checkCounts float64) (string, error) {

	vertifyKeyClaims := VertifyKeyClaims{
		visitCounts,
		checkCounts,
		jwt.StandardClaims{
			Issuer: ip,
		},
	}

	vertifyKeyToken := jwt.NewWithClaims(jwt.SigningMethodHS256, vertifyKeyClaims)
	vertifyKey, err := vertifyKeyToken.SignedString([]byte(Init.Config.Server.Vertify.SigningKey))

	if err != nil {
		return "", err
	}

	err = Redis.Client.Set(ip, vertifyKey, Init.Config.Server.Vertify.ExpiredTime).Err()

	if err != nil {
		return "", nil
	}

	return vertifyKey, nil

}

func GetVertifyKey(ip string) (string, error) {

	vertifyKey, err := Redis.Client.Get(ip).Result()

	if err != nil {

		return GenerateAndSetNewVertifyKey(ip, 1, 0)

	}

	vertifyKeyToken, err := jwt.Parse(vertifyKey, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(Init.Config.Server.Vertify.SigningKey), nil
	})

	if err != nil {

		return GenerateAndSetNewVertifyKey(ip, 1, 0)

	}

	if vertifyKeyClaims, ok := vertifyKeyToken.Claims.(jwt.MapClaims); ok && vertifyKeyToken.Valid {

		visitCounts, ok := vertifyKeyClaims["visitCounts"].(float64)

		if !ok {
			visitCounts = 0
		}

		if int(visitCounts) >= Init.Config.Server.Vertify.VisitCounts {
			return "", errors.New("visit counts limit")
		}

		checkCounts, ok := vertifyKeyClaims["checkCounts"].(float64)
		if !ok {
			checkCounts = 0
		}

		if int(checkCounts) >= Init.Config.Server.Vertify.CheckCounts {
			return "", errors.New("check counts limit")
		}

		return GenerateAndSetNewVertifyKey(ip, visitCounts+1, checkCounts)

	}
	return GenerateAndSetNewVertifyKey(ip, 1, 0)

}
