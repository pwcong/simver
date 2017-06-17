package index

import (
	"errors"

	"strconv"

	"regexp"

	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	Redis "github.com/pwcong/simver/db/redis"
	Init "github.com/pwcong/simver/init"
)

type VertifyKeyClaims struct {
	jwt.StandardClaims
}

func GenerateAndSetNewVertifyKey(ip string, visitCounts int, checkCounts int) (string, error) {

	vertifyKeyClaims := VertifyKeyClaims{
		jwt.StandardClaims{
			Issuer: ip,
		},
	}

	vertifyKeyToken := jwt.NewWithClaims(jwt.SigningMethodHS256, vertifyKeyClaims)
	vertifyKey, err := vertifyKeyToken.SignedString([]byte(Init.Config.Server.Vertify.SigningKey))

	if err != nil {
		return "", err
	}

	err = Redis.Client.Set(ip, strconv.Itoa(visitCounts)+":"+strconv.Itoa(checkCounts), Init.Config.Server.Vertify.ExpiredTime).Err()

	if err != nil {
		return "", err
	}

	return vertifyKey, nil

}

func GetVertifyKey(ip string) (string, error) {

	record, err := Redis.Client.Get(ip).Result()

	if err != nil {

		return GenerateAndSetNewVertifyKey(ip, 1, 0)

	}

	if record == "" {
		return GenerateAndSetNewVertifyKey(ip, 1, 0)

	}

	if matched, err := regexp.Match(`^\d+:\d+$`, []byte(record)); !matched || err != nil {
		return GenerateAndSetNewVertifyKey(ip, 1, 0)
	}

	recordValues := strings.Split(record, ":")

	visitCountsValue := recordValues[0]
	visitCounts, err := strconv.Atoi(visitCountsValue)
	if err != nil {
		visitCounts = 1
	}

	if visitCounts >= Init.Config.Server.Vertify.VisitCounts {
		return "", errors.New("visit counts limit")
	}

	checkCountsValue := recordValues[1]
	checkCounts, err := strconv.Atoi(checkCountsValue)
	if err != nil {
		checkCounts = 0
	}

	if checkCounts >= Init.Config.Server.Vertify.CheckCounts {
		return "", errors.New("check counts limit")
	}

	return GenerateAndSetNewVertifyKey(ip, visitCounts+1, checkCounts)

}
