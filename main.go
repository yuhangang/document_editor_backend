package main

import (
	"echoapp/config"
	"echoapp/container"
	"echoapp/logger"
	"echoapp/migration"
	"echoapp/repo"

	"echoapp/router"
	"embed"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

//go:embed application.*.yml
var yamlFile embed.FS

//go:embed zaplogger.*.yml
var zapYamlFile embed.FS

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgresPass"
	dbname   = "pqgotest"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	if err != nil {
		log.Fatal(err)
		return
	}

	conf, env, commandArgs := config.Load(yamlFile)
	logger := logger.InitLogger(env, zapYamlFile)
	db, err := migration.CreateDB(commandArgs)
	repository := repo.NewDocumentFileRepository(logger, conf, db)

	bigCache, bigCacheInitError := bigcache.NewBigCache(bigcache.DefaultConfig(24 * time.Hour))
	if bigCacheInitError != nil {
		log.Fatal(bigCacheInitError)
		return
	}

	container := container.NewContainer(&repository, conf, bigCache, logger, env)

	migration.InitMasterData(container)
	router.Init(e, container)
	e.Logger.Fatal(e.Start(":1323"))
}
