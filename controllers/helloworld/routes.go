package helloworld

import "github.com/gin-gonic/gin"

func InitRoutes(app *gin.Engine) {
	group := app.Group("api/v1/hello")
	group.GET("world", HelloWorld)
}
