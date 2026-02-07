package routes

import (
	"week12/handle"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ws", handle.HandleConnections)
	r.POST("/upload/image", handle.HandleImageUpload)
}
