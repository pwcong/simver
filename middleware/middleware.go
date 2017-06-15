package middleware

import (
	"log"
	"os"
	"path/filepath"

	InitConfig "github.com/pwcong/simver/init"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Init initialize middlewares
func Init(e *echo.Echo) {

	// initialize log middleware
	if InitConfig.Config.Middlewares.Log.Active {

		if InitConfig.Config.Middlewares.Log.Output == "file" {

			logDir := filepath.Join(filepath.Dir(os.Args[0]), "log")
			if _, err := os.Stat(logDir); err != nil {
				err := os.MkdirAll(logDir, 0666)
				if err != nil {
					log.Fatal(err.Error())
				}

			}

			logPath := filepath.Join(logDir, "server.log")
			logOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatal(err.Error())
			}

			e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
				Output: logOutput,
				Format: InitConfig.Config.Middlewares.Log.Format + "\n",
			}))

		} else if InitConfig.Config.Middlewares.Log.Output == "stdout" {
			e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
				Format: InitConfig.Config.Middlewares.Log.Format + "\n",
			}))
		}
	}

	// initialize cors middleware configuration
	if InitConfig.Config.Middlewares.CORS.Active {

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: InitConfig.Config.Middlewares.CORS.AllowOrigins,
			AllowMethods: InitConfig.Config.Middlewares.CORS.AllowMethods,
		}))
	}

}
