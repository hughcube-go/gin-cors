package cors

import (
	"github.com/gin-gonic/gin"
	baseCors "github.com/hughcube-go/cors"
)

type Cors struct {
	baseCors.Cors
}

func (cors *Cors) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cors.Handler(c.Writer, c.Request) {
			c.Abort()
		} else {
			c.Next()
		}
	}
}
