package routes

import "github.com/gin-gonic/gin"

func Page404(context *gin.Context) {
	context.JSON(404, map[string]string{"msg": "Page not found"})
}
