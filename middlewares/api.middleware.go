package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Hash(secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(secret))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func ApiMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var apiKey string = ctx.Request.Header.Get("x-api-key")

		if apiKey != Hash(os.Getenv("JWT_SECRET")) {
			ctx.JSON(http.StatusUnauthorized,
				gin.H{"status": http.StatusUnauthorized,
					"message": "not authorized"})
			ctx.Abort()
			return
		}
		ctx.Set("authorized", true)
		ctx.Next()
	}
}
