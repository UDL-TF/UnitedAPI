package middleware

import (
	appctx "github.com/UDL-TF/UnitedAPI/internal/context"
	"github.com/gin-gonic/gin"
)

const AppContextKey = "app_context"

// InjectAppContext middleware injects the application context into gin.Context
func InjectAppContext(appContext *appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(AppContextKey, appContext)
		c.Next()
	}
}

// GetAppContext retrieves the application context from gin.Context
func GetAppContext(c *gin.Context) *appctx.AppContext {
	if ctx, exists := c.Get(AppContextKey); exists {
		return ctx.(*appctx.AppContext)
	}
	return nil
}
