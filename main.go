package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/redis.v3"
)

var APILoader = APIDefinitionLoader{}
var config Specification

//loadConfigEnv loads environment variables
func loadConfigEnv() Specification {
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}

// initializeDatabase initializes a DB connection
func initializeDatabase() *DatabaseAccessor {
	accessor, err := NewServer(config.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	return accessor
}

// initializeRedis initializes a Redis connection
func initializeRedis() *redis.Client {
	log.Infof("Trying to connect to redis instance: %s", config.StorageDSN)
	return redis.NewClient(&redis.Options{
		Addr: config.StorageDSN,
	})
}

//loadAPIEndpoints register api endpoints
func loadAPIEndpoints(router *gin.Engine, apiManager *APIManager) {
	log.Debug("Loading API Endpoints")

	handler := AppsAPI{apiManager}
	group := router.Group("/apis")
	{
		group.GET("/", handler.Get())
		group.POST("/", handler.Post())
		group.GET("/:id", handler.GetBy())
		group.PUT("/:id", handler.PutBy())
		group.DELETE("/:id", handler.DeleteBy())
	}
}

func main() {
	loadConfigEnv()
	log.SetOutput(os.Stderr)
	router := gin.Default()

	if config.Debug {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	accessor := initializeDatabase()
	defer accessor.Close()

	database := Database{accessor}
	router.Use(database.Middleware())

	redisStorage := initializeRedis()
	defer redisStorage.Close()

	apiManager := NewAPIManager(router, redisStorage, accessor)
	apiManager.Load()
	loadAPIEndpoints(router, apiManager)

	router.Run(fmt.Sprintf(":%v", config.Port))
}
