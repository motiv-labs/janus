package janus

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellofresh/janus/config"
)

func Home(app config.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, fmt.Sprintf("Welcome to %s, this is version %s", app.Name, app.Version))
	}
}
