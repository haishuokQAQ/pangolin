package utils

import (
	"context"
	"github.com/gin-gonic/gin"
)

func CreateContextFromGinContext(c *gin.Context) context.Context {
	return context.Background()
}
