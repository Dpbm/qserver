package utils

import (
	"github.com/Dpbm/quantumRestAPI/db"
	"github.com/gin-gonic/gin"
)

func GetDBFromContext(context *gin.Context) (*db.DB, bool) {
	db, ok := context.MustGet("db").(*db.DB)
	return db, ok
}
