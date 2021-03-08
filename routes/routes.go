package routes

import (
	"net/http"

	"github.com/towithyou/web_app/controllers/helloworld"

	"github.com/towithyou/web_app/settings"

	"github.com/gin-gonic/gin"
	"github.com/towithyou/web_app/logger"
)

func Setup() (r *gin.Engine) {
	r = gin.New()
	gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})

	helloworld.InitRoutes(r)
	return
}
