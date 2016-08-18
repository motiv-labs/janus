package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/api-gateway/storage"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/redis.v3"
	"os"
	"github.com/gin-gonic/gin"
)

var APILoader = APIDefinitionLoader{}
var config = Specification{}

//loadConfigEnv loads environment variables
func loadConfigEnv() Specification {
	err := envconfig.Process("", &config)

	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}

// initializeDatabase initializes a DB connection
func initializeDatabase() *storage.DatabaseAccessor {
	accessor, err := storage.NewServer(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	return accessor
}

// initializeRedis initializes a Redis connection
func initializeRedis() *redis.Client {
	log.Infof("Trying to connect to redis instance: %s", config.Storage.DSN)
	return redis.NewClient(&redis.Options{
		Addr:     config.Storage.DSN,
		Password: config.Storage.Password,
		DB:       config.Storage.Database,
	})
}

//loadAPIEndpoints register api endpoints
func loadAPIEndpoints(router *gin.Engine) {
	log.Debug("Loading API Endpoints")

	handler := AppsAPI{}
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

	proxyRegister := &ProxyRegister{router}
	apiManager := APIManager{proxyRegister, redisStorage, accessor}
	apiManager.Load()
	loadAPIEndpoints(router)

	router.Run(fmt.Sprintf(":%v", config.Port))
}
