package index

import (
	"net/http"

	"github.com/labstack/echo"

	IndexService "github.com/pwcong/simver/service/index"
)

const (
	URL_INDEX = "/"
)

type JSONResponseWithVertifyKey struct {
	JSONResponse
	Key string `json:"key"`
}

type JSONResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func Default(c echo.Context) error {

	vertifyKey, err := IndexService.GetVertifyKey(c.RealIP())

	if err != nil {
		return c.JSON(http.StatusForbidden, JSONResponse{
			Code:    http.StatusForbidden,
			Message: err.Error(),
		})

	}

	return c.JSON(http.StatusOK, JSONResponseWithVertifyKey{
		JSONResponse: JSONResponse{
			Code:    http.StatusOK,
			Message: "",
		},
		Key: vertifyKey,
	})
}
