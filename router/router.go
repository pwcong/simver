package router

import (
	"github.com/labstack/echo"
	IndexController "github.com/pwcong/simver/controller/index"
)

// Init initialize routes
func Init(e *echo.Echo) {

	e.GET(IndexController.URL_INDEX, IndexController.Default)

}
