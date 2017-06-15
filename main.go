package main

import (
	"log"

	"strconv"

	"github.com/labstack/echo"
	Redis "github.com/pwcong/simver/db/redis"
	Init "github.com/pwcong/simver/init"
	Middleware "github.com/pwcong/simver/middleware"
	Router "github.com/pwcong/simver/router"
	RPC "github.com/pwcong/simver/rpc"
)

func initMiddlewares(e *echo.Echo) {
	Middleware.Init(e)
}

func initRoutes(e *echo.Echo) {

	Router.Init(e)

}

func main() {

	// Intialize database connection
	err := Redis.Open(
		Init.Config.Databases.Redis.IP,
		Init.Config.Databases.Redis.Port,
		Init.Config.Databases.Redis.Password,
		Init.Config.Databases.Redis.DB)

	defer Redis.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Initialize RPC Server
	RPCServer := RPC.Server{
		IP:   Init.Config.Server.RPC.IP,
		Port: Init.Config.Server.RPC.Port,
	}
	go func(rpcServer *RPC.Server) {
		log.Fatal(rpcServer.Start())
	}(&RPCServer)

	// Initialize Resuful Server
	e := echo.New()
	initMiddlewares(e)
	initRoutes(e)

	log.Fatal(e.Start(Init.Config.Server.IP + ":" + strconv.Itoa(Init.Config.Server.Port)))

}
