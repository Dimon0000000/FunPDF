package internal

import "github.com/gin-gonic/gin"

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Setup(e *gin.Engine) {
	api := e.Group("/api")
	{
		api.GET("/index")

		file := api.Group("/files")
		{
			file.POST("/upload")
			file.POST("/:file_id")
			file.GET("/:file_id") // Get the file
		}

		providers := api.Group("/providers")
		{
			providers.GET("/") // GET all providers
			providers.POST("/:provider_name/create")
			providers.PUT("/:provider_name/update")
		}

		translate := api.Group("/translaters")
		{
			translate.POST("/:translater_name/completion")
		}
	}

}
