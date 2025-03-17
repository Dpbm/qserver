package routes

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1
// @version 1.0
// @Summary get health
// @Schemes http
// @Description healthcheck, a route to test if everything is ok (like a ping command)
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func GetHealth(context *gin.Context) {
	context.JSON(200, map[string]string{"status": "ok"})
}
