package init

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type serverRPCConfig struct {
	IP   string
	Port int
}

type serverVertifyConfig struct {
	SigningKey  string
	VisitCounts int
	CheckCounts int
	ExpiredTime time.Duration
}

type serverConfig struct {
	IP      string
	Port    int
	Vertify serverVertifyConfig
	RPC     serverRPCConfig
}

type databaseRedisConfig struct {
	IP       string
	Port     int
	Password string
	DB       int
}
type databasesConfig struct {
	Redis databaseRedisConfig
}

type middlewareCORSConfig struct {
	Active       bool
	AllowOrigins []string
	AllowMethods []string
}
type middlewareLogConfig struct {
	Active bool
	Format string
	Output string
}

type middlewaresConfig struct {
	CORS middlewareCORSConfig
	Log  middlewareLogConfig
}

type config struct {
	Server      serverConfig
	Databases   databasesConfig
	Middlewares middlewaresConfig
}

var Config config

const DEFAULT_CONFIG = `
[server]
ip = "0.0.0.0"
port = 56789

    [server.vertify]
    signingKey = "simver"
    visitCounts = 100
    checkCounts = 1
    expiredTime = 86400000000000

    [server.rpc]
    ip = "0.0.0.0"
    port = 56780

[databases]
    [databases.redis]
    ip = "127.0.0.1"
    port = 6379
    password = ""
    db = 0

[middlewares]

    [middlewares.cors]
    active = true
    allowOrigins = ["*"]
    allowMethods = ["GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"]

    [middlewares.log]
    active = true
    format = "${time_rfc3339_nano} ${remote_ip} ${host} ${method} ${uri} ${status} ${latency_human} ${bytes_in} ${bytes_out}"
    output = "stdout"

`

func initConfig() {

	_, err := toml.DecodeFile(filepath.Join(filepath.Dir(os.Args[0]), "conf", "simver.toml"), &Config)

	if err == nil {

		log.Print("configuration has been loaded successfully.")

	} else {
		_, err := toml.Decode(DEFAULT_CONFIG, &Config)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			log.Print("failed to load custom configuration. use default.")
		}
	}

}

func init() {

	initConfig()

}
