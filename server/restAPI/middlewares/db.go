package middlewares

import (
	db "github.com/Dpbm/quantumRestAPI/db"
	"github.com/gin-gonic/gin"
)

func DB(db *db.DB) gin.HandlerFunc {
	// check: https://stackoverflow.com/questions/34046194/how-to-pass-arguments-to-router-handlers-in-golang-using-gin-web-framework

	return func(context *gin.Context) {
		context.Set("db", db)
		context.Next()
	}
}
