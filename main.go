package main

import (
	"echoapp/config"
	"echoapp/container"
	"echoapp/database"
	"echoapp/logger"
	"echoapp/repo"
	"echoapp/repository"
	"echoapp/router"
	"embed"

	"log"

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
	db, err := database.CreateDB()

	if err != nil {
		log.Fatal(err)
		return
	}
	conf, env := config.Load(yamlFile)
	logger := logger.InitLogger(env, zapYamlFile)
	rep := repository.NewLocationRepository(logger, conf, db)
	repo :=
		repo.NewRepo(db)

	container := container.NewContainer(rep, repo, conf, logger, env)
	//ph := controller.NewContinentHandler(e.StdLogger, container)
	//e.GET("/continents", ph.GetContinents)
	//e.GET("/countries", ph.GetCountries)
	router.Init(e, container, &repo)
	e.Logger.Fatal(e.Start(":1323"))
}
