package routes

import (
	"github.com/gin-gonic/gin"
	"self-hosted-cloud/server/api"
	"self-hosted-cloud/server/middlewares"
)

func run(handler func(c *gin.Context) (int, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, err := handler(c)
		if err != nil {
			c.AbortWithError(code, err)
		}
	}
}

func LoadRoutes(router *gin.Engine) {
	g := router.Group("/admin")
	{
		g.Use(middlewares.AdminMiddleware())

		g.GET("/demo", run(api.GetDemoMode))
		g.POST("/demo", run(api.EnableDemoMode))
		g.POST("/reset", run(api.HardReset))
	}

	g = router.Group("/auth")
	{
		g.POST("/logout", run(api.Logout))
		g.GET("/github/login", run(api.Login))
		g.POST("/github/callback", run(api.Callback))
	}

	g = router.Group("/storage")
	{
		g.GET("/:bucket_uuid", run(api.GetNodes))
		g.GET("/:bucket_uuid/recent", run(api.GetRecentFiles))
		g.GET("/:bucket_uuid/bin", run(api.GetBin))
		g.DELETE("/:bucket_uuid/bin", run(api.EmptyBin))
		g.PUT("/:bucket_uuid", run(api.CreateNode))
		g.DELETE("/:bucket_uuid", run(api.DeleteNodes))
		g.PATCH("/:bucket_uuid", run(api.RenameNode))
		g.GET("/:bucket_uuid/download", run(api.DownloadNodes))
		g.POST("/:bucket_uuid/upload", run(api.UploadNode))
		g.GET("/bucket", run(api.GetBucket))
	}

	g = router.Group("/user")
	{
		g.GET("/", run(api.GetUser))
		g.GET("/:username", run(api.GetUser))
	}
}
