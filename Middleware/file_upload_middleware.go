package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func FileUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Bad request",
			})
			return
		}
		defer file.Close() // close file properly

		//remove the file extension
		fileName := strings.Split(header.Filename, ".")[0]
		c.Set("filePath", fileName)
		c.Set("file", file)

		c.Next()
	}
}
